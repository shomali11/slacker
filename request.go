package slacker

import (
	"github.com/nlopes/slack"
	"github.com/shomali11/slacker/parser"
)

const (
	empty = ""
)

func NewRequest(event *slack.MessageEvent, parameters map[string]string) *Request {
	return &Request{Event: event, parameters: parameters}
}

type Request struct {
	Event      *slack.MessageEvent
	parameters map[string]string
}

func (r *Request) Param(key string) string {
	return r.StringParam(key, empty)
}

func (r *Request) StringParam(key string, defaultValue string) string {
	return parser.StringParam(key, r.parameters, defaultValue)
}

func (r *Request) BooleanParam(key string, defaultValue bool) bool {
	return parser.BooleanParam(key, r.parameters, defaultValue)
}

func (r *Request) IntegerParam(key string, defaultValue int) int {
	return parser.IntegerParam(key, r.parameters, defaultValue)
}

func (r *Request) FloatParam(key string, defaultValue float64) float64 {
	return parser.FloatParam(key, r.parameters, defaultValue)
}
