package slacker

import "github.com/nlopes/slack"

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

// ReplyDefaults configuration
type ReplyDefaults struct {
	Attachments []slack.Attachment
}

func newReplyDefaults(options ...ReplyOption) *ReplyDefaults {
	config := &ReplyDefaults{
		Attachments: []slack.Attachment{},
	}

	for _, option := range options {
		option(config)
	}
	return config
}

// HelpOption an option for help values
type HelpOption func(*HelpDefaults)

// WithHandler sets help handler
func WithHandler(handler func(request Request, response ResponseWriter)) HelpOption {
	return func(defaults *HelpDefaults) {
		defaults.Handler = handler
	}
}

// HelpDefaults configuration
type HelpDefaults struct {
	Handler func(request Request, response ResponseWriter)
}

func newHelpDefaults(options ...HelpOption) *HelpDefaults {
	config := &HelpDefaults{}

	for _, option := range options {
		option(config)
	}
	return config
}
