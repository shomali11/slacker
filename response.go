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
func (r *response) ReportError(err error, options ...ReportErrorOption) {
	defaults := newReportErrorDefaults(options...)

	message := r.rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), r.event.Channel)
	if defaults.ThreadResponse {
		message.ThreadTimestamp = r.event.EventTimestamp
	}

	r.rtm.SendMessage(message)
}

// Typing send a typing indicator
func (r *response) Typing() {
	r.rtm.SendMessage(r.rtm.NewTypingMessage(r.event.Channel))
}

// Reply send a attachments to the current channel with a message
func (r *response) Reply(message string, options ...ReplyOption) {
	defaults := newReplyDefaults(options...)

	if defaults.ThreadResponse {
		r.rtm.PostMessage(
			r.event.Channel,
			slack.MsgOptionText(message, false),
			slack.MsgOptionUser(r.rtm.GetInfo().User.ID),
			slack.MsgOptionAsUser(true),
			slack.MsgOptionAttachments(defaults.Attachments...),
			slack.MsgOptionBlocks(defaults.Blocks...),
			slack.MsgOptionTS(r.event.EventTimestamp),
		)
	} else {
		r.rtm.PostMessage(
			r.event.Channel,
			slack.MsgOptionText(message, false),
			slack.MsgOptionUser(r.rtm.GetInfo().User.ID),
			slack.MsgOptionAsUser(true),
			slack.MsgOptionAttachments(defaults.Attachments...),
			slack.MsgOptionBlocks(defaults.Blocks...),
		)
	}
}

// RTM returns the RTM client
func (r *response) RTM() *slack.RTM {
	return r.rtm
}

// Client returns the slack client
func (r *response) Client() *slack.Client {
	return r.client
}
