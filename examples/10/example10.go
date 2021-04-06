package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

const (
	errorFormat = "> Custom Error: _%s_"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

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
	defaults := slacker.NewReportErrorDefaults(options...)

	client := r.botCtx.Client()
	event := r.botCtx.Event()

	opts := []slack.MsgOption{
		slack.MsgOptionText(fmt.Sprintf(errorFormat, err.Error()), false),
	}
	if defaults.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(event.TimeStamp))
	}

	_, _, err = client.PostMessage(event.Channel, opts...)
	if err != nil {
		fmt.Printf("failed to report error: %v\n", err)
	}
}

// Reply send a attachments to the current channel with a message
func (r *MyCustomResponseWriter) Reply(message string, options ...slacker.ReplyOption) error {
	defaults := slacker.NewReplyDefaults(options...)

	client := r.botCtx.Client()
	event := r.botCtx.Event()
	if event == nil {
		return fmt.Errorf("Unable to get message event details")
	}

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(defaults.Attachments...),
		slack.MsgOptionBlocks(defaults.Blocks...),
	}
	if defaults.ThreadResponse {
		opts = append(opts, slack.MsgOptionTS(event.TimeStamp))
	}

	_, _, err := client.PostMessage(
		event.Channel,
		opts...,
	)
	return err
}
