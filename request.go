package slacker

import (
	"github.com/shomali11/proper"
)

const (
	empty = ""
)

// newRequest creates a new Request structure
func newRequest(properties *proper.Properties) Request {
	return &request{properties: properties}
}

// Request interface that contains the Event received and parameters
type Request interface {
	Param(key string) string
	StringParam(key string, defaultValue string) string
	BooleanParam(key string, defaultValue bool) bool
	IntegerParam(key string, defaultValue int) int
	FloatParam(key string, defaultValue float64) float64
	Properties() *proper.Properties
}

// request contains the Event received and parameters
type request struct {
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

// Properties returns the properties of the request
func (r *request) Properties() *proper.Properties {
	return r.properties
}
