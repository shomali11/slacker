package slacker

import (
	"github.com/shomali11/commander"
	"github.com/shomali11/proper"
)

// CommandDefinition structure contains definition of the bot command
type CommandDefinition struct {
	Usage       string
	Description string
	Examples    []string
	Middlewares []CommandMiddlewareHandler
	Handler     CommandHandler

	// HideHelp will hide this command definition from appearing in the `help` results.
	HideHelp bool
}

// newCommand creates a new bot command object
func newCommand(usage string, definition *CommandDefinition) Command {
	definition.Usage = usage

	return &command{
		definition: definition,
		cmd:        commander.NewCommand(usage),
	}
}

// Command interface
type Command interface {
	Definition() *CommandDefinition

	Match(string) (*proper.Properties, bool)
	Tokenize() []*commander.Token
	Handler(CommandContext, ...CommandMiddlewareHandler)
}

// command structure contains the bot's command, description and handler
type command struct {
	definition *CommandDefinition
	cmd        *commander.Command
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

// Handler executes the handler logic
func (c *command) Handler(ctx CommandContext, middlewares ...CommandMiddlewareHandler) {
	if c.definition == nil || c.definition.Handler == nil {
		return
	}

	handler := c.definition.Handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	handler(ctx)
}
