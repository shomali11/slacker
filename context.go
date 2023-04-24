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
	Usage() string
	Event() *MessageEvent
	Request() Request
	Response() ResponseWriter
	APIClient() *slack.Client
	SocketModeClient() *socketmode.Client
}

// newCommandContext creates a new bot context
func newCommandContext(
	ctx context.Context,
	apiClient *slack.Client,
	socketModeClient *socketmode.Client,
	event *MessageEvent,
	usage string,
	parameters *proper.Properties,
) CommandContext {
	return &commandContext{
		ctx:              ctx,
		event:            event,
		apiClient:        apiClient,
		socketModeClient: socketModeClient,
		usage:            usage,
		request:          newRequest(parameters),
		response:         newResponse(event, apiClient, socketModeClient),
	}
}

type commandContext struct {
	ctx              context.Context
	event            *MessageEvent
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
	usage            string
	request          Request
	response         ResponseWriter
}

// Context returns the context
func (r *commandContext) Context() context.Context {
	return r.ctx
}

// Usage returns the command usage
func (r *commandContext) Usage() string {
	return r.usage
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

// Response returns the command response writer
func (r *commandContext) Response() ResponseWriter {
	return r.response
}

// InteractiveContext interface is interactive bot contexts
type InteractiveContext interface {
	Context() context.Context
	Event() *socketmode.Event
	Callback() *slack.InteractionCallback
	APIClient() *slack.Client
	SocketModeClient() *socketmode.Client
}

// newInteractiveContext creates a new interactive bot context
func newInteractiveContext(
	ctx context.Context,
	apiClient *slack.Client,
	socketModeClient *socketmode.Client,
	event *socketmode.Event,
	callback *slack.InteractionCallback,
) InteractiveContext {
	return &interactiveContext{
		ctx:              ctx,
		event:            event,
		apiClient:        apiClient,
		socketModeClient: socketModeClient,
		callback:         callback,
	}
}

type interactiveContext struct {
	ctx              context.Context
	event            *socketmode.Event
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
	request          *socketmode.Request
	callback         *slack.InteractionCallback
}

// Context returns the context
func (r *interactiveContext) Context() context.Context {
	return r.ctx
}

// Event returns the socket event
func (r *interactiveContext) Event() *socketmode.Event {
	return r.event
}

// APIClient returns the slack API client
func (r *interactiveContext) APIClient() *slack.Client {
	return r.apiClient
}

// SocketModeClient returns the slack socket mode client
func (r *interactiveContext) SocketModeClient() *socketmode.Client {
	return r.socketModeClient
}

// Callback returns the command callback
func (r *interactiveContext) Callback() *slack.InteractionCallback {
	return r.callback
}

// JobContext interface is for job command contexts
type JobContext interface {
	Context() context.Context
	APIClient() *slack.Client
	SocketModeClient() *socketmode.Client
}

// newJobContext creates a new bot context
func newJobContext(ctx context.Context, apiClient *slack.Client, socketModeClient *socketmode.Client) JobContext {
	return &jobContext{ctx: ctx, apiClient: apiClient, socketModeClient: socketModeClient}
}

type jobContext struct {
	ctx              context.Context
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
}

// Context returns the context
func (r *jobContext) Context() context.Context {
	return r.ctx
}

// APIClient returns the slack API client
func (r *jobContext) APIClient() *slack.Client {
	return r.apiClient
}

// SocketModeClient returns the slack socket mode client
func (r *jobContext) SocketModeClient() *socketmode.Client {
	return r.socketModeClient
}
