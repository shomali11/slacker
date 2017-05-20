package slacker

import (
	"github.com/nlopes/slack"
)

func NewResponse(channel string, rtm *slack.RTM) *Response {
	return &Response{channel: channel, rtm: rtm}
}

type Response struct {
	channel string
	rtm     *slack.RTM
}

func (r *Response) Reply(text string) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(text, r.channel))
}

func (r *Response) Typing() {
	r.rtm.SendMessage(r.rtm.NewTypingMessage(r.channel))
}
