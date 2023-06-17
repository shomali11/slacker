package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shomali11/slacker/v2"
)

// Override the default event input cleaning function (to sanitize the messages received by Slacker)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot.SanitizeEventText(func(text string) string {
		fmt.Println("My slack bot does not like backticks!")
		return strings.ReplaceAll(text, "`", "")
	})

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "my-command",
		Handler: func(ctx slacker.CommandContext) {
			ctx.Response().Reply("it works!")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
