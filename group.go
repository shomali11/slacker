package slacker

import (
	"fmt"
	"strings"
)

type Group interface {
	AddCommand(usage string, definition *CommandDefinition)
	PrependCommand(usage string, definition *CommandDefinition)
	AddMiddleware(middleware MiddlewareHandler)

	GetPrefix() string
	GetCommands() []Command
	GetMiddlewares() []MiddlewareHandler
}

func newGroup(prefix string) Group {
	return &group{prefix: prefix}
}

type group struct {
	prefix      string
	middlewares []MiddlewareHandler
	commands    []Command
}

// AddMiddleware define a new middleware and append it to the list of group middlewares
func (g *group) AddMiddleware(middleware MiddlewareHandler) {
	g.middlewares = append(g.middlewares, middleware)
}

// AddCommand define a new command and append it to the list of group bot commands
func (g *group) AddCommand(usage string, definition *CommandDefinition) {
	fullUsage := strings.TrimSpace(fmt.Sprintf("%s %s", g.prefix, usage))
	g.commands = append(g.commands, newCommand(fullUsage, definition))
}

// PrependCommand define a new command and prepend it to the list of group bot commands
func (g *group) PrependCommand(usage string, definition *CommandDefinition) {
	fullUsage := strings.TrimSpace(fmt.Sprintf("%s %s", g.prefix, usage))
	g.commands = append([]Command{newCommand(fullUsage, definition)}, g.commands...)
}

// GetPrefix returns the group's prefix
func (g *group) GetPrefix() string {
	return g.prefix
}

// GetCommands returns Commands
func (g *group) GetCommands() []Command {
	return g.commands
}

// GetMiddlewares returns Middlewares
func (g *group) GetMiddlewares() []MiddlewareHandler {
	return g.middlewares
}
