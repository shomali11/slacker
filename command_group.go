package slacker

import (
	"fmt"
	"strings"
)

// newGroup creates a new CommandGroup with a prefix
func newGroup(prefix string) *CommandGroup {
	return &CommandGroup{prefix: prefix}
}

// CommandGroup groups commands with a common prefix and middlewares
type CommandGroup struct {
	prefix      string
	middlewares []CommandMiddlewareHandler
	commands    []Command
}

// AddMiddleware define a new middleware and append it to the list of group middlewares
func (g *CommandGroup) AddMiddleware(middleware CommandMiddlewareHandler) {
	g.middlewares = append(g.middlewares, middleware)
}

// AddCommand define a new command and append it to the list of group bot commands
func (g *CommandGroup) AddCommand(definition *CommandDefinition) {
	definition.Command = strings.TrimSpace(fmt.Sprintf("%s %s", g.prefix, definition.Command))
	g.commands = append(g.commands, newCommand(definition))
}

// PrependCommand define a new command and prepend it to the list of group bot commands
func (g *CommandGroup) PrependCommand(definition *CommandDefinition) {
	definition.Command = strings.TrimSpace(fmt.Sprintf("%s %s", g.prefix, definition.Command))
	g.commands = append([]Command{newCommand(definition)}, g.commands...)
}

// GetPrefix returns the group's prefix
func (g *CommandGroup) GetPrefix() string {
	return g.prefix
}

// GetCommands returns Commands
func (g *CommandGroup) GetCommands() []Command {
	return g.commands
}

// GetMiddlewares returns Middlewares
func (g *CommandGroup) GetMiddlewares() []CommandMiddlewareHandler {
	return g.middlewares
}
