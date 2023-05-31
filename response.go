package slacker

import (
	"github.com/slack-go/slack"
)

// PosterReplierResponse interface is used to respond to an event
type PosterReplierResponse interface {
	Reply(message string, options ...ReplyOption)
	ReplyError(err error, options ...ReplyOption)
	ReplyBlocks(blocks []slack.Block, options ...ReplyOption)
	Post(channel string, message string, options ...PostOption)
	PostError(channel string, err error, options ...PostOption)
	PostBlocks(channel string, blocks []slack.Block, options ...PostOption)
}

// newPosterReplierResponse creates a new response structure
func newPosterReplierResponse(poster Poster, replier Replier) PosterReplierResponse {
	return &posterReplierResponse{poster: poster, replier: replier}
}

type posterReplierResponse struct {
	poster  Poster
	replier Replier
}

// Reply send a message to the current channel
func (r *posterReplierResponse) Reply(message string, options ...ReplyOption) {
	r.replier.Reply(message, options...)
}

// ReplyError send an error to the current channel
func (r *posterReplierResponse) ReplyError(err error, options ...ReplyOption) {
	r.replier.ReplyError(err, options...)
}

// ReplyBlocks send blocks to the current channel
func (r *posterReplierResponse) ReplyBlocks(blocks []slack.Block, options ...ReplyOption) {
	r.replier.ReplyBlocks(blocks, options...)
}

// Post send a message to a channel
func (r *posterReplierResponse) Post(channel string, message string, options ...PostOption) {
	r.poster.Post(channel, message, options...)
}

// PostError send an error to a channel
func (r *posterReplierResponse) PostError(channel string, err error, options ...PostOption) {
	r.poster.PostError(channel, err, options...)
}

// PostBlocks send blocks to a channel
func (r *posterReplierResponse) PostBlocks(channel string, blocks []slack.Block, options ...PostOption) {
	r.poster.PostBlocks(channel, blocks, options...)
}

// PosterResponse interface is used to respond to an event
type PosterResponse interface {
	Post(channel string, message string, options ...PostOption)
	PostError(channel string, err error, options ...PostOption)
	PostBlocks(channel string, blocks []slack.Block, options ...PostOption)
}

// newPosterResponse creates a new response structure
func newPosterResponse(poster Poster) PosterResponse {
	return &posterReplierResponse{poster: poster}
}

type posterResponse struct {
	poster Poster
}

// Post send a message to a channel
func (r *posterResponse) Post(channel string, message string, options ...PostOption) {
	r.poster.Post(channel, message, options...)
}

// PostError send an error to a channel
func (r *posterResponse) PostError(channel string, err error, options ...PostOption) {
	r.poster.PostError(channel, err, options...)
}

// PostBlocks send blocks to a channel
func (r *posterResponse) PostBlocks(channel string, blocks []slack.Block, options ...PostOption) {
	r.poster.PostBlocks(channel, blocks, options...)
}
