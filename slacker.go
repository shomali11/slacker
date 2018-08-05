package slacker

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/shomali11/proper"
)

const (
	space               = " "
	dash                = "-"
	newLine             = "\n"
	invalidToken        = "invalid token"
	helpCommand         = "help"
	directChannelMarker = "D"
	userMentionFormat   = "<@%s>"
	codeMessageFormat   = "`%s`"
	boldMessageFormat   = "*%s*"
	italicMessageFormat = "_%s_"
	slackBotUser        = "USLACKBOT"
)

// NewClient creates a new client using the Slack API
func NewClient(token string) *Slacker {
	client := slack.New(token)
	slacker := &Slacker{
		client: client,
		rtm:    client.NewRTM(),
	}
	return slacker
}

// Slacker contains the Slack API, botCommands, and handlers
type Slacker struct {
	client                *slack.Client
	rtm                   *slack.RTM
	botCommands           []BotCommand
	requestConstructor    func(ctx context.Context, event *slack.MessageEvent, properties *proper.Properties) Request
	responseConstructor   func(channel string, client *slack.Client, rtm *slack.RTM) ResponseWriter
	initHandler           func()
	errorHandler          func(err string)
	helpHandler           func(request Request, response ResponseWriter)
	defaultMessageHandler func(request Request, response ResponseWriter)
	defaultEventHandler   func(interface{})
}

// BotCommands returns Bot Commands
func (s *Slacker) BotCommands() []BotCommand {
	return s.botCommands
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
func (s *Slacker) CustomRequest(requestConstructor func(ctx context.Context, event *slack.MessageEvent, properties *proper.Properties) Request) {
	s.requestConstructor = requestConstructor
}

// CustomResponse creates a new response writer
func (s *Slacker) CustomResponse(responseConstructor func(channel string, client *slack.Client, rtm *slack.RTM) ResponseWriter) {
	s.responseConstructor = responseConstructor
}

// DefaultCommand handle messages when none of the commands are matched
func (s *Slacker) DefaultCommand(defaultMessageHandler func(request Request, response ResponseWriter)) {
	s.defaultMessageHandler = defaultMessageHandler
}

// DefaultEvent handle events when an unknown event is seen
func (s *Slacker) DefaultEvent(defaultEventHandler func(interface{})) {
	s.defaultEventHandler = defaultEventHandler
}

// Help handle the help message, it will use the default if not set
func (s *Slacker) Help(helpHandler func(request Request, response ResponseWriter)) {
	s.helpHandler = helpHandler
}

// Command define a new command and append it to the list of existing commands
func (s *Slacker) Command(usage string, description string, handler func(request Request, response ResponseWriter)) {
	s.botCommands = append(s.botCommands, NewBotCommand(usage, description, handler))
}

// Listen receives events from Slack and each is handled as needed
func (s *Slacker) Listen(ctx context.Context) error {
	s.prependHelpHandle()

	go s.rtm.ManageConnection()

	for msg := range s.rtm.IncomingEvents {
		select {
		case <-ctx.Done():
			s.rtm.Disconnect()
			return ctx.Err()
		default:
			switch event := msg.Data.(type) {
			case *slack.ConnectedEvent:
				if s.initHandler == nil {
					continue
				}
				go s.initHandler()

			case *slack.MessageEvent:
				if s.isFromBot(event) {
					continue
				}

				if !s.isBotMentioned(event) && !s.isDirectMessage(event) {
					continue
				}
				go s.handleMessage(ctx, event)

			case *slack.RTMError:
				if s.errorHandler == nil {
					continue
				}
				go s.errorHandler(event.Error())

			case *slack.InvalidAuthEvent:
				return errors.New(invalidToken)

			default:
				if s.defaultEventHandler == nil {
					continue
				}
				go s.defaultEventHandler(event)
			}
		}
	}
	return nil
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

func (s *Slacker) handleMessage(ctx context.Context, event *slack.MessageEvent) {
	if s.requestConstructor == nil {
		s.requestConstructor = NewRequest
	}

	if s.responseConstructor == nil {
		s.responseConstructor = NewResponse
	}

	response := s.responseConstructor(event.Channel, s.client, s.rtm)

	for _, cmd := range s.botCommands {
		parameters, isMatch := cmd.Match(event.Text)
		if !isMatch {
			continue
		}

		request := s.requestConstructor(ctx, event, parameters)
		cmd.Execute(request, response)
		return

	}

	if s.defaultMessageHandler != nil {
		request := s.requestConstructor(ctx, event, &proper.Properties{})
		s.defaultMessageHandler(request, response)
	}
}

func (s *Slacker) defaultHelp(request Request, response ResponseWriter) {
	helpMessage := empty
	for _, command := range s.botCommands {
		tokens := command.Tokenize()
		for _, token := range tokens {
			if token.IsParameter {
				helpMessage += fmt.Sprintf(codeMessageFormat, token.Word) + space
			} else {
				helpMessage += fmt.Sprintf(boldMessageFormat, token.Word) + space
			}
		}
		helpMessage += dash + space + fmt.Sprintf(italicMessageFormat, command.Description()) + newLine
	}
	response.Reply(helpMessage)
}

func (s *Slacker) prependHelpHandle() {
	if s.helpHandler == nil {
		s.helpHandler = s.defaultHelp
	}
	s.botCommands = append([]BotCommand{NewBotCommand(helpCommand, helpCommand, s.helpHandler)}, s.botCommands...)
}
