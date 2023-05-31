package slacker

import (
	"github.com/slack-go/slack"
)

// A Replier interface is used to respond to an event
type Replier interface {
	Reply(message string, options ...ReplyOption)
	ReplyError(err error, options ...ReplyOption)
	ReplyBlocks(blocks []slack.Block, options ...ReplyOption)
}

// newReplier creates a new replier structure
func newReplier(channelID string, eventTS string, poster Poster) Replier {
	return &replier{channelID: channelID, eventTS: eventTS, poster: poster}
}

type replier struct {
	channelID string
	eventTS   string
	poster    Poster
}

// Reply send a message to the current channel
func (r *replier) Reply(message string, options ...ReplyOption) {
	responseOptions := r.convertOptions(options...)
	r.poster.Post(r.channelID, message, responseOptions...)
}

// ReplyError send an error to the current channel
func (r *replier) ReplyError(err error, options ...ReplyOption) {
	responseOptions := r.convertOptions(options...)
	r.poster.PostError(r.channelID, err, responseOptions...)
}

// ReplyBlocks send blocks to the current channel
func (r *replier) ReplyBlocks(blocks []slack.Block, options ...ReplyOption) {
	responseOptions := r.convertOptions(options...)
	r.poster.PostBlocks(r.channelID, blocks, responseOptions...)
}

func (r *replier) convertOptions(options ...ReplyOption) []PostOption {
	replyOptions := newReplyOptions(options...)
	responseOptions := []PostOption{
		SetAttachments(replyOptions.Attachments),
	}

	if replyOptions.InThread {
		responseOptions = append(responseOptions, SetThreadTS(r.eventTS))
	}
	return responseOptions
}
