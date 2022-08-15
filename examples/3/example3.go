package main

import (
	"context"
	"log"
	"os"

	"github.com/broxgit/slacker"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.Command("echo {word}", &slacker.CommandDefinition{
		Description: "Echo a word!",
		Examples:    []string{"echo hello"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			word := request.Param("word")
			response.Reply(word)
		},
	})

	bot.Command("say <sentence>", &slacker.CommandDefinition{
		Description: "Say a sentence!",
		Examples:    []string{"say hello there everyone!"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			sentence := request.Param("sentence")
			response.Reply(sentence)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
