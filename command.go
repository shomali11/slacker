package slacker

import (
	"github.com/shomali11/commander"
	"github.com/shomali11/proper"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// CommandDefinition structure contains definition of the bot command
type CommandDefinition struct {
	Description       string
	Examples          []string
	BlockID           string
	Channels 		  []string
	AuthorizationFunc func(BotContext, Request) bool
	Handler           func(BotContext, Request, ResponseWriter)
	Interactive       func(InteractiveBotContext, *socketmode.Request, *slack.InteractionCallback)

	// HideHelp will hide this command definition from appearing in the `help` results.
	HideHelp bool
}

// NewCommand creates a new bot command object
func NewCommand(usage string, definition *CommandDefinition) Command {
	return &command{
		usage:      usage,
		definition: definition,
		cmd:        commander.NewCommand(usage),
	}
}

// Command interface
type Command interface {
	Usage() string
	Definition() *CommandDefinition

	Match(string) (*proper.Properties, bool)
	Tokenize() []*commander.Token
	Execute(BotContext, Request, ResponseWriter)
	Interactive(InteractiveBotContext, *socketmode.Request, *slack.InteractionCallback)
}

// command structure contains the bot's command, description and handler
type command struct {
	usage      string
	definition *CommandDefinition
	cmd        *commander.Command
}

// Usage returns the command usage
func (c *command) Usage() string {
	return c.usage
}

// Definition returns the command definition
func (c *command) Definition() *CommandDefinition {
	return c.definition
}

// Match determines whether the bot should respond based on the text received
func (c *command) Match(text string) (*proper.Properties, bool) {
	return c.cmd.Match(text)
}

// Tokenize returns the command format's tokens
func (c *command) Tokenize() []*commander.Token {
	return c.cmd.Tokenize()
}

// Execute executes the handler logic
func (c *command) Execute(botCtx BotContext, request Request, response ResponseWriter) {
	if c.definition == nil || c.definition.Handler == nil {
		return
	}
	c.definition.Handler(botCtx, request, response)
}

// Interactive executes the interactive logic
func (c *command) Interactive(botContext InteractiveBotContext, request *socketmode.Request, callback *slack.InteractionCallback) {
	if c.definition == nil || c.definition.Interactive == nil {
		return
	}
	c.definition.Interactive(botContext, request, callback)
}
