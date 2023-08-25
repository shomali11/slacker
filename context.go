package slacker

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// BotContext interface is for bot command contexts
type BotContext interface {
	Context() context.Context
	Event() *MessageEvent
	APIClient() *slack.Client
	SocketModeClient() *socketmode.Client
	GetUserName() (string, error)
}

// NewBotContext creates a new bot context
func NewBotContext(ctx context.Context, apiClient *slack.Client, socketModeClient *socketmode.Client, event *MessageEvent) BotContext {
	return &botContext{ctx: ctx, event: event, apiClient: apiClient, socketModeClient: socketModeClient}
}

type botContext struct {
	ctx              context.Context
	event            *MessageEvent
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
}

// Context returns the context
func (r *botContext) Context() context.Context {
	return r.ctx
}

// Event returns the slack message event
func (r *botContext) Event() *MessageEvent {
	return r.event
}

func (r *botContext) GetUserName() (string, error) {
	u, err := r.apiClient.GetUserInfo(r.event.UserID)
	if err != nil {
		return "", err
	}
	return u.Name, nil
}

// APIClient returns the slack API client
func (r *botContext) APIClient() *slack.Client {
	return r.apiClient
}

// SocketModeClient returns the slack socket mode client
func (r *botContext) SocketModeClient() *socketmode.Client {
	return r.socketModeClient
}

// InteractiveBotContext interface is interactive bot command contexts
type InteractiveBotContext interface {
	Context() context.Context
	Event() *socketmode.Event
	APIClient() *slack.Client
	SocketModeClient() *socketmode.Client
}

// NewInteractiveBotContext creates a new interactive bot context
func NewInteractiveBotContext(ctx context.Context, apiClient *slack.Client, socketModeClient *socketmode.Client, event *socketmode.Event) InteractiveBotContext {
	return &interactiveBotContext{ctx: ctx, event: event, apiClient: apiClient, socketModeClient: socketModeClient}
}

type interactiveBotContext struct {
	ctx              context.Context
	event            *socketmode.Event
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
}

// Context returns the context
func (r *interactiveBotContext) Context() context.Context {
	return r.ctx
}

// Event returns the socket event
func (r *interactiveBotContext) Event() *socketmode.Event {
	return r.event
}

// APIClient returns the slack API client
func (r *interactiveBotContext) APIClient() *slack.Client {
	return r.apiClient
}

// SocketModeClient returns the slack socket mode client
func (r *interactiveBotContext) SocketModeClient() *socketmode.Client {
	return r.socketModeClient
}

// JobContext interface is for job command contexts
type JobContext interface {
	Context() context.Context
	APIClient() *slack.Client
	SocketModeClient() *socketmode.Client
}

// NewJobContext creates a new bot context
func NewJobContext(ctx context.Context, apiClient *slack.Client, socketModeClient *socketmode.Client) JobContext {
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
