package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	authorizedUsers := []string{"<User ID>"}

	authorizedDefinition := &slacker.CommandDefinition{
		Description: "Very secret stuff",
		AuthorizationFunc: func(botCtx slacker.BotContext, request slacker.Request) bool {
			return contains(authorizedUsers, botCtx.Event().User)
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("You are authorized!")
		},
	}

	bot.Command("secret", authorizedDefinition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func contains(list []string, element string) bool {
	for _, value := range list {
		if value == element {
			return true
		}
	}
	return false
}
