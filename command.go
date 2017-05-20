package slacker

import "github.com/shomali11/slacker/expression"

func NewCommand(cmd string, description string, handler func(request *Request, response *Response)) *Command {
	return &Command{cmd: cmd, description: description, handler: handler}
}

type Command struct {
	cmd         string
	description string
	handler     func(request *Request, response *Response)
}

func (c *Command) Match(text string) (bool, map[string]string) {
	return expression.Match(c.cmd, text)
}

func (c *Command) Execute(request *Request, response *Response) {
	c.handler(request, response)
}
