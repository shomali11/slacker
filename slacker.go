package slacker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/shomali11/proper"
	"github.com/slack-go/slack"
)

const (
	space               = " "
	dash                = "-"
	star                = "*"
	newLine             = "\n"
	invalidToken        = "invalid token"
	helpCommand         = "help"
	directChannelMarker = "D"
	userMentionFormat   = "<@%s>"
	codeMessageFormat   = "`%s`"
	boldMessageFormat   = "*%s*"
	italicMessageFormat = "_%s_"
	quoteMessageFormat  = ">_*Example:* %s_"
	authorizedUsersOnly = "Authorized users only"
	slackBotUser        = "USLACKBOT"
)

var (
	unAuthorizedError = errors.New("You are not authorized to execute this command")
)

// NewClient creates a new client using the Slack API
func NewClient(token string, options ...ClientOption) *Slacker {
	defaults := newClientDefaults(options...)

	client := slack.New(token, slack.OptionDebug(defaults.Debug))
	slacker := &Slacker{
		client:            client,
		rtm:               client.NewRTM(),
		commandChannel:    make(chan *CommandEvent, 100),
		unAuthorizedError: unAuthorizedError,
		eventHandler:      DefaultEventHandler,
	}
	return slacker
}

// Slacker contains the Slack API, botCommands, and handlers
type Slacker struct {
	client                *slack.Client
	rtm                   *slack.RTM
	botCommands           []BotCommand
	botContextConstructor func(ctx context.Context, event *slack.MessageEvent, client *slack.Client, rtm *slack.RTM) BotContext
	requestConstructor    func(botCtx BotContext, properties *proper.Properties) Request
	responseConstructor   func(botCtx BotContext) ResponseWriter
	initHandler           func()
	errorHandler          func(err string)
	helpDefinition        *CommandDefinition
	defaultMessageHandler func(botCtx BotContext, request Request, response ResponseWriter)
	fallbackEventHandler  func(interface{})
	eventHandler          EventHandler
	unAuthorizedError     error
	commandChannel        chan *CommandEvent
}

// BotCommands returns Bot Commands
func (s *Slacker) BotCommands() []BotCommand {
	return s.botCommands
}

// Client returns the internal slack.Client of Slacker struct
func (s *Slacker) Client() *slack.Client {
	return s.client
}

// RTM returns returns the internal slack.RTM of Slacker struct
func (s *Slacker) RTM() *slack.RTM {
	return s.rtm
}

// Init handle the event when the bot is first connected
func (s *Slacker) Init(initHandler func()) {
	s.initHandler = initHandler
}

// Err handle when errors are encountered
func (s *Slacker) Err(errorHandler func(err string)) {
	s.errorHandler = errorHandler
}

// CustomRequest creates a new request
func (s *Slacker) CustomRequest(requestConstructor func(botCtx BotContext, properties *proper.Properties) Request) {
	s.requestConstructor = requestConstructor
}

// CustomResponse creates a new response writer
func (s *Slacker) CustomResponse(responseConstructor func(botCtx BotContext) ResponseWriter) {
	s.responseConstructor = responseConstructor
}

// DefaultCommand handle messages when none of the commands are matched
func (s *Slacker) DefaultCommand(defaultMessageHandler func(botCtx BotContext, request Request, response ResponseWriter)) {
	s.defaultMessageHandler = defaultMessageHandler
}

// DefaultEvent handle events when an unknown event is seen
// Deprecated. Use FallbackEvent instead.
func (s *Slacker) DefaultEvent(defaultEventHandler func(interface{})) {
	s.fallbackEventHandler = defaultEventHandler
}

// FallbackEvent handle events when an unknown event is seen
func (s *Slacker) FallbackEvent(fallbackEventHandler func(interface{})) {
	s.fallbackEventHandler = fallbackEventHandler
}

// UnAuthorizedError error message
func (s *Slacker) UnAuthorizedError(unAuthorizedError error) {
	s.unAuthorizedError = unAuthorizedError
}

// Help handle the help message, it will use the default if not set
func (s *Slacker) Help(definition *CommandDefinition) {
	s.helpDefinition = definition
}

// Command define a new command and append it to the list of existing commands
func (s *Slacker) Command(usage string, definition *CommandDefinition) {
	s.botCommands = append(s.botCommands, NewBotCommand(usage, definition))
}

// CommandEvents returns read only command events channel
func (s *Slacker) CommandEvents() <-chan *CommandEvent {
	return s.commandChannel
}

// Listen receives events from Slack and each is handled as needed
func (s *Slacker) Listen(ctx context.Context) error {
	s.prependHelpHandle()

	go s.rtm.ManageConnection()
	for {
		select {
		case <-ctx.Done():
			s.rtm.Disconnect()
			return ctx.Err()
		case msg, ok := <-s.rtm.IncomingEvents:
			if !ok {
				return nil
			}
			if err := s.eventHandler(ctx, s, msg); err != nil {
				return err
			}
		}
	}
}

// GetUserInfo retrieve complete user information
func (s *Slacker) GetUserInfo(user string) (*slack.User, error) {
	return s.client.GetUserInfo(user)
}

func (s *Slacker) sendMessage(text string, channel string) {
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(text, channel))
}

func (s *Slacker) isFromBot(event *slack.MessageEvent) bool {
	info := s.rtm.GetInfo()
	return len(event.User) == 0 || event.User == slackBotUser || event.User == info.User.ID || len(event.BotID) > 0
}

func (s *Slacker) isBotMentioned(event *slack.MessageEvent) bool {
	info := s.rtm.GetInfo()
	return strings.Contains(event.Text, fmt.Sprintf(userMentionFormat, info.User.ID))
}

func (s *Slacker) isDirectMessage(event *slack.MessageEvent) bool {
	return strings.HasPrefix(event.Channel, directChannelMarker)
}

func (s *Slacker) handleMessage(ctx context.Context, message *slack.MessageEvent) {
	if s.botContextConstructor == nil {
		s.botContextConstructor = NewBotContext
	}

	if s.requestConstructor == nil {
		s.requestConstructor = NewRequest
	}

	if s.responseConstructor == nil {
		s.responseConstructor = NewResponse
	}

	botCtx := s.botContextConstructor(ctx, message, s.client, s.rtm)
	response := s.responseConstructor(botCtx)

	for _, cmd := range s.botCommands {
		parameters, isMatch := cmd.Match(message.Text)
		if !isMatch {
			continue
		}

		request := s.requestConstructor(botCtx, parameters)
		if cmd.Definition().AuthorizationFunc != nil && !cmd.Definition().AuthorizationFunc(botCtx, request) {
			response.ReportError(s.unAuthorizedError)
			return
		}

		select {
		case s.commandChannel <- NewCommandEvent(cmd.Usage(), parameters, message):
		default:
			// full channel, dropped event
		}

		cmd.Execute(botCtx, request, response)
		return
	}

	if s.defaultMessageHandler != nil {
		request := s.requestConstructor(botCtx, &proper.Properties{})
		s.defaultMessageHandler(botCtx, request, response)
	}
}

func (s *Slacker) defaultHelp(botCtx BotContext, request Request, response ResponseWriter) {
	authorizedCommandAvailable := false
	helpMessage := empty
	for _, command := range s.botCommands {
		tokens := command.Tokenize()
		for _, token := range tokens {
			if token.IsParameter() {
				helpMessage += fmt.Sprintf(codeMessageFormat, token.Word) + space
			} else {
				helpMessage += fmt.Sprintf(boldMessageFormat, token.Word) + space
			}
		}

		if len(command.Definition().Description) > 0 {
			helpMessage += dash + space + fmt.Sprintf(italicMessageFormat, command.Definition().Description)
		}

		if command.Definition().AuthorizationFunc != nil {
			authorizedCommandAvailable = true
			helpMessage += space + fmt.Sprintf(codeMessageFormat, star)
		}

		helpMessage += newLine

		if len(command.Definition().Example) > 0 {
			helpMessage += fmt.Sprintf(quoteMessageFormat, command.Definition().Example) + newLine
		}
	}

	if authorizedCommandAvailable {
		helpMessage += fmt.Sprintf(codeMessageFormat, star+space+authorizedUsersOnly) + newLine
	}
	response.Reply(helpMessage)
}

func (s *Slacker) prependHelpHandle() {
	if s.helpDefinition == nil {
		s.helpDefinition = &CommandDefinition{}
	}

	if s.helpDefinition.Handler == nil {
		s.helpDefinition.Handler = s.defaultHelp
	}

	if len(s.helpDefinition.Description) == 0 {
		s.helpDefinition.Description = helpCommand
	}

	s.botCommands = append([]BotCommand{NewBotCommand(helpCommand, s.helpDefinition)}, s.botCommands...)
}
