package main

import (
	"context"
	"log"
	"os"

	"github.com/shomali11/slacker/v2"
)

// Defining an authorization middleware so that a command can only be executed by authorized users

var authorizedUserNames = []string{"shomali11"}

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	authorizedDefinitionByName := &slacker.CommandDefinition{
		Command:     "secret",
		Description: "Very secret stuff",
		Examples:    []string{"secret"},
		Middlewares: []slacker.CommandMiddlewareHandler{authorizationMiddleware()},
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("You are authorized!")
		},
	}

	bot.AddCommand(authorizedDefinitionByName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func authorizationMiddleware() slacker.CommandMiddlewareHandler {
	return func(next slacker.CommandHandler) slacker.CommandHandler {
		return func(ctx *slacker.CommandContext) {
			if contains(authorizedUserNames, ctx.Event().UserProfile.DisplayName) {
				next(ctx)
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
