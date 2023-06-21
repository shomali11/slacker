package main

import (
	"log"
	"os"

	"context"
	"fmt"

	"github.com/shomali11/slacker/v2"
	"github.com/slack-go/slack/socketmode"
)

// Adding handlers to when the bot is connected, a default for when none of the commands match,
// adding default inner event handler when event type isn't message or app_mention

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.OnHello(func(event socketmode.Event) {
		log.Println("On Hello!")
		fmt.Println(event)
	})

	bot.OnConnected(func(event socketmode.Event) {
		log.Println("On Connected!")
		fmt.Println(event)
	})

	bot.OnConnecting(func(event socketmode.Event) {
		log.Println("On Connecting!")
		fmt.Println(event)
	})

	bot.OnConnectionError(func(event socketmode.Event) {
		log.Println("On Connection Error!")
		fmt.Println(event)
	})

	bot.OnDisconnected(func(event socketmode.Event) {
		log.Println("On Disconnected!")
		fmt.Println(event)
	})

	bot.UnsupportedCommandHandler(func(ctx *slacker.CommandContext) {
		ctx.Response().Reply("Say what?")
	})

	bot.UnsupportedEventHandler(func(event socketmode.Event) {
		fmt.Println(event)
	})

	definition := &slacker.CommandDefinition{
		Command:     "help",
		Description: "help!",
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("Your own help function...")
		},
	}

	bot.Help(definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
