package slacker

import (
	"context"

	"github.com/shomali11/proper"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// CommandContext interface is for bot command contexts
type CommandContext interface {
	Context() context.Context
	Definition() *CommandDefinition
	Event() *MessageEvent
	Request() Request
	Response() PosterReplierResponse
	APIClient() *slack.Client
	SocketModeClient() *socketmode.Client
}

// newCommandContext creates a new bot context
func newCommandContext(
	ctx context.Context,
	apiClient *slack.Client,
	socketModeClient *socketmode.Client,
	event *MessageEvent,
	definition *CommandDefinition,
	parameters *proper.Properties,
) CommandContext {
	request := newRequest(parameters)
	poster := newPoster(apiClient, socketModeClient)
	replier := newReplier(event.ChannelID, event.TimeStamp, poster)
	response := newPosterReplierResponse(poster, replier)

	return &commandContext{
		ctx:              ctx,
		event:            event,
		apiClient:        apiClient,
		socketModeClient: socketModeClient,
		definition:       definition,
		request:          request,
		response:         response,
	}
}

type commandContext struct {
	ctx              context.Context
	event            *MessageEvent
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
	definition       *CommandDefinition
	request          Request
	response         PosterReplierResponse
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

// SocketModeClient returns the slack socket mode client
func (r *commandContext) SocketModeClient() *socketmode.Client {
	return r.socketModeClient
}

// Request returns the command request
func (r *commandContext) Request() Request {
	return r.request
}

// Response returns the command response
func (r *commandContext) Response() PosterReplierResponse {
	return r.response
}

// InteractionContext interface is interaction bot contexts
type InteractionContext interface {
	Context() context.Context
	Event() *socketmode.Event
	Callback() *slack.InteractionCallback
	Response() PosterReplierResponse
	APIClient() *slack.Client
	SocketModeClient() *socketmode.Client
}

// newInteractionContext creates a new interaction bot context
func newInteractionContext(
	ctx context.Context,
	apiClient *slack.Client,
	socketModeClient *socketmode.Client,
	event *socketmode.Event,
	callback *slack.InteractionCallback,
) InteractionContext {
	poster := newPoster(apiClient, socketModeClient)
	replier := newReplier(callback.Channel.ID, callback.MessageTs, poster)
	response := newPosterReplierResponse(poster, replier)
	return &interactionContext{
		ctx:              ctx,
		event:            event,
		apiClient:        apiClient,
		socketModeClient: socketModeClient,
		callback:         callback,
		response:         response,
	}
}

type interactionContext struct {
	ctx              context.Context
	event            *socketmode.Event
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
	callback         *slack.InteractionCallback
	response         PosterReplierResponse
}

// Context returns the context
func (r *interactionContext) Context() context.Context {
	return r.ctx
}

// Event returns the socket event
func (r *interactionContext) Event() *socketmode.Event {
	return r.event
}

// Response returns the command response
func (r *interactionContext) Response() PosterReplierResponse {
	return r.response
}

// APIClient returns the slack API client
func (r *interactionContext) APIClient() *slack.Client {
	return r.apiClient
}

// SocketModeClient returns the slack socket mode client
func (r *interactionContext) SocketModeClient() *socketmode.Client {
	return r.socketModeClient
}

// Callback returns the command callback
func (r *interactionContext) Callback() *slack.InteractionCallback {
	return r.callback
}

// JobContext interface is for job command contexts
type JobContext interface {
	Context() context.Context
	Response() PosterResponse
	APIClient() *slack.Client
	SocketModeClient() *socketmode.Client
}

// newJobContext creates a new bot context
func newJobContext(ctx context.Context, apiClient *slack.Client, socketModeClient *socketmode.Client) JobContext {
	poster := newPoster(apiClient, socketModeClient)
	response := newPosterResponse(poster)
	return &jobContext{
		ctx:              ctx,
		apiClient:        apiClient,
		socketModeClient: socketModeClient,
		response:         response,
	}
}

type jobContext struct {
	ctx              context.Context
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
	response         PosterResponse
}

// Context returns the context
func (r *jobContext) Context() context.Context {
	return r.ctx
}

// Response returns the command response
func (r *jobContext) Response() PosterResponse {
	return r.response
}

// APIClient returns the slack API client
func (r *jobContext) APIClient() *slack.Client {
	return r.apiClient
}

// SocketModeClient returns the slack socket mode client
func (r *jobContext) SocketModeClient() *socketmode.Client {
	return r.socketModeClient
}
