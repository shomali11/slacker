package slacker

import (
	"github.com/nlopes/slack"
)

// NewResponse creates a new response structure
func NewResponse(channel string, rtm *slack.RTM) *Response {
	return &Response{channel: channel, rtm: rtm}
}

// Response contains the channel and Real Time Messaging library
type Response struct {
	channel string
	rtm     *slack.RTM
}

// Reply send a message back to the channel where we received an event from
func (r *Response) Reply(text string) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(text, r.channel))
}

// Typing send a typing indicator
func (r *Response) Typing() {
	r.rtm.SendMessage(r.rtm.NewTypingMessage(r.channel))
}
