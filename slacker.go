package slacker

import (
	"context"
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

const (
	space                = " "
	dash                 = "-"
	newLine              = "\n"
	invalidToken         = "invalid token"
	helpCommand          = "help"
	codeMessageFormat    = "`%s`"
	boldMessageFormat    = "*%s*"
	italicMessageFormat  = "_%s_"
	exampleMessageFormat = "_*Example:*_ %s"
)

// NewClient creates a new client using the Slack API
func NewClient(botToken, appToken string, clientOptions ...ClientOption) *Slacker {
	options := newClientOptions(clientOptions...)
	slackOpts := newSlackOptions(appToken, options)

	slackAPI := slack.New(botToken, slackOpts...)
	socketModeClient := socketmode.New(
		slackAPI,
		socketmode.OptionDebug(options.Debug),
	)

	slacker := &Slacker{
		slackClient:              slackAPI,
		socketModeClient:         socketModeClient,
		cronClient:               cron.New(cron.WithLocation(options.CronLocation)),
		commandGroups:            []*CommandGroup{newGroup("")},
		botInteractionMode:       options.BotMode,
		sanitizeEventTextHandler: defaultEventTextSanitizer,
		logger:                   options.Logger,
		interactions:             make(map[slack.InteractionType][]*Interaction),
	}
	return slacker
}

// Slacker contains the Slack API, botCommands, and handlers
type Slacker struct {
	slackClient                   *slack.Client
	socketModeClient              *socketmode.Client
	cronClient                    *cron.Cron
	commandMiddlewares            []CommandMiddlewareHandler
	commandGroups                 []*CommandGroup
	interactionMiddlewares        []InteractionMiddlewareHandler
	interactions                  map[slack.InteractionType][]*Interaction
	jobMiddlewares                []JobMiddlewareHandler
	jobs                          []*Job
	onHello                       func(socketmode.Event)
	onConnected                   func(socketmode.Event)
	onConnecting                  func(socketmode.Event)
	onConnectionError             func(socketmode.Event)
	onDisconnected                func(socketmode.Event)
	unsupportedInteractionHandler InteractionHandler
	helpDefinition                *CommandDefinition
	unsupportedCommandHandler     CommandHandler
	unsupportedEventHandler       func(socketmode.Event)
	appID                         string
	botInteractionMode            BotMode
	sanitizeEventTextHandler      func(string) string
	logger                        Logger
}

// GetCommandGroups returns Command Groups
func (s *Slacker) GetCommandGroups() []*CommandGroup {
	return s.commandGroups
}

// GetInteractions returns Groups
func (s *Slacker) GetInteractions() map[slack.InteractionType][]*Interaction {
	return s.interactions
}

// GetJobs returns Jobs
func (s *Slacker) GetJobs() []*Job {
	return s.jobs
}

// SlackClient returns the internal slack.Client of Slacker struct
func (s *Slacker) SlackClient() *slack.Client {
	return s.slackClient
}

// SocketModeClient returns the internal socketmode.Client of Slacker struct
func (s *Slacker) SocketModeClient() *socketmode.Client {
	return s.socketModeClient
}

// OnHello handle the event when slack sends the bot "hello"
func (s *Slacker) OnHello(onHello func(socketmode.Event)) {
	s.onHello = onHello
}

// OnConnected handle the event when the bot is connected
func (s *Slacker) OnConnected(onConnected func(socketmode.Event)) {
	s.onConnected = onConnected
}

// OnConnecting handle the event when the bot is connecting
func (s *Slacker) OnConnecting(onConnecting func(socketmode.Event)) {
	s.onConnecting = onConnecting
}

// OnConnectionError handle the event when the bot fails to connect
func (s *Slacker) OnConnectionError(onConnectionError func(socketmode.Event)) {
	s.onConnectionError = onConnectionError
}

// OnDisconnected handle the event when the bot is disconnected
func (s *Slacker) OnDisconnected(onDisconnected func(socketmode.Event)) {
	s.onDisconnected = onDisconnected
}

// UnsupportedInteractionHandler handles interactions when none of the callbacks are matched
func (s *Slacker) UnsupportedInteractionHandler(unsupportedInteractionHandler InteractionHandler) {
	s.unsupportedInteractionHandler = unsupportedInteractionHandler
}

// UnsupportedCommandHandler handles messages when none of the commands are matched
func (s *Slacker) UnsupportedCommandHandler(unsupportedCommandHandler CommandHandler) {
	s.unsupportedCommandHandler = unsupportedCommandHandler
}

// UnsupportedEventHandler handles events when an unknown event is seen
func (s *Slacker) UnsupportedEventHandler(unsupportedEventHandler func(socketmode.Event)) {
	s.unsupportedEventHandler = unsupportedEventHandler
}

// SanitizeEventTextHandler overrides the default event text sanitization
func (s *Slacker) SanitizeEventTextHandler(sanitizeEventTextHandler func(in string) string) {
	s.sanitizeEventTextHandler = sanitizeEventTextHandler
}

// Help handle the help message, it will use the default if not set
func (s *Slacker) Help(definition *CommandDefinition) {
	if len(definition.Command) == 0 {
		s.logger.Error("missing `Command`")
		return
	}
	s.helpDefinition = definition
}

// AddCommand define a new command and append it to the list of bot commands
func (s *Slacker) AddCommand(definition *CommandDefinition) {
	if len(definition.Command) == 0 {
		s.logger.Error("missing `Command`")
		return
	}
	s.commandGroups[0].AddCommand(definition)
}

// AddCommandMiddleware appends a new command middleware to the list of root level command middlewares
func (s *Slacker) AddCommandMiddleware(middleware CommandMiddlewareHandler) {
	s.commandMiddlewares = append(s.commandMiddlewares, middleware)
}

// AddCommandGroup define a new group and append it to the list of groups
func (s *Slacker) AddCommandGroup(prefix string) *CommandGroup {
	group := newGroup(prefix)
	s.commandGroups = append(s.commandGroups, group)
	return group
}

// AddInteraction define a new interaction and append it to the list of interactions
func (s *Slacker) AddInteraction(definition *InteractionDefinition) {
	if len(definition.InteractionID) == 0 {
		s.logger.Error("missing `ID`")
		return
	}
	if len(definition.Type) == 0 {
		s.logger.Error("missing `Type`")
		return
	}
	s.interactions[definition.Type] = append(s.interactions[definition.Type], newInteraction(definition))
}

// AddInteractionMiddleware appends a new interaction middleware to the list of root level interaction middlewares
func (s *Slacker) AddInteractionMiddleware(middleware InteractionMiddlewareHandler) {
	s.interactionMiddlewares = append(s.interactionMiddlewares, middleware)
}

// AddJob define a new cron job and append it to the list of jobs
func (s *Slacker) AddJob(definition *JobDefinition) {
	if len(definition.CronExpression) == 0 {
		s.logger.Error("missing `CronExpression`")
		return
	}
	s.jobs = append(s.jobs, newJob(definition))
}

// AddJobMiddleware appends a new job middleware to the list of root level job middlewares
func (s *Slacker) AddJobMiddleware(middleware JobMiddlewareHandler) {
	s.jobMiddlewares = append(s.jobMiddlewares, middleware)
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
					s.logger.Info("connecting to Slack with Socket Mode...")

					if s.onConnecting == nil {
						continue
					}
					go s.onConnecting(socketEvent)

				case socketmode.EventTypeConnectionError:
					s.logger.Info("connection failed. Retrying later...")

					if s.onConnectionError == nil {
						continue
					}
					go s.onConnectionError(socketEvent)

				case socketmode.EventTypeConnected:
					s.logger.Info("connected to Slack with Socket Mode.")

					if s.onConnected == nil {
						continue
					}
					go s.onConnected(socketEvent)

				case socketmode.EventTypeHello:
					s.appID = socketEvent.Request.ConnectionInfo.AppID
					s.logger.Info("connected as App ID %v", s.appID)

					if s.onHello == nil {
						continue
					}
					go s.onHello(socketEvent)

				case socketmode.EventTypeDisconnect:
					s.logger.Info("disconnected due to %v", socketEvent.Request.Reason)

					if s.onDisconnected == nil {
						continue
					}
					go s.onDisconnected(socketEvent)

				case socketmode.EventTypeEventsAPI:
					event, ok := socketEvent.Data.(slackevents.EventsAPIEvent)
					if !ok {
						s.logger.Debug("ignored %+v", socketEvent)
						continue
					}

					// Acknowledge receiving the request
					s.socketModeClient.Ack(*socketEvent.Request)

					if event.Type != slackevents.CallbackEvent {
						if s.unsupportedEventHandler != nil {
							s.unsupportedEventHandler(socketEvent)
						} else {
							s.logger.Debug("unsupported event received %+v", socketEvent)
						}
						continue
					}

					switch event.InnerEvent.Type {
					case "message", "app_mention": // message-based events
						go s.handleMessageEvent(ctx, event.InnerEvent.Data)

					default:
						if s.unsupportedEventHandler != nil {
							s.unsupportedEventHandler(socketEvent)
						} else {
							s.logger.Debug("unsupported event received %+v", socketEvent)
						}
					}

				case socketmode.EventTypeSlashCommand:
					event, ok := socketEvent.Data.(slack.SlashCommand)
					if !ok {
						s.logger.Debug("ignored %+v", socketEvent)
						continue
					}

					// Acknowledge receiving the request
					s.socketModeClient.Ack(*socketEvent.Request)

					go s.handleMessageEvent(ctx, &event)

				case socketmode.EventTypeInteractive:
					callback, ok := socketEvent.Data.(slack.InteractionCallback)
					if !ok {
						s.logger.Debug("ignored %+v", socketEvent)
						continue
					}

					// Acknowledge receiving the request
					s.socketModeClient.Ack(*socketEvent.Request)

					go s.handleInteractionEvent(ctx, &callback)

				default:
					if s.unsupportedEventHandler != nil {
						s.unsupportedEventHandler(socketEvent)
					} else {
						s.logger.Debug("unsupported event received %+v", socketEvent)
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

func (s *Slacker) defaultHelp(ctx *CommandContext) {
	blocks := []slack.Block{}

	for _, group := range s.GetCommandGroups() {
		for _, command := range group.GetCommands() {
			if command.Definition().HideHelp {
				continue
			}

			helpMessage := empty
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

			blocks = append(blocks,
				slack.NewSectionBlock(
					slack.NewTextBlockObject(slack.MarkdownType, helpMessage, false, false),
					nil, nil,
				))

			if len(command.Definition().Examples) > 0 {
				examplesMessage := empty
				for _, example := range command.Definition().Examples {
					examplesMessage += fmt.Sprintf(exampleMessageFormat, example) + newLine
				}

				blocks = append(blocks, slack.NewContextBlock("",
					slack.NewTextBlockObject(slack.MarkdownType, examplesMessage, false, false),
				))
			}
		}
	}

	if len(s.GetJobs()) == 0 {
		ctx.Response().ReplyBlocks(blocks)
		return
	}

	blocks = append(blocks, slack.NewDividerBlock())
	for _, job := range s.GetJobs() {
		if job.Definition().HideHelp {
			continue
		}

		helpMessage := fmt.Sprintf(codeMessageFormat, job.Definition().CronExpression)

		if len(job.Definition().Name) > 0 {
			helpMessage += space + dash + space + fmt.Sprintf(codeMessageFormat, job.Definition().Name)
		}

		if len(job.Definition().Description) > 0 {
			helpMessage += space + dash + space + fmt.Sprintf(italicMessageFormat, job.Definition().Description)
		}

		blocks = append(blocks,
			slack.NewSectionBlock(
				slack.NewTextBlockObject(slack.MarkdownType, helpMessage, false, false),
				nil, nil,
			))
	}

	ctx.Response().ReplyBlocks(blocks)
}

func (s *Slacker) prependHelpHandle() {
	if s.helpDefinition == nil {
		s.helpDefinition = &CommandDefinition{
			Command:     helpCommand,
			Description: helpCommand,
			Handler:     s.defaultHelp,
		}
	}

	s.commandGroups[0].PrependCommand(s.helpDefinition)
}

func (s *Slacker) startCronJobs(ctx context.Context) {
	middlewares := make([]JobMiddlewareHandler, 0)
	middlewares = append(middlewares, s.jobMiddlewares...)

	for _, job := range s.jobs {
		definition := job.Definition()
		middlewares = append(middlewares, definition.Middlewares...)
		jobCtx := newJobContext(ctx, s.logger, s.slackClient, definition)
		_, err := s.cronClient.AddFunc(definition.CronExpression, executeJob(jobCtx, definition.Handler, middlewares...))
		if err != nil {
			s.logger.Error(err.Error())
		}

	}

	s.cronClient.Start()
}

func (s *Slacker) handleInteractionEvent(ctx context.Context, callback *slack.InteractionCallback) {
	middlewares := make([]InteractionMiddlewareHandler, 0)
	middlewares = append(middlewares, s.interactionMiddlewares...)

	var interaction *Interaction
	var definition *InteractionDefinition

	switch callback.Type {
	case slack.InteractionTypeBlockActions:
		for _, i := range s.interactions[callback.Type] {
			for _, a := range callback.ActionCallback.BlockActions {
				definition = i.Definition()
				if a.BlockID == definition.InteractionID {
					interaction = i
					break
				}
			}
			if interaction != nil {
				break
			}
		}
	case slack.InteractionTypeViewClosed, slack.InteractionTypeViewSubmission:
		for _, i := range s.interactions[callback.Type] {
			definition = i.Definition()
			if definition.InteractionID == callback.View.CallbackID {
				interaction = i
				break
			}
		}
	case slack.InteractionTypeShortcut, slack.InteractionTypeMessageAction:
		for _, i := range s.interactions[callback.Type] {
			definition = i.Definition()
			if definition.InteractionID == callback.CallbackID {
				interaction = i
				break
			}
		}
	}

	if interaction != nil {
		interactionCtx := newInteractionContext(ctx, s.logger, s.slackClient, callback, definition)
		middlewares = append(middlewares, definition.Middlewares...)
		executeInteraction(interactionCtx, definition.Handler, middlewares...)
		return
	}

	s.logger.Debug("unsupported interaction type received %s\n", callback.Type)
	if s.unsupportedInteractionHandler != nil {
		interactionCtx := newInteractionContext(ctx, s.logger, s.slackClient, callback, nil)
		executeInteraction(interactionCtx, s.unsupportedInteractionHandler, middlewares...)
	}
}

func (s *Slacker) handleMessageEvent(ctx context.Context, event any) {
	messageEvent := newMessageEvent(s.logger, s.slackClient, event)
	if messageEvent == nil {
		// event doesn't appear to be a valid message type
		return
	}

	if messageEvent.IsBot() {
		if s.ignoreBotMessage(messageEvent) {
			return
		}
	}

	middlewares := make([]CommandMiddlewareHandler, 0)
	middlewares = append(middlewares, s.commandMiddlewares...)

	eventText := s.sanitizeEventTextHandler(messageEvent.Text)
	for _, group := range s.commandGroups {
		for _, cmd := range group.GetCommands() {
			parameters, isMatch := cmd.Match(eventText)
			if !isMatch {
				continue
			}

			definition := cmd.Definition()
			ctx := newCommandContext(ctx, s.logger, s.slackClient, messageEvent, definition, parameters)

			middlewares = append(middlewares, group.GetMiddlewares()...)
			middlewares = append(middlewares, definition.Middlewares...)
			executeCommand(ctx, definition.Handler, middlewares...)
			return
		}
	}

	if s.unsupportedCommandHandler != nil {
		ctx := newCommandContext(ctx, s.logger, s.slackClient, messageEvent, nil, nil)
		executeCommand(ctx, s.unsupportedCommandHandler, middlewares...)
	}
}

func (s *Slacker) ignoreBotMessage(messageEvent *MessageEvent) bool {
	switch s.botInteractionMode {
	case BotModeIgnoreApp:
		bot, err := s.slackClient.GetBotInfo(messageEvent.BotID)
		if err != nil {
			if err.Error() == "missing_scope" {
				s.logger.Error("unable to determine if bot response is from me -- please add users:read scope to your app")
			} else {
				s.logger.Debug("unable to get information on the bot that sent message: %v", err)
			}
			return true
		}
		if bot.AppID == s.appID {
			s.logger.Debug("ignoring event that originated from my App ID: %v", bot.AppID)
			return true
		}
	case BotModeIgnoreAll:
		s.logger.Debug("ignoring event that originated from Bot ID: %v", messageEvent.BotID)
		return true
	default:
		// BotInteractionModeIgnoreNone is handled in the default case
	}
	return false
}

func newSlackOptions(appToken string, options *clientOptions) []slack.Option {
	slackOptions := []slack.Option{
		slack.OptionDebug(options.Debug),
		slack.OptionAppLevelToken(appToken),
	}

	if len(options.APIURL) > 0 {
		slackOptions = append(slackOptions, slack.OptionAPIURL(options.APIURL))
	}
	return slackOptions
}

func defaultEventTextSanitizer(msg string) string {
	return strings.ReplaceAll(msg, "\u00a0", " ")
}
