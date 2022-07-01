package slacker

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// A BotContext interface is used to respond to an event
type BotContext interface {
	Context() context.Context
	Event() *MessageEvent
	SocketMode() *socketmode.Client
	Client() *slack.Client
}

// NewBotContext creates a new bot context
func NewBotContext(ctx context.Context, client *slack.Client, socketmode *socketmode.Client, evt *MessageEvent) BotContext {
	return &botContext{ctx: ctx, event: evt, client: client, socketmode: socketmode}
}

type botContext struct {
	ctx        context.Context
	event      *MessageEvent
	client     *slack.Client
	socketmode *socketmode.Client
}

// Context returns the context
func (r *botContext) Context() context.Context {
	return r.ctx
}

// Event returns the slack message event
func (r *botContext) Event() *MessageEvent {
	return r.event
}

// SocketMode returns the SocketMode client
func (r *botContext) SocketMode() *socketmode.Client {
	return r.socketmode
}

// Client returns the slack client
func (r *botContext) Client() *slack.Client {
	return r.client
}
