package slacker

import "github.com/slack-go/slack"

// ClientOption an option for client values
type ClientOption func(*ClientDefaults)

// WithAPIURL sets the API URL (for testing)
func WithAPIURL(url string) ClientOption {
	return func(defaults *ClientDefaults) {
		defaults.APIURL = url
	}
}

// WithDebug sets debug toggle
func WithDebug(debug bool) ClientOption {
	return func(defaults *ClientDefaults) {
		defaults.Debug = debug
	}
}

// WithBotInteractionMode instructs Slacker on how to handle message events coming from a bot.
func WithBotInteractionMode(mode BotInteractionMode) ClientOption {
	return func(defaults *ClientDefaults) {
		defaults.BotMode = mode
	}
}

// WithoutFontFormatting disables font formatting in messages
func WithoutFontFormatting() ClientOption {
	return func(defaults *ClientDefaults) {
		defaults.MessageWithoutFontFormatting = true
	}
}

// WithoutAllFormatting disables font formatting in messages
func WithoutAllFormatting() ClientOption {
	return func(defaults *ClientDefaults) {
		defaults.MessageWithoutAllFormatting = true
	}
}

// ClientDefaults configuration
type ClientDefaults struct {
	APIURL                       string
	Debug                        bool
	BotMode                      BotInteractionMode
	MessageWithoutFontFormatting bool
	MessageWithoutAllFormatting  bool
}

func newClientDefaults(options ...ClientOption) *ClientDefaults {
	config := &ClientDefaults{
		APIURL:  "", // Empty string will not override default from slack package
		Debug:   false,
		BotMode: BotInteractionModeIgnoreAll,
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

// WithThreadReply specifies the reply to be inside a thread of the original message
func WithThreadReply(useThread bool) ReplyOption {
	return func(defaults *ReplyDefaults) {
		defaults.ThreadResponse = useThread
	}
}

// ReplyDefaults configuration
type ReplyDefaults struct {
	Attachments    []slack.Attachment
	Blocks         []slack.Block
	ThreadResponse bool
}

// NewReplyDefaults builds our ReplyDefaults from zero or more ReplyOption.
func NewReplyDefaults(options ...ReplyOption) *ReplyDefaults {
	config := &ReplyDefaults{
		Attachments:    []slack.Attachment{},
		Blocks:         []slack.Block{},
		ThreadResponse: false,
	}

	for _, option := range options {
		option(config)
	}
	return config
}

// ReportErrorOption an option for report error values
type ReportErrorOption func(*ReportErrorDefaults)

// ReportErrorDefaults configuration
type ReportErrorDefaults struct {
	ThreadResponse bool
}

// WithThreadReplyError specifies the reply to be inside a thread of the original message
func WithThreadReplyError(useThread bool) ReportErrorOption {
	return func(defaults *ReportErrorDefaults) {
		defaults.ThreadResponse = useThread
	}
}

// NewReportErrorDefaults builds our ReportErrorDefaults from zero or more ReportErrorOption.
func NewReportErrorDefaults(options ...ReportErrorOption) *ReportErrorDefaults {
	config := &ReportErrorDefaults{
		ThreadResponse: false,
	}

	for _, option := range options {
		option(config)
	}
	return config
}
