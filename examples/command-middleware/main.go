package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

// Defining an authorization middleware so that a command can only be executed by authorized users

var authorizedUserNames = []string{"shomali11"}

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	authorizedDefinitionByName := &slacker.CommandDefinition{
		Description: "Very secret stuff",
		Examples:    []string{"secret"},
		Middlewares: []slacker.MiddlewareHandler{authorizationMiddleware()},
		Handler: func(botCtx slacker.CommandContext) {
			botCtx.Response().Reply("You are authorized!")
		},
	}

	bot.AddCommand("secret", authorizedDefinitionByName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func authorizationMiddleware() slacker.MiddlewareHandler {
	return func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(botCtx slacker.CommandContext) {
			if contains(authorizedUserNames, botCtx.Event().UserProfile.DisplayName) {
				next(botCtx)
			}
		}
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
