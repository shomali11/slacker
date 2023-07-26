package slacker

import (
	"context"

	"github.com/shomali11/proper"
	"github.com/slack-go/slack"
)

// newCommandContext creates a new command context
func newCommandContext(
	ctx context.Context,
	logger Logger,
	slackClient *slack.Client,
	event *MessageEvent,
	definition *CommandDefinition,
	parameters *proper.Properties,
) *CommandContext {
	request := newRequest(parameters)
	writer := newWriter(ctx, logger, slackClient)
	replier := newReplier(event.ChannelID, event.UserID, event.InThread(), event.TimeStamp, writer)
	response := newResponseReplier(writer, replier)

	return &CommandContext{
		ctx:         ctx,
		event:       event,
		slackClient: slackClient,
		definition:  definition,
		request:     request,
		response:    response,
		logger:      logger,
	}
}

// CommandContext contains information relevant to the executed command
type CommandContext struct {
	ctx         context.Context
	event       *MessageEvent
	slackClient *slack.Client
	definition  *CommandDefinition
	request     *Request
	response    *ResponseReplier
	logger      Logger
}

// Context returns the context
func (r *CommandContext) Context() context.Context {
	return r.ctx
}

// Definition returns the command definition
func (r *CommandContext) Definition() *CommandDefinition {
	return r.definition
}

// Event returns the slack message event
func (r *CommandContext) Event() *MessageEvent {
	return r.event
}

// SlackClient returns the slack API client
func (r *CommandContext) SlackClient() *slack.Client {
	return r.slackClient
}

// Request returns the command request
func (r *CommandContext) Request() *Request {
	return r.request
}

// Response returns the response writer
func (r *CommandContext) Response() *ResponseReplier {
	return r.response
}

// Logger returns the logger
func (r *CommandContext) Logger() Logger {
	return r.logger
}

// newInteractionContext creates a new interaction context
func newInteractionContext(
	ctx context.Context,
	logger Logger,
	slackClient *slack.Client,
	callback *slack.InteractionCallback,
	definition *InteractionDefinition,
) *InteractionContext {
	inThread := isMessageInThread(callback.OriginalMessage.ThreadTimestamp, callback.OriginalMessage.Timestamp)
	writer := newWriter(ctx, logger, slackClient)
	replier := newReplier(callback.Channel.ID, callback.User.ID, inThread, callback.MessageTs, writer)
	response := newResponseReplier(writer, replier)
	return &InteractionContext{
		ctx:         ctx,
		definition:  definition,
		callback:    callback,
		slackClient: slackClient,
		response:    response,
		logger:      logger,
	}
}

// InteractionContext contains information relevant to the executed interaction
type InteractionContext struct {
	ctx         context.Context
	definition  *InteractionDefinition
	callback    *slack.InteractionCallback
	slackClient *slack.Client
	response    *ResponseReplier
	logger      Logger
}

// Context returns the context
func (r *InteractionContext) Context() context.Context {
	return r.ctx
}

// Definition returns the interaction definition
func (r *InteractionContext) Definition() *InteractionDefinition {
	return r.definition
}

// Callback returns the interaction callback
func (r *InteractionContext) Callback() *slack.InteractionCallback {
	return r.callback
}

// Response returns the response writer
func (r *InteractionContext) Response() *ResponseReplier {
	return r.response
}

// SlackClient returns the slack API client
func (r *InteractionContext) SlackClient() *slack.Client {
	return r.slackClient
}

// Logger returns the logger
func (r *InteractionContext) Logger() Logger {
	return r.logger
}

// newJobContext creates a new bot context
func newJobContext(ctx context.Context, logger Logger, slackClient *slack.Client, definition *JobDefinition) *JobContext {
	writer := newWriter(ctx, logger, slackClient)
	response := newWriterResponse(writer)
	return &JobContext{
		ctx:         ctx,
		definition:  definition,
		slackClient: slackClient,
		response:    response,
		logger:      logger,
	}
}

// JobContext contains information relevant to the executed job
type JobContext struct {
	ctx         context.Context
	definition  *JobDefinition
	slackClient *slack.Client
	response    *ResponseWriter
	logger      Logger
}

// Context returns the context
func (r *JobContext) Context() context.Context {
	return r.ctx
}

// Definition returns the job definition
func (r *JobContext) Definition() *JobDefinition {
	return r.definition
}

// Response returns the response writer
func (r *JobContext) Response() *ResponseWriter {
	return r.response
}

// SlackClient returns the slack API client
func (r *JobContext) SlackClient() *slack.Client {
	return r.slackClient
}

// Logger returns the logger
func (r *JobContext) Logger() Logger {
	return r.logger
}
