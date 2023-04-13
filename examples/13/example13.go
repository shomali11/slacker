package main

import (
	"log"
	"os"

	"context"
	"fmt"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.Init(func() {
		log.Println("Connected!")
	})

	bot.DefaultCommand(func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Say what?")
	})

	bot.DefaultEvent(func(event interface{}) {
		fmt.Println(event)
	})

	bot.DefaultInnerEvent(func(ctx context.Context, evt interface{}, request *socketmode.Request) {
		fmt.Printf("Handling inner event: %s", evt)
	})

	definition := &slacker.CommandDefinition{
		Description: "help!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("Your own help function...")
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
