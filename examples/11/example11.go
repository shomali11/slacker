package main

import (
	"log"

	"context"
	"errors"
	"fmt"

	"github.com/shomali11/slacker"
)

const (
	errorFormat = "> Custom Error: _%s_"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	bot.CustomResponse(NewCustomResponseWriter)

	definition := &slacker.CommandDefinition{
		Description: "Custom!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
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
func NewCustomResponseWriter(botCtx slacker.BotContext) slacker.ResponseWriter {
	return &MyCustomResponseWriter{botCtx: botCtx}
}

// MyCustomResponseWriter a custom response writer
type MyCustomResponseWriter struct {
	botCtx slacker.BotContext
}

// ReportError sends back a formatted error message to the channel where we received the event from
func (r *MyCustomResponseWriter) ReportError(err error, options ...slacker.ReportErrorOption) {
	rtm := r.botCtx.RTM()
	event := r.botCtx.Event()
	rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf(errorFormat, err.Error()), event.Channel))
}

// Typing send a typing indicator
func (r *MyCustomResponseWriter) Typing() {
	rtm := r.botCtx.RTM()
	event := r.botCtx.Event()
	rtm.SendMessage(rtm.NewTypingMessage(event.Channel))
}

// Reply send a attachments to the current channel with a message
func (r *MyCustomResponseWriter) Reply(message string, options ...slacker.ReplyOption) error {
	rtm := r.botCtx.RTM()
	event := r.botCtx.Event()
	rtm.SendMessage(rtm.NewOutgoingMessage(message, event.Channel))
	return nil
}
