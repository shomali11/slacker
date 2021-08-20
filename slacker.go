package slacker

import (
	"context"
	"errors"
	"fmt"

	"github.com/shomali11/proper"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
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
func NewClient(botToken, appToken string, options ...ClientOption) *Slacker {
	defaults := newClientDefaults(options...)

	api := slack.New(
		botToken,
		slack.OptionDebug(defaults.Debug),
		slack.OptionAppLevelToken(appToken),
		slack.OptionLog(defaults.Logger),
	)

	smc := socketmode.New(
		api,
		socketmode.OptionDebug(defaults.Debug),
	)
	slacker := &Slacker{
		client:            api,
		socketModeClient:  smc,
		commandChannel:    make(chan *CommandEvent, 100),
		unAuthorizedError: unAuthorizedError,
		log:               internalLog{logger: defaults.Logger},
	}
	return slacker
}

// Slacker contains the Slack API, botCommands, and handlers
type Slacker struct {
	client                  *slack.Client
	socketModeClient        *socketmode.Client
	botCommands             []BotCommand
	botContextConstructor   func(ctx context.Context, api *slack.Client, client *socketmode.Client, evt *MessageEvent) BotContext
	requestConstructor      func(botCtx BotContext, properties *proper.Properties) Request
	responseConstructor     func(botCtx BotContext) ResponseWriter
	initHandler             func()
	errorHandler            func(err string)
	interactiveEventHandler func(*Slacker, *socketmode.Event, *slack.InteractionCallback)
	helpDefinition          *CommandDefinition
	defaultMessageHandler   func(botCtx BotContext, request Request, response ResponseWriter)
	defaultEventHandler     func(interface{})
	unAuthorizedError       error
	commandChannel          chan *CommandEvent
	appID                   string
	log                     internalLog
}

// BotCommands returns Bot Commands
func (s *Slacker) BotCommands() []BotCommand {
	return s.botCommands
}

// Client returns the internal slack.Client of Slacker struct
func (s *Slacker) Client() *slack.Client {
	return s.client
}

// SocketMode returns the internal socketmode.Client of Slacker struct
func (s *Slacker) SocketMode() *socketmode.Client {
	return s.socketModeClient
}

// Init handle the event when the bot is first connected
func (s *Slacker) Init(initHandler func()) {
	s.initHandler = initHandler
}

// Err handle when errors are encountered
func (s *Slacker) Err(errorHandler func(err string)) {
	s.errorHandler = errorHandler
}

func (s *Slacker) Interactive(interactiveEventHandler func(*Slacker, *socketmode.Event, *slack.InteractionCallback)) {
	s.interactiveEventHandler = interactiveEventHandler
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
func (s *Slacker) DefaultEvent(defaultEventHandler func(interface{})) {
	s.defaultEventHandler = defaultEventHandler
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

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case evt, ok := <-s.socketModeClient.Events:
				if !ok {
					return
				}

				switch evt.Type {
				case socketmode.EventTypeConnecting:
					s.log.Println("Connecting to Slack with Socket Mode.")
					if s.initHandler == nil {
						continue
					}
					go s.initHandler()
				case socketmode.EventTypeConnectionError:
					s.log.Println("Connection failed. Retrying later...")
				case socketmode.EventTypeConnected:
					s.log.Println("Connected to Slack with Socket Mode.")
				case socketmode.EventTypeEventsAPI:
					ev, ok := evt.Data.(slackevents.EventsAPIEvent)
					if !ok {
						s.log.Printf("Ignored %+v\n", evt)
						continue
					}

					switch ev.InnerEvent.Type {
					case "message", "app_mention": // message-based events
						go s.handleMessageEvent(ctx, ev.InnerEvent.Data)

					default:
						s.log.Printf("unsupported inner event: %+v\n", ev.InnerEvent.Type)
					}

					s.socketModeClient.Ack(*evt.Request)
				case socketmode.EventTypeInteractive:
					if s.interactiveEventHandler == nil {
						s.unsupportedEventReceived()
						continue
					}

					callback, ok := evt.Data.(slack.InteractionCallback)
					if !ok {
						s.log.Printf("Ignored %+v\n", evt)
						continue
					}

					go s.interactiveEventHandler(s, &evt, &callback)
				default:
					s.unsupportedEventReceived()
				}
			}
		}
	}()

	// blocking call that handles listening for events and placing them in the
	// Events channel as well as handling outgoing events.
	return s.socketModeClient.Run()
}

func (s *Slacker) unsupportedEventReceived() {
	s.socketModeClient.Debugf("unsupported Events API event received")
}

// GetUserInfo retrieve complete user information
func (s *Slacker) GetUserInfo(user string) (*slack.User, error) {
	return s.client.GetUserInfo(user)
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

func (s *Slacker) handleMessageEvent(ctx context.Context, evt interface{}) {
	if s.botContextConstructor == nil {
		s.botContextConstructor = NewBotContext
	}

	if s.requestConstructor == nil {
		s.requestConstructor = NewRequest
	}

	if s.responseConstructor == nil {
		s.responseConstructor = NewResponse
	}

	ev := newMessageEvent(evt)
	if ev == nil {
		// event doesn't appear to be a valid message type
		return
	}

	botCtx := s.botContextConstructor(ctx, s.client, s.socketModeClient, ev)
	response := s.responseConstructor(botCtx)

	for _, cmd := range s.botCommands {
		parameters, isMatch := cmd.Match(ev.Text)
		if !isMatch {
			continue
		}

		request := s.requestConstructor(botCtx, parameters)
		if cmd.Definition().AuthorizationFunc != nil && !cmd.Definition().AuthorizationFunc(botCtx, request) {
			response.ReportError(s.unAuthorizedError)
			return
		}

		select {
		case s.commandChannel <- NewCommandEvent(cmd.Usage(), parameters, ev):
		default:
			// full channel, dropped event
		}

		cmd.Execute(botCtx, request, response)
		return
	}

	if s.defaultMessageHandler != nil {
		request := s.requestConstructor(botCtx, nil)
		s.defaultMessageHandler(botCtx, request, response)
	}
}

func newMessageEvent(evt interface{}) *MessageEvent {
	var me *MessageEvent

	switch evt.(type) {
	case *slackevents.MessageEvent:
		ev := evt.(*slackevents.MessageEvent)
		me = &MessageEvent{
			Channel:         ev.Channel,
			User:            ev.User,
			Text:            ev.Text,
			Data:            evt,
			Type:            ev.Type,
			TimeStamp:       ev.TimeStamp,
			ThreadTimeStamp: ev.ThreadTimeStamp,
			BotID:           ev.BotID,
		}
	case *slackevents.AppMentionEvent:
		ev := evt.(*slackevents.AppMentionEvent)
		me = &MessageEvent{
			Channel:         ev.Channel,
			User:            ev.User,
			Text:            ev.Text,
			Data:            evt,
			Type:            ev.Type,
			TimeStamp:       ev.TimeStamp,
			ThreadTimeStamp: ev.ThreadTimeStamp,
			BotID:           ev.BotID,
		}
	}

	// Filter out other bots. At the very least this is needed for MessageEvent
	// to prevent the bot from self-triggering and causing loops. However better
	// logic should be in place to prevent repeated self-triggering / bot-storms
	// if we want to enable this later.
	if me.IsBot() {
		return nil
	}

	return me
}
