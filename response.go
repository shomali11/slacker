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
	Reply(text string, options ...ReplyOption) error
	ReportError(err error, options ...ReportErrorOption)
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
	defaults := NewReportErrorDefaults(options...)

	client := r.botCtx.Client()
	ev := r.botCtx.Event()

	opts := []slack.MsgOption{
		slack.MsgOptionText(fmt.Sprintf(errorFormat, err.Error()), false),
	}
	if defaults.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(ev.TimeStamp))
	}
	_, _, err = client.PostMessage(ev.Channel, opts...)
	if err != nil {
		fmt.Printf("failed posting message: %v\n", err)
	}
}

// Reply send a attachments to the current channel with a message
func (r *response) Reply(message string, options ...ReplyOption) error {
	defaults := NewReplyDefaults(options...)

	client := r.botCtx.Client()
	ev := r.botCtx.Event()
	if ev == nil {
		return fmt.Errorf("unable to get message event details")
	}

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(defaults.Attachments...),
		slack.MsgOptionBlocks(defaults.Blocks...),
	}
	if defaults.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(ev.TimeStamp))
	}

	_, _, err := client.PostMessage(
		ev.Channel,
		opts...,
	)
	return err
}
