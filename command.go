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
	Example           string
	BlockID           string
	AuthorizationFunc func(botCtx BotContext, request Request) bool
	Handler           func(botCtx BotContext, request Request, response ResponseWriter)
	Interactive       func(*Slacker, *socketmode.Event, *slack.InteractionCallback, *socketmode.Request)

	// HideHelp will cause this command to not be shown when a user requests
	// help.
	HideHelp bool
}

// NewBotCommand creates a new bot command object
func NewBotCommand(usage string, definition *CommandDefinition) BotCommand {
	command := commander.NewCommand(usage)
	return &botCommand{
		usage:      usage,
		definition: definition,
		command:    command,
	}
}

// BotCommand interface
type BotCommand interface {
	Usage() string
	Definition() *CommandDefinition

	Match(text string) (*proper.Properties, bool)
	Tokenize() []*commander.Token
	Execute(botCtx BotContext, request Request, response ResponseWriter)
	Interactive(*Slacker, *socketmode.Event, *slack.InteractionCallback, *socketmode.Request)
}

// botCommand structure contains the bot's command, description and handler
type botCommand struct {
	usage      string
	definition *CommandDefinition
	command    *commander.Command
}

// Usage returns the command usage
func (c *botCommand) Usage() string {
	return c.usage
}

// Description returns the command description
func (c *botCommand) Definition() *CommandDefinition {
	return c.definition
}

// Match determines whether the bot should respond based on the text received
func (c *botCommand) Match(text string) (*proper.Properties, bool) {
	return c.command.Match(text)
}

// Tokenize returns the command format's tokens
func (c *botCommand) Tokenize() []*commander.Token {
	return c.command.Tokenize()
}

// Execute executes the handler logic
func (c *botCommand) Execute(botCtx BotContext, request Request, response ResponseWriter) {
	if c.definition == nil || c.definition.Handler == nil {
		return
	}
	c.definition.Handler(botCtx, request, response)
}

// Interactive executes the interactive logic
func (c *botCommand) Interactive(slacker *Slacker, evt *socketmode.Event, callback *slack.InteractionCallback, req *socketmode.Request) {
	if c.definition == nil || c.definition.Interactive == nil {
		return
	}
	c.definition.Interactive(slacker, evt, callback, req)
}
