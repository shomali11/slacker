package slacker

import (
	"fmt"

	"github.com/slack-go/slack"
)

const (
	errorFormat = "*Error:* _%s_"
)

// A ResponseWriter interface is used to respond to an event
type ResponseWriter interface {
	Reply(text string, options ...ReplyOption)
	ReportError(err error, options ...ReportErrorOption)
	Typing()
}

// NewResponse creates a new response structure
func NewResponse(botCtx BotContext) ResponseWriter {
	return &response{botCtx: botCtx}
}

type response struct {
	botCtx BotContext
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *response) ReportError(err error, options ...ReportErrorOption) {
	defaults := newReportErrorDefaults(options...)

	rtm := r.botCtx.RTM()
	event := r.botCtx.Event()
	message := rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), event.Channel)
	if defaults.ThreadResponse {
		message.ThreadTimestamp = event.EventTimestamp
	}

	rtm.SendMessage(message)
}

// Typing send a typing indicator
func (r *response) Typing() {
	rtm := r.botCtx.RTM()
	event := r.botCtx.Event()
	rtm.SendMessage(rtm.NewTypingMessage(event.Channel))
}

// Reply send a attachments to the current channel with a message
func (r *response) Reply(message string, options ...ReplyOption) {
	defaults := newReplyDefaults(options...)

	rtm := r.botCtx.RTM()
	event := r.botCtx.Event()
	if defaults.ThreadResponse {
		rtm.PostMessage(
			event.Channel,
			slack.MsgOptionText(message, false),
			slack.MsgOptionUser(rtm.GetInfo().User.ID),
			slack.MsgOptionAsUser(true),
			slack.MsgOptionAttachments(defaults.Attachments...),
			slack.MsgOptionBlocks(defaults.Blocks...),
			slack.MsgOptionTS(event.EventTimestamp),
		)
	} else {
		rtm.PostMessage(
			event.Channel,
			slack.MsgOptionText(message, false),
			slack.MsgOptionUser(rtm.GetInfo().User.ID),
			slack.MsgOptionAsUser(true),
			slack.MsgOptionAttachments(defaults.Attachments...),
			slack.MsgOptionBlocks(defaults.Blocks...),
		)
	}
}
