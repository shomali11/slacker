package slacker

import (
	"context"
	"fmt"
	"strings"

	"github.com/robfig/cron"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
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
	quoteMessageFormat  = ">_*Example:* %s_"
	slackBotUser        = "USLACKBOT"
)

func defaultCleanEventInput(msg string) string {
	return strings.ReplaceAll(msg, "\u00a0", " ")
}

// NewClient creates a new client using the Slack API
func NewClient(botToken, appToken string, clientOptions ...ClientOption) *Slacker {
	options := newClientOptions(clientOptions...)

	slackOpts := []slack.Option{
		slack.OptionDebug(options.Debug),
		slack.OptionAppLevelToken(appToken),
	}

	if options.APIURL != "" {
		slackOpts = append(slackOpts, slack.OptionAPIURL(options.APIURL))
	}

	api := slack.New(
		botToken,
		slackOpts...,
	)

	socketModeClient := socketmode.New(
		api,
		socketmode.OptionDebug(options.Debug),
	)

	setLogDebugMode(options.Debug)

	slacker := &Slacker{
		apiClient:          api,
		socketModeClient:   socketModeClient,
		cronClient:         cron.New(),
		groups:             []Group{newGroup("")},
		botInteractionMode: options.BotMode,
		sanitizeEventText:  defaultCleanEventInput,
		debug:              options.Debug,
	}
	return slacker
}

// Slacker contains the Slack API, botCommands, and handlers
type Slacker struct {
	apiClient                    *slack.Client
	socketModeClient             *socketmode.Client
	cronClient                   *cron.Cron
	middlewares                  []MiddlewareHandler
	groups                       []Group
	jobs                         []Job
	initHandler                  func()
	unhandledInteractiveCallback InteractiveHandler
	helpDefinition               *CommandDefinition
	defaultMessageHandler        CommandHandler
	defaultEventHandler          func(socketmode.Event)
	unhandledInnerEventHandler   func(context.Context, interface{}, *socketmode.Request)
	appID                        string
	botInteractionMode           BotInteractionMode
	sanitizeEventText            func(string) string
	debug                        bool
}

// GetGroups returns Groups
func (s *Slacker) GetGroups() []Group {
	return s.groups
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

// SanitizeEventText allows the api consumer to override the default event text sanitization
func (s *Slacker) SanitizeEventText(sanitizeEventText func(in string) string) {
	s.sanitizeEventText = sanitizeEventText
}

// Interactive assigns an interactive event handler
func (s *Slacker) Interactive(interactiveEventHandler InteractiveHandler) {
	s.unhandledInteractiveCallback = interactiveEventHandler
}

// DefaultCommand handle messages when none of the commands are matched
func (s *Slacker) DefaultCommand(defaultMessageHandler CommandHandler) {
	s.defaultMessageHandler = defaultMessageHandler
}

// DefaultEvent handle events when an unknown event is seen
func (s *Slacker) DefaultEvent(defaultEventHandler func(socketmode.Event)) {
	s.defaultEventHandler = defaultEventHandler
}

// DefaultInnerEvent handle events when an unknown inner event is seen
func (s *Slacker) DefaultInnerEvent(defaultInnerEventHandler func(ctx context.Context, evt interface{}, request *socketmode.Request)) {
	s.unhandledInnerEventHandler = defaultInnerEventHandler
}

// Help handle the help message, it will use the default if not set
func (s *Slacker) Help(definition *CommandDefinition) {
	s.helpDefinition = definition
}

// AddCommand define a new command and append it to the list of bot commands
func (s *Slacker) AddCommand(usage string, definition *CommandDefinition) {
	s.groups[0].AddCommand(usage, definition)
}

// AddMiddleware define a new middleware and append it to the list of root level middlewares
func (s *Slacker) AddMiddleware(middleware MiddlewareHandler) {
	s.middlewares = append(s.middlewares, middleware)
}

// AddGroup define a new group and append it to the list of groups
func (s *Slacker) AddGroup(prefix string) Group {
	group := newGroup(prefix)
	s.groups = append(s.groups, group)
	return group
}

// AddJob define a new cron job and append it to the list of jobs
func (s *Slacker) AddJob(spec string, definition *JobDefinition) {
	s.jobs = append(s.jobs, newJob(spec, definition))
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
					infof("Connecting to Slack with Socket Mode.")

				case socketmode.EventTypeConnectionError:
					infof("Connection failed. Retrying later...")

				case socketmode.EventTypeConnected:
					infof("Connected to Slack with Socket Mode.")

					if s.initHandler == nil {
						continue
					}
					go s.initHandler()

				case socketmode.EventTypeHello:
					s.appID = socketEvent.Request.ConnectionInfo.AppID
					infof("Connected as App ID %v\n", s.appID)

				case socketmode.EventTypeEventsAPI:
					event, ok := socketEvent.Data.(slackevents.EventsAPIEvent)
					if !ok {
						debugf("Ignored %+v\n", socketEvent)
						continue
					}

					switch event.InnerEvent.Type {
					case "message", "app_mention": // message-based events
						go s.handleMessageEvent(ctx, event.InnerEvent.Data, nil)

					default:
						if s.unhandledInnerEventHandler != nil {
							s.unhandledInnerEventHandler(ctx, event.InnerEvent.Data, socketEvent.Request)
						} else {
							debugf("unsupported inner event: %+v\n", event.InnerEvent.Type)
						}
					}

					s.socketModeClient.Ack(*socketEvent.Request)

				case socketmode.EventTypeSlashCommand:
					callback, ok := socketEvent.Data.(slack.SlashCommand)
					if !ok {
						debugf("Ignored %+v\n", socketEvent)
						continue
					}
					s.socketModeClient.Ack(*socketEvent.Request)
					go s.handleMessageEvent(ctx, &callback, socketEvent.Request)

				case socketmode.EventTypeInteractive:
					callback, ok := socketEvent.Data.(slack.InteractionCallback)
					if !ok {
						debugf("Ignored %+v\n", socketEvent)
						continue
					}

					go s.handleInteractiveEvent(ctx, &socketEvent, &callback)

				default:
					if s.defaultEventHandler != nil {
						s.defaultEventHandler(socketEvent)
					} else {
						debugf("unsupported Events API event received")
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

func (s *Slacker) defaultHelp(botCtx CommandContext) {
	helpMessage := empty

	for _, group := range s.groups {
		for _, command := range group.GetCommands() {
			if command.Definition().HideHelp {
				continue
			}

			helpMessage += "•" + space

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

			helpMessage += newLine

			for _, example := range command.Definition().Examples {
				helpMessage += fmt.Sprintf(quoteMessageFormat, example) + newLine
			}
		}

		helpMessage += newLine
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

	botCtx.Response().Reply(helpMessage)
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

	s.groups[0].PrependCommand(helpCommand, s.helpDefinition)
}

func (s *Slacker) startCronJobs(ctx context.Context) {
	jobCtx := newJobContext(ctx, s.apiClient, s.socketModeClient)
	for _, jobCommand := range s.jobs {
		s.cronClient.AddFunc(jobCommand.Spec(), jobCommand.Callback(jobCtx))
	}

	s.cronClient.Start()
}

func (s *Slacker) handleInteractiveEvent(ctx context.Context, event *socketmode.Event, callback *slack.InteractionCallback) {
	botCtx := newInteractiveContext(ctx, s.apiClient, s.socketModeClient, event, callback)

	for _, group := range s.groups {
		for _, cmd := range group.GetCommands() {
			for _, action := range callback.ActionCallback.BlockActions {
				if action.BlockID != cmd.Definition().BlockID {
					continue
				}

				cmd.InteractiveCallback(botCtx)
				return
			}
		}
	}

	if s.unhandledInteractiveCallback != nil {
		s.unhandledInteractiveCallback(botCtx)
	}
}

func (s *Slacker) handleMessageEvent(ctx context.Context, event interface{}, request *socketmode.Request) {
	messageEvent := newMessageEvent(s.apiClient, event, request)
	if messageEvent == nil {
		// event doesn't appear to be a valid message type
		return
	}

	if messageEvent.IsBot() {
		if s.ignoreBotMessage(messageEvent) {
			return
		}
	}

	eventText := s.sanitizeEventText(messageEvent.Text)
	for _, group := range s.groups {
		for _, cmd := range group.GetCommands() {
			parameters, isMatch := cmd.Match(eventText)
			if !isMatch {
				continue
			}

			botCtx := newCommandContext(ctx, s.apiClient, s.socketModeClient, messageEvent, cmd.Usage(), parameters)
			middlewares := make([]MiddlewareHandler, 0)
			middlewares = append(middlewares, s.middlewares...)
			middlewares = append(middlewares, group.GetMiddlewares()...)
			middlewares = append(middlewares, cmd.Definition().Middlewares...)
			cmd.Handler(botCtx, middlewares...)
			return
		}
	}

	if s.defaultMessageHandler != nil {
		botCtx := newCommandContext(ctx, s.apiClient, s.socketModeClient, messageEvent, "", nil)
		s.defaultMessageHandler(botCtx)
	}
}

func (s *Slacker) ignoreBotMessage(messageEvent *MessageEvent) bool {
	switch s.botInteractionMode {
	case BotInteractionModeIgnoreApp:
		bot, err := s.apiClient.GetBotInfo(messageEvent.BotID)
		if err != nil {
			if err.Error() == "missing_scope" {
				infof("unable to determine if bot response is from me -- please add users:read scope to your app\n")
			} else {
				debugf("unable to get information on the bot that sent message: %v\n", err)
			}
			return true
		}
		if bot.AppID == s.appID {
			debugf("Ignoring event that originated from my App ID: %v\n", bot.AppID)
			return true
		}
	case BotInteractionModeIgnoreAll:
		debugf("Ignoring event that originated from Bot ID: %v\n", messageEvent.BotID)
		return true
	default:
		// BotInteractionModeIgnoreNone is handled in the default case
	}
	return false
}
