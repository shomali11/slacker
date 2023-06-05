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
	APIClient() *slack.Client
}

// newCommandContext creates a new bot context
func newCommandContext(
	ctx context.Context,
	logger Logger,
	apiClient *slack.Client,
	event *MessageEvent,
	definition *CommandDefinition,
	parameters *proper.Properties,
) CommandContext {
	request := newRequest(parameters)
	poster := newWriter(logger, apiClient)
	replier := newReplier(event.ChannelID, event.UserID, event.TimeStamp, poster)
	response := newWriterReplierResponse(poster, replier)

	return &commandContext{
		ctx:        ctx,
		event:      event,
		apiClient:  apiClient,
		definition: definition,
		request:    request,
		response:   response,
	}
}

type commandContext struct {
	ctx        context.Context
	event      *MessageEvent
	apiClient  *slack.Client
	definition *CommandDefinition
	request    Request
	response   WriterReplierResponse
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

// APIClient returns the slack API client
func (r *commandContext) APIClient() *slack.Client {
	return r.apiClient
}

// Request returns the command request
func (r *commandContext) Request() Request {
	return r.request
}

// Response returns the command response
func (r *commandContext) Response() WriterReplierResponse {
	return r.response
}

// InteractionContext interface is interaction bot contexts
type InteractionContext interface {
	Context() context.Context
	Callback() *slack.InteractionCallback
	Response() WriterReplierResponse
	APIClient() *slack.Client
}

// newInteractionContext creates a new interaction bot context
func newInteractionContext(
	ctx context.Context,
	logger Logger,
	apiClient *slack.Client,
	callback *slack.InteractionCallback,
) InteractionContext {
	poster := newWriter(logger, apiClient)
	replier := newReplier(callback.Channel.ID, callback.User.ID, callback.MessageTs, poster)
	response := newWriterReplierResponse(poster, replier)
	return &interactionContext{
		ctx:       ctx,
		callback:  callback,
		apiClient: apiClient,
		response:  response,
	}
}

type interactionContext struct {
	ctx       context.Context
	callback  *slack.InteractionCallback
	apiClient *slack.Client
	response  WriterReplierResponse
}

// Context returns the context
func (r *interactionContext) Context() context.Context {
	return r.ctx
}

// Callback returns the interaction callback
func (r *interactionContext) Callback() *slack.InteractionCallback {
	return r.callback
}

// Response returns the command response
func (r *interactionContext) Response() WriterReplierResponse {
	return r.response
}

// APIClient returns the slack API client
func (r *interactionContext) APIClient() *slack.Client {
	return r.apiClient
}

// JobContext interface is for job command contexts
type JobContext interface {
	Context() context.Context
	Response() WriterResponse
	APIClient() *slack.Client
}

// newJobContext creates a new bot context
func newJobContext(ctx context.Context, logger Logger, apiClient *slack.Client) JobContext {
	poster := newWriter(logger, apiClient)
	response := newWriterResponse(poster)
	return &jobContext{
		ctx:       ctx,
		apiClient: apiClient,
		response:  response,
	}
}

type jobContext struct {
	ctx       context.Context
	apiClient *slack.Client
	response  WriterResponse
}

// Context returns the context
func (r *jobContext) Context() context.Context {
	return r.ctx
}

// Response returns the command response
func (r *jobContext) Response() WriterResponse {
	return r.response
}

// APIClient returns the slack API client
func (r *jobContext) APIClient() *slack.Client {
	return r.apiClient
}
