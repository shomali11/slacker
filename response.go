package slacker

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// A ResponseWriter interface is used to respond to an event
type ResponseWriter interface {
	Post(channel string, message string, options ...ReplyOption)
	Reply(text string, options ...ReplyOption)
	Error(err error, options ...ErrorOption)
}

// newResponse creates a new response structure
func newResponse(event *MessageEvent, apiClient *slack.Client, socketModeClient *socketmode.Client) ResponseWriter {
	return &response{event: event, apiClient: apiClient, socketModeClient: socketModeClient}
}

type response struct {
	event            *MessageEvent
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
}

// Error sends back a formatted error message to the channel where we received the event from
func (r *response) Error(err error, options ...ErrorOption) {
	errorOptions := newErrorOptions(options...)

	opts := []slack.MsgOption{
		slack.MsgOptionText(fmt.Sprintf(errorOptions.Format, err.Error()), false),
	}

	if errorOptions.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(r.event.TimeStamp))
	}

	_, _, err = r.apiClient.PostMessage(r.event.ChannelID, opts...)
	if err != nil {
		infof("failed to post message: %v\n", err)
	}
}

// Reply send a message to the current channel
func (r *response) Reply(message string, options ...ReplyOption) {
	if r.event == nil {
		infof("unable to get message event details\n")
		return
	}
	r.Post(r.event.ChannelID, message, options...)
}

// Post send a message to a channel
func (r *response) Post(channel string, message string, options ...ReplyOption) {
	replyOptions := newReplyOptions(options...)

	if r.event == nil {
		infof("unable to get message event details\n")
		return
	}

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(replyOptions.Attachments...),
		slack.MsgOptionBlocks(replyOptions.Blocks...),
	}

	if replyOptions.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(r.event.TimeStamp))
	}

	_, _, err := r.apiClient.PostMessage(
		channel,
		opts...,
	)
	if err != nil {
		infof("failed to post message: %v\n", err)
	}
}
