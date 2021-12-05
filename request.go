package slacker

import (
	allot "github.com/sdslabs/allot/pkg"
)

const (
	empty = ""
)

// NewRequest creates a new Request structure
func NewRequest(botCtx BotContext, parameters []allot.Parameter, match allot.MatchInterface) Request {
	return &request{botCtx: botCtx, parameters: parameters, match: match}
}

// Request interface that contains the Event received and parameters
type Request interface {
	Param(key string) string
	StringParam(key string, defaultValue string) string
	IntegerParam(key string, defaultValue int) int
	Parameters() []allot.Parameter
}

// request contains the Event received and parameters
type request struct {
	botCtx     BotContext
	parameters []allot.Parameter
	match      allot.MatchInterface
}

// Param attempts to look up a string value by key. If not found, return the an empty string
func (r *request) Param(key string) string {
	return r.StringParam(key, empty)
}

// StringParam attempts to look up a string value by key. If not found, return the default string value
func (r *request) StringParam(key string, defaultValue string) string {
	re, err := r.match.String(key)
	if err != nil {
		return defaultValue
	}
	return re
}

// IntegerParam attempts to look up a integer value by key. If not found, return the default integer value
func (r *request) IntegerParam(key string, defaultValue int) int {
	re, err := r.match.Integer(key)
	if err != nil {
		return defaultValue
	}
	return re
}

// Parameters returns the Parameters of the request
func (r *request) Parameters() []allot.Parameter {
	return r.parameters
}
