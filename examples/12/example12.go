package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	authorizedUserIds := []string{"<User ID>"}
	authorizedUserNames := []string{"shomali11"}

	authorizedDefinitionById := &slacker.CommandDefinition{
		Description: "Very secret stuff",
		Examples:    []string{"secret-id"},
		AuthorizationFunc: func(botCtx slacker.BotContext, request slacker.Request) bool {
			return contains(authorizedUserIds, botCtx.Event().UserID)
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("You are authorized!")
		},
	}

	authorizedDefinitionByName := &slacker.CommandDefinition{
		Description: "Very secret stuff",
		Examples:    []string{"secret-name"},
		AuthorizationFunc: func(botCtx slacker.BotContext, request slacker.Request) bool {
			return contains(authorizedUserNames, botCtx.Event().UserProfile.DisplayName)
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("You are authorized!")
		},
	}

	bot.Command("secret-id", authorizedDefinitionById)
	bot.Command("secret-name", authorizedDefinitionByName)

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
