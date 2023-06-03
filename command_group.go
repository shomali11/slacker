package slacker

import (
	"fmt"
	"strings"
)

// CommandGroup a group of commands
type CommandGroup interface {
	AddCommand(usage string, definition *CommandDefinition)
	PrependCommand(usage string, definition *CommandDefinition)
	AddMiddleware(middleware CommandMiddlewareHandler)

	GetPrefix() string
	GetCommands() []Command
	GetMiddlewares() []CommandMiddlewareHandler
}

func newGroup(prefix string) CommandGroup {
	return &commandGroup{prefix: prefix}
}

type commandGroup struct {
	prefix      string
	middlewares []CommandMiddlewareHandler
	commands    []Command
}

// AddMiddleware define a new middleware and append it to the list of group middlewares
func (g *commandGroup) AddMiddleware(middleware CommandMiddlewareHandler) {
	g.middlewares = append(g.middlewares, middleware)
}

// AddCommand define a new command and append it to the list of group bot commands
func (g *commandGroup) AddCommand(usage string, definition *CommandDefinition) {
	fullUsage := strings.TrimSpace(fmt.Sprintf("%s %s", g.prefix, usage))
	g.commands = append(g.commands, newCommand(fullUsage, definition))
}

// PrependCommand define a new command and prepend it to the list of group bot commands
func (g *commandGroup) PrependCommand(usage string, definition *CommandDefinition) {
	fullUsage := strings.TrimSpace(fmt.Sprintf("%s %s", g.prefix, usage))
	g.commands = append([]Command{newCommand(fullUsage, definition)}, g.commands...)
}

// GetPrefix returns the group's prefix
func (g *commandGroup) GetPrefix() string {
	return g.prefix
}

// GetCommands returns Commands
func (g *commandGroup) GetCommands() []Command {
	return g.commands
}

// GetMiddlewares returns Middlewares
func (g *commandGroup) GetMiddlewares() []CommandMiddlewareHandler {
	return g.middlewares
}
