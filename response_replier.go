package slacker

import (
	"github.com/slack-go/slack"
)

// A Replier interface is used to respond to an event
type Replier interface {
	Reply(message string, options ...ReplyOption) (string, error)
	ReplyError(err error, options ...ReplyOption) (string, error)
	ReplyBlocks(blocks []slack.Block, options ...ReplyOption) (string, error)
}

// newReplier creates a new replier structure
func newReplier(channelID string, userID string, eventTS string, poster Writer) Replier {
	return &replier{channelID: channelID, userID: userID, eventTS: eventTS, poster: poster}
}

type replier struct {
	channelID string
	userID    string
	eventTS   string
	poster    Writer
}

// Reply send a message to the current channel
func (r *replier) Reply(message string, options ...ReplyOption) (string, error) {
	responseOptions := r.convertOptions(options...)
	return r.poster.Post(r.channelID, message, responseOptions...)
}

// ReplyError send an error to the current channel
func (r *replier) ReplyError(err error, options ...ReplyOption) (string, error) {
	responseOptions := r.convertOptions(options...)
	return r.poster.PostError(r.channelID, err, responseOptions...)
}

// ReplyBlocks send blocks to the current channel
func (r *replier) ReplyBlocks(blocks []slack.Block, options ...ReplyOption) (string, error) {
	responseOptions := r.convertOptions(options...)
	return r.poster.PostBlocks(r.channelID, blocks, responseOptions...)
}

func (r *replier) convertOptions(options ...ReplyOption) []PostOption {
	replyOptions := newReplyOptions(options...)
	responseOptions := []PostOption{
		SetAttachments(replyOptions.Attachments),
	}

	if replyOptions.InThread {
		responseOptions = append(responseOptions, SetThreadTS(r.eventTS))
	}

	if len(replyOptions.ReplaceMessageTS) > 0 {
		responseOptions = append(responseOptions, SetReplace(replyOptions.ReplaceMessageTS))
	}

	if replyOptions.IsEphemeral {
		responseOptions = append(responseOptions, SetEphemeral(r.userID))
	}
	return responseOptions
}
