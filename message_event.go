package slacker

// MessageEvent contains details common to message based events, including the
// raw event as returned from Slack along with the corresponding event type.
// The struct should be kept minimal and only include data that is commonly
// used to prevent frequent type assertions when evaluating the event.
type MessageEvent struct {
	// Channel ID where the message was sent
	Channel string

	// ChannelName where the message was sent
	ChannelName string

	// User ID of the sender
	User string

	// UserName of the the sender
	UserName string

	// Text is the unalterted text of the message, as returned by Slack
	Text string

	// TimeStamp is the message timestamp. For events that do not support
	// threading (eg. slash commands) this will be unset.
	// will be left unset.
	TimeStamp string

	// ThreadTimeStamp is the message thread timestamp. For events that do not
	// support threading (eg. slash commands) this will be unset.
	ThreadTimeStamp string

	// Data is the raw event data returned from slack. Using Type, you can assert
	// this into a slackevents *Event struct.
	Data interface{}

	// Type is the type of the event, as returned by Slack. For instance,
	// `app_mention` or `message`
	Type string

	// BotID of the bot that sent this message. If a bot did not send this
	// message, this will be an empty string.
	BotID string
}

// IsThread indicates if a message event took place in a thread.
func (e *MessageEvent) IsThread() bool {
	if e.ThreadTimeStamp == "" || e.ThreadTimeStamp == e.TimeStamp {
		return false
	}
	return true
}

// IsBot indicates if the message was sent by a bot
func (e *MessageEvent) IsBot() bool {
	return e.BotID != ""
}
