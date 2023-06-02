package slacker

import (
	"github.com/slack-go/slack"
)

// WriterReplierResponse interface is used to respond to an event
type WriterReplierResponse interface {
	Reply(message string, options ...ReplyOption) (string, error)
	ReplyError(err error, options ...ReplyOption) (string, error)
	ReplyBlocks(blocks []slack.Block, options ...ReplyOption) (string, error)
	Post(channel string, message string, options ...PostOption) (string, error)
	PostError(channel string, err error, options ...PostOption) (string, error)
	PostBlocks(channel string, blocks []slack.Block, options ...PostOption) (string, error)
	Delete(channel string, messageTimestamp string) (string, error)
}

// newWriterReplierResponse creates a new response structure
func newWriterReplierResponse(poster Writer, replier Replier) WriterReplierResponse {
	return &writerReplierResponse{poster: poster, replier: replier}
}

type writerReplierResponse struct {
	poster  Writer
	replier Replier
}

// Reply send a message to the current channel
func (r *writerReplierResponse) Reply(message string, options ...ReplyOption) (string, error) {
	return r.replier.Reply(message, options...)
}

// ReplyError send an error to the current channel
func (r *writerReplierResponse) ReplyError(err error, options ...ReplyOption) (string, error) {
	return r.replier.ReplyError(err, options...)
}

// ReplyBlocks send blocks to the current channel
func (r *writerReplierResponse) ReplyBlocks(blocks []slack.Block, options ...ReplyOption) (string, error) {
	return r.replier.ReplyBlocks(blocks, options...)
}

// Post send a message to a channel
func (r *writerReplierResponse) Post(channel string, message string, options ...PostOption) (string, error) {
	return r.poster.Post(channel, message, options...)
}

// PostError send an error to a channel
func (r *writerReplierResponse) PostError(channel string, err error, options ...PostOption) (string, error) {
	return r.poster.PostError(channel, err, options...)
}

// PostBlocks send blocks to a channel
func (r *writerReplierResponse) PostBlocks(channel string, blocks []slack.Block, options ...PostOption) (string, error) {
	return r.poster.PostBlocks(channel, blocks, options...)
}

// Delete deletes a message in a channel
func (r *writerReplierResponse) Delete(channel string, messageTimestamp string) (string, error) {
	return r.poster.Delete(channel, messageTimestamp)
}

// WriterResponse interface is used to respond to an event
type WriterResponse interface {
	Post(channel string, message string, options ...PostOption) (string, error)
	PostError(channel string, err error, options ...PostOption) (string, error)
	PostBlocks(channel string, blocks []slack.Block, options ...PostOption) (string, error)
	Delete(channel string, messageTimestamp string) (string, error)
}

// newWriterResponse creates a new response structure
func newWriterResponse(poster Writer) WriterResponse {
	return &writerResponse{poster: poster}
}

type writerResponse struct {
	poster Writer
}

// Post send a message to a channel
func (r *writerResponse) Post(channel string, message string, options ...PostOption) (string, error) {
	return r.poster.Post(channel, message, options...)
}

// PostError send an error to a channel
func (r *writerResponse) PostError(channel string, err error, options ...PostOption) (string, error) {
	return r.poster.PostError(channel, err, options...)
}

// PostBlocks send blocks to a channel
func (r *writerResponse) PostBlocks(channel string, blocks []slack.Block, options ...PostOption) (string, error) {
	return r.poster.PostBlocks(channel, blocks, options...)
}

// Delete deletes a message in a channel
func (r *writerResponse) Delete(channel string, messageTimestamp string) (string, error) {
	return r.poster.Delete(channel, messageTimestamp)
}
