package slacker

import (
	"time"

	"github.com/slack-go/slack"
)

// ClientOption an option for client values
type ClientOption func(*clientOptions)

// WithAPIURL sets the API URL (for testing)
func WithAPIURL(url string) ClientOption {
	return func(defaults *clientOptions) {
		defaults.APIURL = url
	}
}

// WithDebug sets debug toggle
func WithDebug(debug bool) ClientOption {
	return func(defaults *clientOptions) {
		defaults.Debug = debug
	}
}

// WithBotMode instructs Slacker on how to handle message events coming from a bot.
func WithBotMode(mode BotMode) ClientOption {
	return func(defaults *clientOptions) {
		defaults.BotMode = mode
	}
}

// WithLogger sets slacker logger
func WithLogger(logger Logger) ClientOption {
	return func(defaults *clientOptions) {
		defaults.Logger = logger
	}
}

// WithCronLocation overrides the timezone of the cron instance.
func WithCronLocation(location *time.Location) ClientOption {
	return func(defaults *clientOptions) {
		defaults.CronLocation = location
	}
}

type clientOptions struct {
	APIURL       string
	Debug        bool
	BotMode      BotMode
	Logger       Logger
	CronLocation *time.Location
}

func newClientOptions(options ...ClientOption) *clientOptions {
	config := &clientOptions{
		APIURL:       slack.APIURL,
		Debug:        false,
		BotMode:      BotModeIgnoreAll,
		CronLocation: time.Local,
	}

	for _, option := range options {
		option(config)
	}

	if config.Logger == nil {
		config.Logger = newBuiltinLogger(config.Debug)
	}
	return config
}

// ReplyOption an option for reply values
type ReplyOption func(*replyOptions)

// WithAttachments sets message attachments
func WithAttachments(attachments []slack.Attachment) ReplyOption {
	return func(defaults *replyOptions) {
		defaults.Attachments = attachments
	}
}

// WithInThread specifies whether to reply inside a thread of the original message
func WithInThread() ReplyOption {
	return func(defaults *replyOptions) {
		defaults.InThread = true
	}
}

// WithReplace replaces the original message
func WithReplace(originalMessageTS string) ReplyOption {
	return func(defaults *replyOptions) {
		defaults.ReplaceMessageTS = originalMessageTS
	}
}

// WithEphemeral sets the message as ephemeral
func WithEphemeral() ReplyOption {
	return func(defaults *replyOptions) {
		defaults.IsEphemeral = true
	}
}

// WithSchedule sets message's schedule
func WithSchedule(timestamp time.Time) ReplyOption {
	return func(defaults *replyOptions) {
		defaults.ScheduleTime = &timestamp
	}
}

type replyOptions struct {
	Attachments      []slack.Attachment
	InThread         bool
	ReplaceMessageTS string
	IsEphemeral      bool
	ScheduleTime     *time.Time
}

// newReplyOptions builds our ReplyOptions from zero or more ReplyOption.
func newReplyOptions(options ...ReplyOption) *replyOptions {
	config := &replyOptions{
		Attachments: []slack.Attachment{},
		InThread:    false,
	}

	for _, option := range options {
		option(config)
	}
	return config
}

// PostOption an option for post values
type PostOption func(*postOptions)

// SetAttachments sets message attachments
func SetAttachments(attachments []slack.Attachment) PostOption {
	return func(defaults *postOptions) {
		defaults.Attachments = attachments
	}
}

// SetThreadTS specifies whether to reply inside a thread
func SetThreadTS(threadTS string) PostOption {
	return func(defaults *postOptions) {
		defaults.ThreadTS = threadTS
	}
}

// SetReplace sets message url to be replaced
func SetReplace(originalMessageTS string) PostOption {
	return func(defaults *postOptions) {
		defaults.ReplaceMessageTS = originalMessageTS
	}
}

// SetEphemeral sets the user who receives the ephemeral message
func SetEphemeral(userID string) PostOption {
	return func(defaults *postOptions) {
		defaults.EphemeralUserID = userID
	}
}

// SetSchedule sets message's schedule
func SetSchedule(timestamp time.Time) PostOption {
	return func(defaults *postOptions) {
		defaults.ScheduleTime = &timestamp
	}
}

type postOptions struct {
	Attachments      []slack.Attachment
	ThreadTS         string
	ReplaceMessageTS string
	EphemeralUserID  string
	ScheduleTime     *time.Time
}

// newPostOptions builds our PostOptions from zero or more PostOption.
func newPostOptions(options ...PostOption) *postOptions {
	config := &postOptions{
		Attachments: []slack.Attachment{},
	}

	for _, option := range options {
		option(config)
	}
	return config
}
