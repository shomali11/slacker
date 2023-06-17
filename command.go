package slacker

import (
	"github.com/shomali11/commander"
	"github.com/shomali11/proper"
)

// CommandDefinition structure contains definition of the bot command
type CommandDefinition struct {
	Command     string
	Description string
	Examples    []string
	Middlewares []CommandMiddlewareHandler
	Handler     CommandHandler

	// HideHelp will hide this command definition from appearing in the `help` results.
	HideHelp bool
}

// newCommand creates a new bot command object
func newCommand(definition *CommandDefinition) Command {
	return &command{
		definition: definition,
		cmd:        commander.NewCommand(definition.Command),
	}
}

// Command interface
type Command interface {
	Definition() *CommandDefinition

	Match(string) (*proper.Properties, bool)
	Tokenize() []*commander.Token
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
