package slacker

import (
	"fmt"

	"github.com/nlopes/slack"
)

const (
	errorFormat = "*Error:* _%s_"
)

// A ResponseWriter interface is used to respond to an event
type ResponseWriter interface {
	Reply(text string)
	PostMessage(channel, message string, attachments *[]slack.Attachment)
	ReportError(err error)
	Typing()
	RTM() *slack.RTM
	Client() *slack.Client
}

// NewResponse creates a new response structure
func NewResponse(channel string, Client *slack.Client, RTM *slack.RTM) *Response {
	return &Response{channel: channel, client: Client, rtm: RTM}
}

// Response contains the channel and Real Time Messaging library
type Response struct {
	channel string
	client  *slack.Client
	rtm     *slack.RTM
}

// Reply send a message back to the channel where we received the event from
func (r *Response) Reply(text string) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(text, r.channel))
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *Response) ReportError(err error) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), r.channel))
}

// Typing send a typing indicator
func (r *Response) Typing() {
	r.rtm.SendMessage(r.rtm.NewTypingMessage(r.channel))
}

// PostMessage send a attachments to the current channel with a message
func (r *Response) PostMessage(channel, message string, attachments *[]slack.Attachment) {
	params := slack.PostMessageParameters{}
	params.User = r.rtm.GetInfo().User.ID
	params.AsUser = true
	params.Attachments = *attachments
	r.rtm.PostMessage(channel, message, params)
}

// RTM returns the RTM client
func (r *Response) RTM() *slack.RTM {
	return r.rtm
}

// Client returns the slack client
func (r *Response) Client() *slack.Client {
	return r.client
}
