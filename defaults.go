package slacker

import "github.com/nlopes/slack"

// DefaultsOption an option for default values
type DefaultsOption func(*Defaults)

// WithAttachments sets message attachments
func WithAttachments(attachments []slack.Attachment) DefaultsOption {
	return func(defaults *Defaults) {
		defaults.Attachments = attachments
	}
}

// WithDebug sets message attachments
func WithDebug(debug bool) DefaultsOption {
	return func(defaults *Defaults) {
		defaults.Debug = debug
	}
}

// Defaults configuration
type Defaults struct {
	Attachments []slack.Attachment
	Debug       bool
}

func newDefaults(options ...DefaultsOption) *Defaults {
	config := &Defaults{
		Attachments: []slack.Attachment{},
	}

	for _, option := range options {
		option(config)
	}
	return config
}
