package main

import (
	"context"
	"github.com/shomali11/slacker"
	"log"
)

func main() {
	bot := slacker.NewClient("<YOUR SLACK BOT TOKEN>")

	authorizedDefinition := &slacker.CommandDefinition{
		Description:           "Very secret stuff",
		AuthorizationRequired: true,
		// Either AuthorizationFunc OR AuthorizedUsers can be used to grant access
		// They are OR'ed together.
		AuthorizationFunc: func(request slacker.Request) { return true },
		AuthorizedUsers:   []string{},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
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
