package slacker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/robfig/cron"
	"github.com/shomali11/proper"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

const (
	space               = " "
	dash                = "-"
	newLine             = "\n"
	lock                = ":lock:"
	invalidToken        = "invalid token"
	helpCommand         = "help"
	directChannelMarker = "D"
	userMentionFormat   = "<@%s>"
	codeMessageFormat   = "`%s`"
	boldMessageFormat   = "*%s*"
	italicMessageFormat = "_%s_"
	quoteMessageFormat  = ">_*Example:* %s_"
	slackBotUser        = "USLACKBOT"
)

var (
	errUnauthorized = errors.New("you are not authorized to execute this command")
)

func defaultCleanEventInput(msg string) string {
	return strings.ReplaceAll(msg, "\u00a0", " ")
}

// fontFormattingCleanEventInput cleans up the message from slack to remove font formatting
func fontFormattingCleanEventInput(msg string) string {
	msg = defaultCleanEventInput(msg)
	re := regexp.MustCompile("[*_~]")
	return re.ReplaceAllString(msg, "")
}

// allFormattingCleanEventInput cleans up the message from slack to remove font formatting and code blocks
func allFormattingCleanEventInput(msg string) string {
	msg = defaultCleanEventInput(msg)
	re := regexp.MustCompile("[*_~`>]")
	return re.ReplaceAllString(msg, "")
}

// NewClient creates a new client using the Slack API
func NewClient(botToken, appToken string, options ...ClientOption) *Slacker {
	defaults := newClientDefaults(options...)

	slackOpts := []slack.Option{
		slack.OptionDebug(defaults.Debug),
		slack.OptionAppLevelToken(appToken),
	}

	if defaults.APIURL != "" {
		slackOpts = append(slackOpts, slack.OptionAPIURL(defaults.APIURL))
	}

	// set sanitize event
	sanitizeEvent := defaultCleanEventInput
	if defaults.MessageWithoutFontFormatting {
		sanitizeEvent = fontFormattingCleanEventInput
	}
	if defaults.MessageWithoutAllFormatting {
		sanitizeEvent = allFormattingCleanEventInput
	}

	api := slack.New(
		botToken,
		slackOpts...,
	)

	socketModeClient := socketmode.New(
		api,
		socketmode.OptionDebug(defaults.Debug),
	)

	slacker := &Slacker{
		apiClient:          api,
		socketModeClient:   socketModeClient,
		cronClient:         cron.New(),
		commandChannel:     make(chan *CommandEvent, 100),
		errUnauthorized:    errUnauthorized,
		botInteractionMode: defaults.BotMode,
		sanitizeEventText:  sanitizeEvent,
		debug:              defaults.Debug,
	}
	return slacker
}

// Slacker contains the Slack API, botCommands, and handlers
type Slacker struct {
	apiClient                        *slack.Client
	socketModeClient                 *socketmode.Client
	cronClient                       *cron.Cron
	commands                         []Command
	botContextConstructor            func(context.Context, *slack.Client, *socketmode.Client, *MessageEvent) BotContext
	interactiveBotContextConstructor func(context.Context, *slack.Client, *socketmode.Client, *socketmode.Event) InteractiveBotContext
	commandConstructor               func(string, *CommandDefinition) Command
	requestConstructor               func(BotContext, *proper.Properties) Request
	responseConstructor              func(BotContext) ResponseWriter
	jobs                             []Job
	jobContextConstructor            func(context.Context, *slack.Client, *socketmode.Client) JobContext
	jobConstructor                   func(string, *JobDefinition) Job
	initHandler                      func()
	errorHandler                     func(err string)
	interactiveEventHandler          func(InteractiveBotContext, *slack.InteractionCallback)
	helpDefinition                   *CommandDefinition
	defaultMessageHandler            func(BotContext, Request, ResponseWriter)
	defaultEventHandler              func(interface{})
	defaultInnerEventHandler         func(context.Context, interface{}, *socketmode.Request)
	errUnauthorized                  error
	commandChannel                   chan *CommandEvent
	appID                            string
	botInteractionMode               BotInteractionMode
	sanitizeEventText                func(string) string
	debug                            bool
}

// BotCommands returns Bot Commands
func (s *Slacker) BotCommands() []Command {
	return s.commands
}

// APIClient returns the internal slack.Client of Slacker struct
func (s *Slacker) APIClient() *slack.Client {
	return s.apiClient
}

// SocketModeClient returns the internal socketmode.Client of Slacker struct
func (s *Slacker) SocketModeClient() *socketmode.Client {
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

// SanitizeEventText allows the api consumer to override the default event text sanitization
func (s *Slacker) SanitizeEventText(sanitizeEventText func(in string) string) {
	s.sanitizeEventText = sanitizeEventText
}

// Interactive assigns an interactive event handler
func (s *Slacker) Interactive(interactiveEventHandler func(InteractiveBotContext, *slack.InteractionCallback)) {
	s.interactiveEventHandler = interactiveEventHandler
}

// CustomBotContext creates a new bot context
func (s *Slacker) CustomBotContext(botContextConstructor func(context.Context, *slack.Client, *socketmode.Client, *MessageEvent) BotContext) {
	s.botContextConstructor = botContextConstructor
}

// CustomInteractiveBotContext creates a new interactive bot context
func (s *Slacker) CustomInteractiveBotContext(interactiveBotContextConstructor func(context.Context, *slack.Client, *socketmode.Client, *socketmode.Event) InteractiveBotContext) {
	s.interactiveBotContextConstructor = interactiveBotContextConstructor
}

// CustomJobContext creates a new job context
func (s *Slacker) CustomJobContext(jobContextConstructor func(context.Context, *slack.Client, *socketmode.Client) JobContext) {
	s.jobContextConstructor = jobContextConstructor
}

// CustomCommand creates a new BotCommand
func (s *Slacker) CustomCommand(commandConstructor func(usage string, definition *CommandDefinition) Command) {
	s.commandConstructor = commandConstructor
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

// DefaultInnerEvent handle events when an unknown inner event is seen
func (s *Slacker) DefaultInnerEvent(defaultInnerEventHandler func(ctx context.Context, evt interface{}, request *socketmode.Request)) {
	s.defaultInnerEventHandler = defaultInnerEventHandler
}

// UnAuthorizedError error message
func (s *Slacker) UnAuthorizedError(errUnauthorized error) {
	s.errUnauthorized = errUnauthorized
}

// Help handle the help message, it will use the default if not set
func (s *Slacker) Help(definition *CommandDefinition) {
	s.helpDefinition = definition
}

// Command define a new command and append it to the list of existing bot commands
func (s *Slacker) Command(usage string, definition *CommandDefinition) {
	if s.commandConstructor == nil {
		s.commandConstructor = NewCommand
	}
	s.commands = append(s.commands, s.commandConstructor(usage, definition))
}

// Job define a new cron job and append it to the list of existing jobs
func (s *Slacker) Job(spec string, definition *JobDefinition) {
	if s.jobConstructor == nil {
		s.jobConstructor = NewJob
	}
	s.jobs = append(s.jobs, s.jobConstructor(spec, definition))
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
			case socketEvent, ok := <-s.socketModeClient.Events:
				if !ok {
					return
				}

				switch socketEvent.Type {
				case socketmode.EventTypeConnecting:
					s.logf("Connecting to Slack with Socket Mode.")
					if s.initHandler == nil {
						continue
					}
					go s.initHandler()

				case socketmode.EventTypeConnectionError:
					s.logf("Connection failed. Retrying later...")

				case socketmode.EventTypeConnected:
					s.logf("Connected to Slack with Socket Mode.")

				case socketmode.EventTypeHello:
					s.appID = socketEvent.Request.ConnectionInfo.AppID
					s.logf("Connected as App ID %v\n", s.appID)

				case socketmode.EventTypeEventsAPI:
					event, ok := socketEvent.Data.(slackevents.EventsAPIEvent)
					if !ok {
						s.debugf("Ignored %+v\n", socketEvent)
						continue
					}

					switch event.InnerEvent.Type {
					case "message", "app_mention": // message-based events
						go s.handleMessageEvent(ctx, event.InnerEvent.Data, nil)

					default:
						if s.defaultInnerEventHandler != nil {
							s.defaultInnerEventHandler(ctx, event.InnerEvent.Data, socketEvent.Request)
						} else {
							s.debugf("unsupported inner event: %+v\n", event.InnerEvent.Type)
						}
					}

					s.socketModeClient.Ack(*socketEvent.Request)

				case socketmode.EventTypeSlashCommand:
					callback, ok := socketEvent.Data.(slack.SlashCommand)
					if !ok {
						s.debugf("Ignored %+v\n", socketEvent)
						continue
					}
					s.socketModeClient.Ack(*socketEvent.Request)
					go s.handleMessageEvent(ctx, &callback, socketEvent.Request)

				case socketmode.EventTypeInteractive:
					callback, ok := socketEvent.Data.(slack.InteractionCallback)
					if !ok {
						s.debugf("Ignored %+v\n", socketEvent)
						continue
					}

					go s.handleInteractiveEvent(ctx, &socketEvent, &callback)

				default:
					if s.defaultEventHandler != nil {
						s.defaultEventHandler(socketEvent)
					} else {
						s.unsupportedEventReceived()
					}
				}
			}
		}
	}()

	s.startCronJobs(ctx)
	defer s.cronClient.Stop()

	// blocking call that handles listening for events and placing them in the
	// Events channel as well as handling outgoing events.
	return s.socketModeClient.RunContext(ctx)
}

func (s *Slacker) unsupportedEventReceived() {
	s.socketModeClient.Debugf("unsupported Events API event received")
}

func (s *Slacker) defaultHelp(botCtx BotContext, request Request, response ResponseWriter) {
	helpMessage := empty
	for _, command := range s.commands {
		if command.Definition().HideHelp {
			continue
		}
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
			helpMessage += space + lock
		}

		helpMessage += newLine

		for _, example := range command.Definition().Examples {
			helpMessage += fmt.Sprintf(quoteMessageFormat, example) + newLine
		}
	}

	for _, command := range s.jobs {
		if command.Definition().HideHelp {
			continue
		}

		helpMessage += fmt.Sprintf(codeMessageFormat, command.Spec()) + space

		if len(command.Definition().Description) > 0 {
			helpMessage += dash + space + fmt.Sprintf(italicMessageFormat, command.Definition().Description)
		}

		helpMessage += newLine
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

	s.commands = append([]Command{NewCommand(helpCommand, s.helpDefinition)}, s.commands...)
}

func (s *Slacker) startCronJobs(ctx context.Context) {
	if s.jobContextConstructor == nil {
		s.jobContextConstructor = NewJobContext
	}

	jobCtx := s.jobContextConstructor(ctx, s.apiClient, s.socketModeClient)
	for _, jobCommand := range s.jobs {
		s.cronClient.AddFunc(jobCommand.Spec(), jobCommand.Callback(jobCtx))
	}

	s.cronClient.Start()
}

func (s *Slacker) handleInteractiveEvent(ctx context.Context, event *socketmode.Event, callback *slack.InteractionCallback) {
	if s.interactiveBotContextConstructor == nil {
		s.interactiveBotContextConstructor = NewInteractiveBotContext
	}

	botCtx := s.interactiveBotContextConstructor(ctx, s.apiClient, s.socketModeClient, event)
	for _, cmd := range s.commands {
		for _, action := range callback.ActionCallback.BlockActions {
			if action.BlockID != cmd.Definition().BlockID {
				continue
			}

			cmd.Interactive(botCtx, event.Request, callback)
			return
		}
	}

	if s.interactiveEventHandler != nil {
		s.interactiveEventHandler(botCtx, callback)
	}
}

func (s *Slacker) handleMessageEvent(ctx context.Context, event interface{}, req *socketmode.Request) {
	if s.botContextConstructor == nil {
		s.botContextConstructor = NewBotContext
	}

	if s.requestConstructor == nil {
		s.requestConstructor = NewRequest
	}

	if s.responseConstructor == nil {
		s.responseConstructor = NewResponse
	}

	messageEvent := NewMessageEvent(s, event, req)
	if messageEvent == nil {
		// event doesn't appear to be a valid message type
		return
	} else if messageEvent.IsBot() {
		switch s.botInteractionMode {
		case BotInteractionModeIgnoreApp:
			bot, err := s.apiClient.GetBotInfo(messageEvent.BotID)
			if err != nil {
				if err.Error() == "missing_scope" {
					s.logf("unable to determine if bot response is from me -- please add users:read scope to your app\n")
				} else {
					s.debugf("unable to get bot that sent message information: %v\n", err)
				}
				return
			}
			if bot.AppID == s.appID {
				s.debugf("Ignoring event that originated from my App ID: %v\n", bot.AppID)
				return
			}
		case BotInteractionModeIgnoreAll:
			s.debugf("Ignoring event that originated from Bot ID: %v\n", messageEvent.BotID)
			return
		default:
			// BotInteractionModeIgnoreNone is handled in the default case
		}
	}

	botCtx := s.botContextConstructor(ctx, s.apiClient, s.socketModeClient, messageEvent)
	response := s.responseConstructor(botCtx)

	eventText := s.sanitizeEventText(messageEvent.Text)
	for _, cmd := range s.commands {
		parameters, isMatch := cmd.Match(eventText)
		if !isMatch {
			continue
		}

		request := s.requestConstructor(botCtx, parameters)
		if cmd.Definition().AuthorizationFunc != nil && !cmd.Definition().AuthorizationFunc(botCtx, request) {
			response.ReportError(s.errUnauthorized)
			return
		}

		select {
		case s.commandChannel <- NewCommandEvent(cmd.Usage(), parameters, messageEvent):
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

func (s *Slacker) logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (s *Slacker) debugf(format string, v ...interface{}) {
	if s.debug {
		log.Printf(format, v...)
	}
}
