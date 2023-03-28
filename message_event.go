package slacker

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

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

// NewMessageEvent creates a new message event structure 
func NewMessageEvent(slacker *Slacker, evt interface{}, req *socketmode.Request) *MessageEvent {
	var me *MessageEvent

	switch ev := evt.(type) {
	case *slackevents.MessageEvent:
		me = &MessageEvent{
			Channel:         ev.Channel,
			ChannelName:     getChannelName(slacker, ev.Channel),
			User:            ev.User,
			UserName:        getUserName(slacker, ev.User),
			Text:            ev.Text,
			Data:            evt,
			Type:            ev.Type,
			TimeStamp:       ev.TimeStamp,
			ThreadTimeStamp: ev.ThreadTimeStamp,
			BotID:           ev.BotID,
		}
	case *slackevents.AppMentionEvent:
		me = &MessageEvent{
			Channel:         ev.Channel,
			ChannelName:     getChannelName(slacker, ev.Channel),
			User:            ev.User,
			UserName:        getUserName(slacker, ev.User),
			Text:            ev.Text,
			Data:            evt,
			Type:            ev.Type,
			TimeStamp:       ev.TimeStamp,
			ThreadTimeStamp: ev.ThreadTimeStamp,
			BotID:           ev.BotID,
		}
	case *slack.SlashCommand:
		me = &MessageEvent{
			Channel:     ev.ChannelID,
			ChannelName: ev.ChannelName,
			User:        ev.UserID,
			UserName:    ev.UserName,
			Text:        fmt.Sprintf("%s %s", ev.Command[1:], ev.Text),
			Data:        req,
			Type:        req.Type,
		}
	}

	// Filter out other bots. At the very least this is needed for MessageEvent
	// to prevent the bot from self-triggering and causing loops. However better
	// logic should be in place to prevent repeated self-triggering / bot-storms
	// if we want to enable this later.
	if me.IsBot() {
		return nil
	}
	return me
}

func getChannelName(slacker *Slacker, channelID string) string {
	channel, err := slacker.client.GetConversationInfo(channelID, true)
	if err != nil {
		fmt.Printf("unable to get channel info for %s: %v\n", channelID, err)
		return channelID
	}
	return channel.Name
}

func getUserName(slacker *Slacker, userID string) string {
	user, err := slacker.client.GetUserInfo(userID)
	if err != nil {
		fmt.Printf("unable to get user info for %s: %v\n", userID, err)
		return userID
	}
	return user.Name
}
