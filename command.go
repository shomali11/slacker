package slacker

import "github.com/shomali11/slacker/expression"

// NewCommand creates a new command structure
func NewCommand(cmd string, description string, handler func(request *Request, response *Response)) *Command {
	return &Command{cmd: cmd, description: description, handler: handler}
}

// Command structure contains the command, description and handler
type Command struct {
	cmd         string
	description string
	handler     func(request *Request, response *Response)
}

// Match determines whether the bot should respond based on the text received
func (c *Command) Match(text string) (bool, map[string]string) {
	return expression.Match(c.cmd, text)
}

// Execute executes the handler logic
func (c *Command) Execute(request *Request, response *Response) {
	c.handler(request, response)
}
