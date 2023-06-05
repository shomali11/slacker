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
	directChannelMarker  = "D"
	userMentionFormat    = "<@%s>"
	codeMessageFormat    = "`%s`"
	boldMessageFormat    = "*%s*"
	italicMessageFormat  = "_%s_"
	exampleMessageFormat = "_*Example:*_ %s"
	slackBotUser         = "USLACKBOT"
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
		slackClient:          slackAPI,
		socketModeClient:   socketModeClient,
		cronClient:         cron.New(),
		commandGroups:      []CommandGroup{newGroup("")},
		botInteractionMode: options.BotMode,
		sanitizeEventText:  defaultEventTextSanitizer,
		logger:             options.Logger,
	}
	return slacker
}

// Slacker contains the Slack API, botCommands, and handlers
type Slacker struct {
	slackClient                   *slack.Client
	socketModeClient            *socketmode.Client
	cronClient                  *cron.Cron
	commandMiddlewares          []CommandMiddlewareHandler
	commandGroups               []CommandGroup
	interactionMiddlewares      []InteractionMiddlewareHandler
	interactions                []Interaction
	jobMiddlewares              []JobMiddlewareHandler
	jobs                        []Job
	initHandler                 func()
	unhandledInteractionHandler InteractionHandler
	helpDefinition              *CommandDefinition
	unhandledMessageHandler     CommandHandler
	unhandledEventHandler       func(socketmode.Event)
	appID                       string
	botInteractionMode          BotInteractionMode
	sanitizeEventText           func(string) string
	logger                      Logger
}

// GetGroups returns Groups
func (s *Slacker) GetGroups() []CommandGroup {
	return s.commandGroups
}

// SlackClient returns the internal slack.Client of Slacker struct
func (s *Slacker) SlackClient() *slack.Client {
	return s.slackClient
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

// UnhandledInteractionHandler assigns an interaction handler when none of the interactions are handled
func (s *Slacker) UnhandledInteractionHandler(unhanldedInteractionHandler InteractionHandler) {
	s.unhandledInteractionHandler = unhanldedInteractionHandler
}

// UnhandledMessageHandler handle messages when none of the commands are matched
func (s *Slacker) UnhandledMessageHandler(unhandledMessageHandler CommandHandler) {
	s.unhandledMessageHandler = unhandledMessageHandler
}

// UnhandledEventHandler handle events when an unknown event is seen
func (s *Slacker) UnhandledEventHandler(unhandledEventHandler func(socketmode.Event)) {
	s.unhandledEventHandler = unhandledEventHandler
}

// Help handle the help message, it will use the default if not set
func (s *Slacker) Help(definition *CommandDefinition) {
	s.helpDefinition = definition
}

// AddCommand define a new command and append it to the list of bot commands
func (s *Slacker) AddCommand(usage string, definition *CommandDefinition) {
	s.commandGroups[0].AddCommand(usage, definition)
}

// AddCommandMiddleware appends a new command middleware to the list of root level command middlewares
func (s *Slacker) AddCommandMiddleware(middleware CommandMiddlewareHandler) {
	s.commandMiddlewares = append(s.commandMiddlewares, middleware)
}

// AddCommandGroup define a new group and append it to the list of groups
func (s *Slacker) AddCommandGroup(prefix string) CommandGroup {
	group := newGroup(prefix)
	s.commandGroups = append(s.commandGroups, group)
	return group
}

// AddInteraction define a new interaction and append it to the list of interactions
func (s *Slacker) AddInteraction(blockID string, definition *InteractionDefinition) {
	s.interactions = append(s.interactions, newInteraction(blockID, definition))
}

// AddInteractionMiddleware appends a new interaction middleware to the list of root level interaction middlewares
func (s *Slacker) AddInteractionMiddleware(middleware InteractionMiddlewareHandler) {
	s.interactionMiddlewares = append(s.interactionMiddlewares, middleware)
}

// AddJob define a new cron job and append it to the list of jobs
func (s *Slacker) AddJob(spec string, definition *JobDefinition) {
	s.jobs = append(s.jobs, newJob(spec, definition))
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
					s.logger.Infof("connecting to Slack with Socket Mode...\n")

				case socketmode.EventTypeConnectionError:
					s.logger.Infof("connection failed. Retrying later...\n")

				case socketmode.EventTypeConnected:
					s.logger.Infof("connected to Slack with Socket Mode.\n")

					if s.initHandler == nil {
						continue
					}
					go s.initHandler()

				case socketmode.EventTypeHello:
					s.appID = socketEvent.Request.ConnectionInfo.AppID
					s.logger.Infof("connected as App ID %v\n", s.appID)

				case socketmode.EventTypeDisconnect:
					s.logger.Infof("disconnected due to %v\n", socketEvent.Request.Reason)

				case socketmode.EventTypeEventsAPI:
					event, ok := socketEvent.Data.(slackevents.EventsAPIEvent)
					if !ok {
						s.logger.Debugf("ignored %+v\n", socketEvent)
						continue
					}

					// Acknowledge receiving the request
					s.socketModeClient.Ack(*socketEvent.Request)

					switch event.InnerEvent.Type {
					case "message", "app_mention": // message-based events
						go s.handleMessageEvent(ctx, event.InnerEvent.Data)

					default:
						if s.unhandledEventHandler != nil {
							s.unhandledEventHandler(socketEvent)
						} else {
							s.logger.Debugf("unsupported event received %+v\n", socketEvent)
						}
					}

				case socketmode.EventTypeSlashCommand:
					event, ok := socketEvent.Data.(slack.SlashCommand)
					if !ok {
						s.logger.Debugf("ignored %+v\n", socketEvent)
						continue
					}

					// Acknowledge receiving the request
					s.socketModeClient.Ack(*socketEvent.Request)

					go s.handleMessageEvent(ctx, &event)

				case socketmode.EventTypeInteractive:
					callback, ok := socketEvent.Data.(slack.InteractionCallback)
					if !ok {
						s.logger.Debugf("ignored %+v\n", socketEvent)
						continue
					}

					// Acknowledge receiving the request
					s.socketModeClient.Ack(*socketEvent.Request)

					go s.handleInteractionEvent(ctx, &callback)

				default:
					if s.unhandledEventHandler != nil {
						s.unhandledEventHandler(socketEvent)
					} else {
						s.logger.Debugf("unsupported event received %+v\n", socketEvent)
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

func (s *Slacker) defaultHelp(ctx CommandContext) {
	blocks := []slack.Block{}

	for _, group := range s.GetGroups() {
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

	ctx.Response().ReplyBlocks(blocks)
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

	s.commandGroups[0].PrependCommand(helpCommand, s.helpDefinition)
}

func (s *Slacker) startCronJobs(ctx context.Context) {
	jobCtx := newJobContext(ctx, s.logger, s.slackClient)

	middlewares := make([]JobMiddlewareHandler, 0)
	middlewares = append(middlewares, s.jobMiddlewares...)

	for _, job := range s.jobs {
		definition := job.Definition()
		_, err := s.cronClient.AddFunc(definition.Spec, executeJob(jobCtx, definition.Handler, middlewares...))
		if err != nil {
			s.logger.Errorf(err.Error())
		}

	}

	s.cronClient.Start()
}

func (s *Slacker) handleInteractionEvent(ctx context.Context, callback *slack.InteractionCallback) {
	interactionCtx := newInteractionContext(ctx, s.logger, s.slackClient, callback)

	middlewares := make([]InteractionMiddlewareHandler, 0)
	middlewares = append(middlewares, s.interactionMiddlewares...)

	for _, interaction := range s.interactions {
		for _, action := range callback.ActionCallback.BlockActions {
			definition := interaction.Definition()
			if action.BlockID != definition.BlockID {
				continue
			}

			middlewares = append(middlewares, definition.Middlewares...)
			executeInteraction(interactionCtx, definition.Handler, middlewares...)
			return
		}
	}

	if s.unhandledInteractionHandler != nil {
		executeInteraction(interactionCtx, s.unhandledInteractionHandler, middlewares...)
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

	eventText := s.sanitizeEventText(messageEvent.Text)
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

	if s.unhandledMessageHandler != nil {
		ctx := newCommandContext(ctx, s.logger, s.slackClient, messageEvent, nil, nil)
		executeCommand(ctx, s.unhandledMessageHandler, middlewares...)
	}
}

func (s *Slacker) ignoreBotMessage(messageEvent *MessageEvent) bool {
	switch s.botInteractionMode {
	case BotInteractionModeIgnoreApp:
		bot, err := s.slackClient.GetBotInfo(messageEvent.BotID)
		if err != nil {
			if err.Error() == "missing_scope" {
				s.logger.Errorf("unable to determine if bot response is from me -- please add users:read scope to your app\n")
			} else {
				s.logger.Debugf("unable to get information on the bot that sent message: %v\n", err)
			}
			return true
		}
		if bot.AppID == s.appID {
			s.logger.Debugf("ignoring event that originated from my App ID: %v\n", bot.AppID)
			return true
		}
	case BotInteractionModeIgnoreAll:
		s.logger.Debugf("ignoring event that originated from Bot ID: %v\n", messageEvent.BotID)
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
