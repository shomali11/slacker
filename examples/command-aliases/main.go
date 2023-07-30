package main

import (
	"context"
	"log"
	//"os"

	"github.com/shomali11/slacker/v2"
)

// Defining a command with aliases

func main() {
	//bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	bot := slacker.NewClient("xoxb-13360094916-2243791173942-I6FwX8kkBNPvANakdfVWm8PF", "xapp-1-A027JMM1RV2-5373599545492-06fa58c972681603e29d44cac9c0f62aa895bc5ccc013256b180dee7961eab23")

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "echo {word}",
		Aliases: []string{
			"repeat {word}",
			"mimic {word}",
		},
		Description: "Echo a word!",
		Examples: []string{
			"echo hello",
			"repeat hello",
			"mimic hello",
		},
		Handler: func(ctx *slacker.CommandContext) {
			word := ctx.Request().Param("word")
			ctx.Response().Reply(word)
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
