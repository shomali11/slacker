package slacker

import "github.com/slack-go/slack"

// ClientOption an option for client values
type ClientOption func(*ClientDefaults)

// WithDebug sets debug toggle
func WithDebug(debug bool) ClientOption {
	return func(defaults *ClientDefaults) {
		defaults.Debug = debug
	}
}

// ClientDefaults configuration
type ClientDefaults struct {
	Debug bool
}

func newClientDefaults(options ...ClientOption) *ClientDefaults {
	config := &ClientDefaults{
		Debug: false,
	}

	for _, option := range options {
		option(config)
	}
	return config
}

// ReplyOption an option for reply values
type ReplyOption func(*ReplyDefaults)

// WithAttachments sets message attachments
func WithAttachments(attachments []slack.Attachment) ReplyOption {
	return func(defaults *ReplyDefaults) {
		defaults.Attachments = attachments
	}
}

// WithBlocks sets message blocks
func WithBlocks(blocks []slack.Block) ReplyOption {
	return func(defaults *ReplyDefaults) {
		defaults.Blocks = blocks
	}
}

// ReplyDefaults configuration
type ReplyDefaults struct {
	Attachments []slack.Attachment
	Blocks      []slack.Block
}

func newReplyDefaults(options ...ReplyOption) *ReplyDefaults {
	config := &ReplyDefaults{
		Attachments: []slack.Attachment{},
		Blocks:      []slack.Block{},
	}

	for _, option := range options {
		option(config)
	}
	return config
}
