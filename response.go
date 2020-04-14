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
	ReplyInThread(text string, options ...ReplyOption)
	ReportError(err error)
	ReportErrorInThread(err error)
	Typing()
	TypingInThread()
	RTM() *slack.RTM
	Client() *slack.Client
}

// NewResponse creates a new response structure
func NewResponse(event *slack.MessageEvent, client *slack.Client, rtm *slack.RTM) ResponseWriter {
	return &response{event: event, client: client, rtm: rtm}
}

type response struct {
	event  *slack.MessageEvent
	client *slack.Client
	rtm    *slack.RTM
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *response) ReportError(err error) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), r.event.Channel))
}

// ReportErrorInThread sends back a formatted error message to the channel where we received the event from inside a thread to the previous message
func (r *response) ReportErrorInThread(err error) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), r.event.Channel, slack.RTMsgOptionTS(r.event.EventTimestamp)))

}

// Typing send a typing indicator
func (r *response) Typing() {
	r.rtm.SendMessage(r.rtm.NewTypingMessage(r.event.Channel))
}

// TypingInThread send a typing indicator in a thread
func (r *response) TypingInThread() {
	message := r.rtm.NewTypingMessage(r.event.Channel)
	message.ThreadTimestamp = r.event.EventTimestamp
	r.rtm.SendMessage(message)
}

// Reply send a attachments to the current channel with a message
func (r *response) Reply(message string, options ...ReplyOption) {
	defaults := newReplyDefaults(options...)

	r.rtm.PostMessage(
		r.event.Channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionUser(r.rtm.GetInfo().User.ID),
		slack.MsgOptionAsUser(true),
		slack.MsgOptionAttachments(defaults.Attachments...),
		slack.MsgOptionBlocks(defaults.Blocks...),
	)
}

// ReplyInThread send a attachments to the current channel with a message in a thread to the previous message
func (r *response) ReplyInThread(message string, options ...ReplyOption) {
	defaults := newReplyDefaults(options...)

	r.rtm.PostMessage(
		r.event.Channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionUser(r.rtm.GetInfo().User.ID),
		slack.MsgOptionAsUser(true),
		slack.MsgOptionAttachments(defaults.Attachments...),
		slack.MsgOptionBlocks(defaults.Blocks...),
		slack.MsgOptionTS(r.event.EventTimestamp),
	)
}

// RTM returns the RTM client
func (r *response) RTM() *slack.RTM {
	return r.rtm
}

// Client returns the slack client
func (r *response) Client() *slack.Client {
	return r.client
}
