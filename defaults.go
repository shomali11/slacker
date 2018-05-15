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

// Defaults configuration
type Defaults struct {
	Attachments []slack.Attachment
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
