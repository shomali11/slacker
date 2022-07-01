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

// MessageEvent contains details common to message based events, including the
// raw event as returned from Slack along with the corresponding event type.
// The struct should be kept minimal and only include data that is commonly
// used to prevent freqeuent type assertions when evaluating the event.
type MessageEvent struct {
	// Channel ID where the message was sent
	Channel string

	// ChannelName where the message was sent
	ChannelName string

	// User ID of the sender
	User string

	// UserName of the the sender
	UserName string

	// Text is the unalterted text of the message, as returned by Slack
	Text string

	// TimeStamp is the message timestamp
	TimeStamp string

	// ThreadTimeStamp is the message thread timestamp.
	ThreadTimeStamp string

	// Data is the raw event data returned from slack. Using Type, you can assert
	// this into a slackevents *Event struct.
	Data interface{}

	// Type is the type of the event, as returned by Slack. For instance,
	// `app_mention` or `message`
	Type string

	// BotID of the bot that sent this message. If a bot did not send this
	// message, this will be an empty string.
	BotID string
}

// IsThread indicates if a message event took place in a thread.
func (e *MessageEvent) IsThread() bool {
	if e.ThreadTimeStamp == "" || e.ThreadTimeStamp == e.TimeStamp {
		return false
	}
	return true
}

// IsBot indicates if the message was sent by a bot
func (e *MessageEvent) IsBot() bool {
	return e.BotID != ""
}
