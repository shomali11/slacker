package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shomali11/commander"
	"github.com/shomali11/proper"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"), slacker.WithDebug(true))
	bot.CustomCommand(func(usage string, definition *slacker.CommandDefinition) slacker.BotCommand {
		return &cmd{
			usage:      usage,
			definition: definition,
			command:    commander.NewCommand(fmt.Sprintf("custom-prefix %s", usage)),
		}
	})

	// Invoked by `custom-prefix ping`
	bot.Command("ping", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			_ = response.Reply("it works!")
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

type cmd struct {
	usage      string
	definition *slacker.CommandDefinition
	command    *commander.Command
}

func (c *cmd) Usage() string {
	return c.usage
}

func (c *cmd) Definition() *slacker.CommandDefinition {
	return c.definition
}

func (c *cmd) Match(text string) (*proper.Properties, bool) {
	return c.command.Match(text)
}

func (c *cmd) Tokenize() []*commander.Token {
	return c.command.Tokenize()
}

func (c *cmd) Execute(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	log.Printf("Executing command [%s] invoked by %s", c.usage, botCtx.Event().User)
	c.definition.Handler(botCtx, request, response)
}

func (c *cmd) Interactive(*slacker.Slacker, *socketmode.Event, *slack.InteractionCallback, *socketmode.Request) {
}
