package slacker

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// Writer interface is used to respond to an event
type Writer interface {
	Post(channel string, message string, options ...PostOption) (string, error)
	PostError(channel string, err error, options ...PostOption) (string, error)
	PostBlocks(channel string, blocks []slack.Block, options ...PostOption) (string, error)

	Delete(channel string, messageTimestamp string) (string, error)
}

// newWriter creates a new poster structure
func newWriter(apiClient *slack.Client, socketModeClient *socketmode.Client) Writer {
	return &writer{apiClient: apiClient, socketModeClient: socketModeClient}
}

type writer struct {
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
}

// Post send a message to a channel
func (r *writer) Post(channel string, message string, options ...PostOption) (string, error) {
	return r.post(channel, message, []slack.Block{}, options...)
}

// PostError send an error to a channel
func (r *writer) PostError(channel string, err error, options ...PostOption) (string, error) {
	blocks := []slack.Block{
		slack.NewContextBlock("", slack.NewTextBlockObject(slack.MarkdownType, err.Error(), false, false)),
	}
	return r.PostBlocks(channel, blocks, options...)
}

// PostBlocks send blocks to a channel
func (r *writer) PostBlocks(channel string, blocks []slack.Block, options ...PostOption) (string, error) {
	return r.post(channel, "", blocks, options...)
}

// Delete deletes message
func (r *writer) Delete(channel string, messageTimestamp string) (string, error) {
	_, timestamp, err := r.apiClient.DeleteMessage(
		channel,
		messageTimestamp,
	)
	if err != nil {
		infof("failed to delete message: %v\n", err)
	}
	return timestamp, err
}

func (r *writer) post(channel string, message string, blocks []slack.Block, options ...PostOption) (string, error) {
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

	_, timestamp, err := r.apiClient.PostMessage(
		channel,
		opts...,
	)
	if err != nil {
		infof("failed to post message: %v\n", err)
	}
	return timestamp, err
}
