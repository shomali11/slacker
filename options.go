package slacker

import "github.com/slack-go/slack"

const (
	errorFormat = "```%s```"
)

// ClientOption an option for client values
type ClientOption func(*ClientOptions)

// WithDebug sets debug toggle
func WithDebug(debug bool) ClientOption {
	return func(defaults *ClientOptions) {
		defaults.Debug = debug
	}
}

// WithBotInteractionMode instructs Slacker on how to handle message events coming from a bot.
func WithBotInteractionMode(mode BotInteractionMode) ClientOption {
	return func(defaults *ClientOptions) {
		defaults.BotMode = mode
	}
}

// ClientOptions configuration
type ClientOptions struct {
	Debug   bool
	BotMode BotInteractionMode
}

func newClientOptions(options ...ClientOption) *ClientOptions {
	config := &ClientOptions{
		Debug:   false,
		BotMode: BotInteractionModeIgnoreAll,
	}

	for _, option := range options {
		option(config)
	}
	return config
}

// ReplyOption an option for reply values
type ReplyOption func(*ReplyOptions)

// WithAttachments sets message attachments
func WithAttachments(attachments []slack.Attachment) ReplyOption {
	return func(defaults *ReplyOptions) {
		defaults.Attachments = attachments
	}
}

// WithBlocks sets message blocks
func WithBlocks(blocks []slack.Block) ReplyOption {
	return func(defaults *ReplyOptions) {
		defaults.Blocks = blocks
	}
}

// WithThreadReply specifies whether to reply inside a thread of the original message
func WithThreadReply(useThread bool) ReplyOption {
	return func(defaults *ReplyOptions) {
		defaults.ThreadResponse = useThread
	}
}

// ReplyOptions configuration
type ReplyOptions struct {
	Attachments    []slack.Attachment
	Blocks         []slack.Block
	ThreadResponse bool
}

// newReplyOptions builds our ReplyOptionss from zero or more ReplyOption.
func newReplyOptions(options ...ReplyOption) *ReplyOptions {
	config := &ReplyOptions{
		Attachments:    []slack.Attachment{},
		Blocks:         []slack.Block{},
		ThreadResponse: false,
	}

	for _, option := range options {
		option(config)
	}
	return config
}

// ErrorOption an option for error values
type ErrorOption func(*ErrorOptions)

// ErrorOptions configuration
type ErrorOptions struct {
	ThreadResponse bool
	Format         string
}

// WithThreadError specifies whether to error inside a thread of the original message
func WithThreadError(useThread bool) ErrorOption {
	return func(defaults *ErrorOptions) {
		defaults.ThreadResponse = useThread
	}
}

// WithFormat specifies the format of the error
func WithFormat(format string) ErrorOption {
	return func(defaults *ErrorOptions) {
		defaults.Format = format
	}
}

// newErrorOptions builds our ReportErrorOptions from zero or more ReportErrorOption.
func newErrorOptions(options ...ErrorOption) *ErrorOptions {
	config := &ErrorOptions{
		ThreadResponse: false,
		Format:         errorFormat,
	}

	for _, option := range options {
		option(config)
	}
	return config
}
