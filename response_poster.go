package slacker

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// Poster interface is used to respond to an event
type Poster interface {
	Post(channel string, message string, options ...PostOption)
	PostError(channel string, err error, options ...PostOption)
	PostBlocks(channel string, blocks []slack.Block, options ...PostOption)
}

// newPoster creates a new poster structure
func newPoster(apiClient *slack.Client, socketModeClient *socketmode.Client) Poster {
	return &poster{apiClient: apiClient, socketModeClient: socketModeClient}
}

type poster struct {
	apiClient        *slack.Client
	socketModeClient *socketmode.Client
}

// Post send a message to a channel
func (r *poster) Post(channel string, message string, options ...PostOption) {
	r.post(channel, message, []slack.Block{}, options...)
}

// PostError send an error to a channel
func (r *poster) PostError(channel string, err error, options ...PostOption) {
	blocks := []slack.Block{
		slack.NewContextBlock("", slack.NewTextBlockObject(slack.MarkdownType, err.Error(), false, false)),
	}
	r.PostBlocks(channel, blocks, options...)
}

// PostBlocks send blocks to a channel
func (r *poster) PostBlocks(channel string, blocks []slack.Block, options ...PostOption) {
	r.post(channel, "", blocks, options...)
}

func (r *poster) post(channel string, message string, blocks []slack.Block, options ...PostOption) {
	postOptions := newPostOptions(options...)

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(postOptions.Attachments...),
		slack.MsgOptionBlocks(blocks...),
	}

	if len(postOptions.ThreadTS) > 0 {
		opts = append(opts, slack.MsgOptionTS(postOptions.ThreadTS))
	}

	_, _, err := r.apiClient.PostMessage(
		channel,
		opts...,
	)
	if err != nil {
		infof("failed to post message: %v\n", err)
	}
}
