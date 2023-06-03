package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/shomali11/slacker/v2"
)

// Scheduling messages

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := &slacker.CommandDefinition{
		Handler: func(ctx slacker.CommandContext) {
			now := time.Now()
			later := now.Add(time.Second * 20)

			ctx.Response().Reply("pong")
			ctx.Response().Reply("pong 20 seconds later", slacker.WithSchedule(later))
		},
	}

	bot.AddCommand("ping", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
