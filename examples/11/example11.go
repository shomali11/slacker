package main

import (
	"log"

	"context"
	"errors"
	"fmt"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

const (
	errorFormat = "> Custom Error: _%s_"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.CustomResponse(NewCustomResponseWriter)

	definition := &slacker.CommandDefinition{
		Description: "Custom!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("custom")
			response.ReportError(errors.New("oops"))
		},
	}

	bot.Command("custom", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// NewCustomResponseWriter creates a new ResponseWriter structure
func NewCustomResponseWriter(channel string, client *slack.Client, rtm *slack.RTM) slacker.ResponseWriter {
	return &MyCustomResponseWriter{channel: channel, client: client, rtm: rtm}
}

// MyCustomResponseWriter a custom response writer
type MyCustomResponseWriter struct {
	channel string
	client  *slack.Client
	rtm     *slack.RTM
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *MyCustomResponseWriter) ReportError(err error) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), r.channel))
}

// Typing send a typing indicator
func (r *MyCustomResponseWriter) Typing() {
	r.rtm.SendMessage(r.rtm.NewTypingMessage(r.channel))
}

// Reply send a attachments to the current channel with a message
func (r *MyCustomResponseWriter) Reply(message string, options ...slacker.ReplyOption) {
	r.rtm.SendMessage(r.rtm.NewOutgoingMessage(message, r.channel))
}

// RTM returns the RTM client
func (r *MyCustomResponseWriter) RTM() *slack.RTM {
	return r.rtm
}

// Client returns the slack client
func (r *MyCustomResponseWriter) Client() *slack.Client {
	return r.client
}
