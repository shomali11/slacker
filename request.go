package slacker

import (
	"github.com/nlopes/slack"
	"github.com/shomali11/slacker/parser"
)

const (
	empty = ""
)

// NewRequest creates a new Request structure
func NewRequest(event *slack.MessageEvent, parameters map[string]string) *Request {
	return &Request{Event: event, parameters: parameters}
}

// Request contains the Event received and parameters
type Request struct {
	Event      *slack.MessageEvent
	parameters map[string]string
}

// Param attempts to look up a string value by key. If not found, return the an empty string
func (r *Request) Param(key string) string {
	return r.StringParam(key, empty)
}

// StringParam attempts to look up a string value by key. If not found, return the default string value
func (r *Request) StringParam(key string, defaultValue string) string {
	return parser.StringParam(key, r.parameters, defaultValue)
}

// BooleanParam attempts to look up a boolean value by key. If not found, return the default boolean value
func (r *Request) BooleanParam(key string, defaultValue bool) bool {
	return parser.BooleanParam(key, r.parameters, defaultValue)
}

// IntegerParam attempts to look up a integer value by key. If not found, return the default integer value
func (r *Request) IntegerParam(key string, defaultValue int) int {
	return parser.IntegerParam(key, r.parameters, defaultValue)
}

// FloatParam attempts to look up a float value by key. If not found, return the default float value
func (r *Request) FloatParam(key string, defaultValue float64) float64 {
	return parser.FloatParam(key, r.parameters, defaultValue)
}
