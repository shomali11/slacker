package slacker

import (
	"github.com/slack-go/slack"
)

// newReplier creates a new replier structure
func newReplier(channelID string, userID string, eventTS string, writer *Writer) *Replier {
	return &Replier{channelID: channelID, userID: userID, eventTS: eventTS, writer: writer}
}

// Replier sends messages to the same channel the event came from
type Replier struct {
	channelID string
	userID    string
	eventTS   string
	writer    *Writer
}

// Reply send a message to the current channel
func (r *Replier) Reply(message string, options ...ReplyOption) (string, error) {
	responseOptions := r.convertOptions(options...)
	return r.writer.Post(r.channelID, message, responseOptions...)
}

// ReplyError send an error to the current channel
func (r *Replier) ReplyError(err error, options ...ReplyOption) (string, error) {
	responseOptions := r.convertOptions(options...)
	return r.writer.PostError(r.channelID, err, responseOptions...)
}

// ReplyBlocks send blocks to the current channel
func (r *Replier) ReplyBlocks(blocks []slack.Block, options ...ReplyOption) (string, error) {
	responseOptions := r.convertOptions(options...)
	return r.writer.PostBlocks(r.channelID, blocks, responseOptions...)
}

func (r *Replier) convertOptions(options ...ReplyOption) []PostOption {
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

	if replyOptions.ScheduleTime != nil {
		responseOptions = append(responseOptions, SetSchedule(*replyOptions.ScheduleTime))
	}
	return responseOptions
}
