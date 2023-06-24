package slacker

import (
	"github.com/slack-go/slack"
)

// newResponseReplier creates a new response structure
func newResponseReplier(writer *Writer, replier *Replier) *ResponseReplier {
	return &ResponseReplier{writer: writer, replier: replier}
}

// ResponseReplier sends messages to Slack
type ResponseReplier struct {
	writer  *Writer
	replier *Replier
}

// Reply send a message to the current channel
func (r *ResponseReplier) Reply(message string, options ...ReplyOption) (string, error) {
	return r.replier.Reply(message, options...)
}

// ReplyError send an error to the current channel
func (r *ResponseReplier) ReplyError(err error, options ...ReplyOption) (string, error) {
	return r.replier.ReplyError(err, options...)
}

// ReplyBlocks send blocks to the current channel
func (r *ResponseReplier) ReplyBlocks(blocks []slack.Block, options ...ReplyOption) (string, error) {
	return r.replier.ReplyBlocks(blocks, options...)
}

// Post send a message to a channel
func (r *ResponseReplier) Post(channel string, message string, options ...PostOption) (string, error) {
	return r.writer.Post(channel, message, options...)
}

// PostError send an error to a channel
func (r *ResponseReplier) PostError(channel string, err error, options ...PostOption) (string, error) {
	return r.writer.PostError(channel, err, options...)
}

// PostBlocks send blocks to a channel
func (r *ResponseReplier) PostBlocks(channel string, blocks []slack.Block, options ...PostOption) (string, error) {
	return r.writer.PostBlocks(channel, blocks, options...)
}

// Delete deletes a message in a channel
func (r *ResponseReplier) Delete(channel string, messageTimestamp string) (string, error) {
	return r.writer.Delete(channel, messageTimestamp)
}

// newWriterResponse creates a new response structure
func newWriterResponse(writer *Writer) *ResponseWriter {
	return &ResponseWriter{writer: writer}
}

// ResponseWriter sends messages to slack
type ResponseWriter struct {
	writer *Writer
}

// Post send a message to a channel
func (r *ResponseWriter) Post(channel string, message string, options ...PostOption) (string, error) {
	return r.writer.Post(channel, message, options...)
}

// PostError send an error to a channel
func (r *ResponseWriter) PostError(channel string, err error, options ...PostOption) (string, error) {
	return r.writer.PostError(channel, err, options...)
}

// PostBlocks send blocks to a channel
func (r *ResponseWriter) PostBlocks(channel string, blocks []slack.Block, options ...PostOption) (string, error) {
	return r.writer.PostBlocks(channel, blocks, options...)
}

// Delete deletes a message in a channel
func (r *ResponseWriter) Delete(channel string, messageTimestamp string) (string, error) {
	return r.writer.Delete(channel, messageTimestamp)
}
