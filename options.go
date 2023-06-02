package slacker

import "github.com/slack-go/slack"

// ClientOption an option for client values
type ClientOption func(*ClientOptions)

// WithAPIURL sets the API URL (for testing)
func WithAPIURL(url string) ClientOption {
	return func(defaults *ClientOptions) {
		defaults.APIURL = url
	}
}

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
	APIURL  string
	Debug   bool
	BotMode BotInteractionMode
}

func newClientOptions(options ...ClientOption) *ClientOptions {
	config := &ClientOptions{
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
type ReplyOption func(*ReplyOptions)

// WithAttachments sets message attachments
func WithAttachments(attachments []slack.Attachment) ReplyOption {
	return func(defaults *ReplyOptions) {
		defaults.Attachments = attachments
	}
}

// WithInThread specifies whether to reply inside a thread of the original message
func WithInThread() ReplyOption {
	return func(defaults *ReplyOptions) {
		defaults.InThread = true
	}
}

// WithReplace replaces the original message
func WithReplace(originalMessageTS string) ReplyOption {
	return func(defaults *ReplyOptions) {
		defaults.ReplaceMessageTS = originalMessageTS
	}
}

// WithEphemeral sets the message as ephemeral
func WithEphemeral() ReplyOption {
	return func(defaults *ReplyOptions) {
		defaults.IsEphemeral = true
	}
}

// ReplyOptions configuration
type ReplyOptions struct {
	Attachments      []slack.Attachment
	InThread         bool
	ReplaceMessageTS string
	IsEphemeral      bool
}

// newReplyOptions builds our ReplyOptions from zero or more ReplyOption.
func newReplyOptions(options ...ReplyOption) *ReplyOptions {
	config := &ReplyOptions{
		Attachments: []slack.Attachment{},
		InThread:    false,
	}

	for _, option := range options {
		option(config)
	}
	return config
}

// PostOption an option for post values
type PostOption func(*PostOptions)

// SetAttachments sets message attachments
func SetAttachments(attachments []slack.Attachment) PostOption {
	return func(defaults *PostOptions) {
		defaults.Attachments = attachments
	}
}

// SetThreadTS specifies whether to reply inside a thread
func SetThreadTS(threadTS string) PostOption {
	return func(defaults *PostOptions) {
		defaults.ThreadTS = threadTS
	}
}

// SetReplace sets message url to be replaced
func SetReplace(originalMessageTS string) PostOption {
	return func(defaults *PostOptions) {
		defaults.ReplaceMessageTS = originalMessageTS
	}
}

// SetEphemeral sets the user who receives the ephemeral message
func SetEphemeral(userID string) PostOption {
	return func(defaults *PostOptions) {
		defaults.EphemeralUserID = userID
	}
}

// PostOptions configuration
type PostOptions struct {
	Attachments      []slack.Attachment
	ThreadTS         string
	ReplaceMessageTS string
	EphemeralUserID  string
}

// newPostOptions builds our PostOptions from zero or more PostOption.
func newPostOptions(options ...PostOption) *PostOptions {
	config := &PostOptions{
		Attachments: []slack.Attachment{},
		ThreadTS:    "",
	}

	for _, option := range options {
		option(config)
	}
	return config
}
