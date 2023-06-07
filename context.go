package slacker

import (
	"context"

	"github.com/shomali11/proper"
	"github.com/slack-go/slack"
)

// CommandContext interface is for bot command contexts
type CommandContext interface {
	Context() context.Context
	Definition() *CommandDefinition
	Event() *MessageEvent
	Request() Request
	Response() WriterReplierResponse
	SlackClient() *slack.Client
}

// newCommandContext creates a new bot context
func newCommandContext(
	ctx context.Context,
	logger Logger,
	slackClient *slack.Client,
	event *MessageEvent,
	definition *CommandDefinition,
	parameters *proper.Properties,
) CommandContext {
	request := newRequest(parameters)
	writer := newWriter(ctx, logger, slackClient)
	replier := newReplier(event.ChannelID, event.UserID, event.TimeStamp, writer)
	response := newWriterReplierResponse(writer, replier)

	return &commandContext{
		ctx:         ctx,
		event:       event,
		slackClient: slackClient,
		definition:  definition,
		request:     request,
		response:    response,
	}
}

type commandContext struct {
	ctx         context.Context
	event       *MessageEvent
	slackClient *slack.Client
	definition  *CommandDefinition
	request     Request
	response    WriterReplierResponse
}

// Context returns the context
func (r *commandContext) Context() context.Context {
	return r.ctx
}

// Definition returns the command definition
func (r *commandContext) Definition() *CommandDefinition {
	return r.definition
}

// Event returns the slack message event
func (r *commandContext) Event() *MessageEvent {
	return r.event
}

// SlackClient returns the slack API client
func (r *commandContext) SlackClient() *slack.Client {
	return r.slackClient
}

// Request returns the command request
func (r *commandContext) Request() Request {
	return r.request
}

// Response returns the response writer
func (r *commandContext) Response() WriterReplierResponse {
	return r.response
}

// InteractionContext interface is interaction bot contexts
type InteractionContext interface {
	Context() context.Context
	Definition() *InteractionDefinition
	Callback() *slack.InteractionCallback
	Response() WriterReplierResponse
	SlackClient() *slack.Client
}

// newInteractionContext creates a new interaction bot context
func newInteractionContext(
	ctx context.Context,
	logger Logger,
	slackClient *slack.Client,
	callback *slack.InteractionCallback,
	definition *InteractionDefinition,
) InteractionContext {
	writer := newWriter(ctx, logger, slackClient)
	replier := newReplier(callback.Channel.ID, callback.User.ID, callback.MessageTs, writer)
	response := newWriterReplierResponse(writer, replier)
	return &interactionContext{
		ctx:         ctx,
		definition:  definition,
		callback:    callback,
		slackClient: slackClient,
		response:    response,
	}
}

type interactionContext struct {
	ctx         context.Context
	definition  *InteractionDefinition
	callback    *slack.InteractionCallback
	slackClient *slack.Client
	response    WriterReplierResponse
}

// Context returns the context
func (r *interactionContext) Context() context.Context {
	return r.ctx
}

// Definition returns the interaction definition
func (r *interactionContext) Definition() *InteractionDefinition {
	return r.definition
}

// Callback returns the interaction callback
func (r *interactionContext) Callback() *slack.InteractionCallback {
	return r.callback
}

// Response returns the response writer
func (r *interactionContext) Response() WriterReplierResponse {
	return r.response
}

// SlackClient returns the slack API client
func (r *interactionContext) SlackClient() *slack.Client {
	return r.slackClient
}

// JobContext interface is for job command contexts
type JobContext interface {
	Context() context.Context
	Definition() *JobDefinition
	Response() WriterResponse
	SlackClient() *slack.Client
}

// newJobContext creates a new bot context
func newJobContext(ctx context.Context, logger Logger, slackClient *slack.Client, definition *JobDefinition) JobContext {
	writer := newWriter(ctx, logger, slackClient)
	response := newWriterResponse(writer)
	return &jobContext{
		ctx:         ctx,
		definition:  definition,
		slackClient: slackClient,
		response:    response,
	}
}

type jobContext struct {
	ctx         context.Context
	definition  *JobDefinition
	slackClient *slack.Client
	response    WriterResponse
}

// Context returns the context
func (r *jobContext) Context() context.Context {
	return r.ctx
}

// Definition returns the job definition
func (r *jobContext) Definition() *JobDefinition {
	return r.definition
}

// Response returns the response writer
func (r *jobContext) Response() WriterResponse {
	return r.response
}

// SlackClient returns the slack API client
func (r *jobContext) SlackClient() *slack.Client {
	return r.slackClient
}
