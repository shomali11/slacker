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
	Post(channel string, message string, options ...ReplyOption) error
	Reply(text string, options ...ReplyOption) error
	ReplyWithMention(text string, options ...ReplyOption) error
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

	apiClient := r.botCtx.APIClient()
	event := r.botCtx.Event()

	opts := []slack.MsgOption{
		slack.MsgOptionText(fmt.Sprintf(errorFormat, err.Error()), false),
	}

	if defaults.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(event.TimeStamp))
	}

	_, _, err = apiClient.PostMessage(event.ChannelID, opts...)
	if err != nil {
		fmt.Printf("failed posting message: %v\n", err)
	}
}

// Reply send a message to the current channel
func (r *response) Reply(message string, options ...ReplyOption) error {
	ev := r.botCtx.Event()
	if ev == nil {
		return fmt.Errorf("unable to get message event details")
	}
	return r.Post(ev.ChannelID, message, options...)
}

func (r *response) ReplyWithMention(message string, options ...ReplyOption) error {
	ev := r.botCtx.Event()
	if ev == nil {
		return fmt.Errorf("unable to get message event details")
	}
	mentionMessage := "<@" + r.botCtx.Event().UserID + ">" + message
	return r.Post(ev.ChannelID, mentionMessage, options...)
}

// Post send a message to a channel
func (r *response) Post(channel string, message string, options ...ReplyOption) error {
	defaults := NewReplyDefaults(options...)

	apiClient := r.botCtx.APIClient()
	event := r.botCtx.Event()
	if event == nil {
		return fmt.Errorf("unable to get message event details")
	}

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(defaults.Attachments...),
		slack.MsgOptionBlocks(defaults.Blocks...),
	}

	if defaults.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(event.TimeStamp))
	}

	_, _, err := apiClient.PostMessage(
		channel,
		opts...,
	)
	return err
}
