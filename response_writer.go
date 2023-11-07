package slacker

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
)

// newWriter creates a new poster structure
func newWriter(ctx context.Context, logger Logger, slackClient *slack.Client) *Writer {
	return &Writer{ctx: ctx, logger: logger, slackClient: slackClient}
}

// Writer sends messages to Slack
type Writer struct {
	ctx         context.Context
	logger      Logger
	slackClient *slack.Client
}

// Post send a message to a channel
func (r *Writer) Post(channel string, message string, options ...PostOption) (string, error) {
	return r.post(channel, message, []slack.Block{}, options...)
}

// PostError send an error to a channel
func (r *Writer) PostError(channel string, err error, options ...PostOption) (string, error) {
	attachments := []slack.Attachment{}
	attachments = append(attachments, slack.Attachment{
		Color: "danger",
		Text:  err.Error(),
	})
	return r.post(channel, "", []slack.Block{}, SetAttachments(attachments))
}

// PostBlocks send blocks to a channel
func (r *Writer) PostBlocks(channel string, blocks []slack.Block, options ...PostOption) (string, error) {
	return r.post(channel, "", blocks, options...)
}

// Delete deletes message
func (r *Writer) Delete(channel string, messageTimestamp string) (string, error) {
	_, timestamp, err := r.slackClient.DeleteMessage(
		channel,
		messageTimestamp,
	)
	if err != nil {
		r.logger.Error("failed to delete message", "error", err)
	}
	return timestamp, err
}

func (r *Writer) post(channel string, message string, blocks []slack.Block, options ...PostOption) (string, error) {
	postOptions := newPostOptions(options...)

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(postOptions.Attachments...),
		slack.MsgOptionBlocks(blocks...),
	}

	if len(postOptions.ThreadTS) > 0 {
		opts = append(opts, slack.MsgOptionTS(postOptions.ThreadTS))
	}

	if len(postOptions.ReplaceMessageTS) > 0 {
		opts = append(opts, slack.MsgOptionUpdate(postOptions.ReplaceMessageTS))
	}

	if len(postOptions.EphemeralUserID) > 0 {
		opts = append(opts, slack.MsgOptionPostEphemeral(postOptions.EphemeralUserID))
	}

	if postOptions.ScheduleTime != nil {
		postAt := fmt.Sprintf("%d", postOptions.ScheduleTime.Unix())
		opts = append(opts, slack.MsgOptionSchedule(postAt))
	}

	_, timestamp, err := r.slackClient.PostMessageContext(
		r.ctx,
		channel,
		opts...,
	)
	if err != nil {
		r.logger.Error("failed to post message", "error", err)
	}
	return timestamp, err
}
