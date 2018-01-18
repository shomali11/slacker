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
	ReportError(err error)
	Typing()
}

// NewResponse creates a new response structure
func NewResponse(channel string, RTM *slack.RTM) *Response {
	return &Response{channel: channel, RTM: RTM}
}

// Response contains the channel and Real Time Messaging library
type Response struct {
	channel string
	RTM     *slack.RTM
}

// Reply send a message back to the channel where we received the event from
func (r *Response) Reply(text string) {
	r.RTM.SendMessage(r.RTM.NewOutgoingMessage(text, r.channel))
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *Response) ReportError(err error) {
	r.RTM.SendMessage(r.RTM.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), r.channel))
}

// Typing send a typing indicator
func (r *Response) Typing() {
	r.RTM.SendMessage(r.RTM.NewTypingMessage(r.channel))
}
