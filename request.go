package slacker

import (
	"context"

	"github.com/nlopes/slack"
	"github.com/shomali11/proper"
)

const (
	empty = ""
)

// NewRequest creates a new Request structure
func NewRequest(ctx context.Context, event *slack.MessageEvent, properties *proper.Properties) Request {
	return &request{ctx: ctx, event: event, properties: properties}
}

// Request interface that contains the Event received and parameters
type Request interface {
	Param(key string) string
	StringParam(key string, defaultValue string) string
	BooleanParam(key string, defaultValue bool) bool
	IntegerParam(key string, defaultValue int) int
	FloatParam(key string, defaultValue float64) float64
	Context() context.Context
	Event() *slack.MessageEvent
	Properties() *proper.Properties
}

// request contains the Event received and parameters
type request struct {
	ctx        context.Context
	event      *slack.MessageEvent
	properties *proper.Properties
}

// Param attempts to look up a string value by key. If not found, return the an empty string
func (r *request) Param(key string) string {
	return r.StringParam(key, empty)
}

// StringParam attempts to look up a string value by key. If not found, return the default string value
func (r *request) StringParam(key string, defaultValue string) string {
	return r.properties.StringParam(key, defaultValue)
}

// BooleanParam attempts to look up a boolean value by key. If not found, return the default boolean value
func (r *request) BooleanParam(key string, defaultValue bool) bool {
	return r.properties.BooleanParam(key, defaultValue)
}

// IntegerParam attempts to look up a integer value by key. If not found, return the default integer value
func (r *request) IntegerParam(key string, defaultValue int) int {
	return r.properties.IntegerParam(key, defaultValue)
}

// FloatParam attempts to look up a float value by key. If not found, return the default float value
func (r *request) FloatParam(key string, defaultValue float64) float64 {
	return r.properties.FloatParam(key, defaultValue)
}

// Context returns the current context of the request
func (r *request) Context() context.Context {
	return r.ctx
}

// Event returns the current event of the request
func (r *request) Event() *slack.MessageEvent {
	return r.event
}

// Properties returns the properties of the request
func (r *request) Properties() *proper.Properties {
	return r.properties
}
