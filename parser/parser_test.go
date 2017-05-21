package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBooleanParam(t *testing.T) {
	emptyParameters := make(map[string]string)

	parameters := make(map[string]string)
	parameters["boolean"] = "true"
	parameters["bad"] = "bad"

	assert.Equal(t, BooleanParam("boolean", emptyParameters, false), false)
	assert.Equal(t, BooleanParam("boolean", parameters, false), true)
	assert.Equal(t, BooleanParam("bad", parameters, false), false)
}

func TestFloatParam(t *testing.T) {
	emptyParameters := make(map[string]string)

	parameters := make(map[string]string)
	parameters["float"] = "1.2"
	parameters["bad"] = "bad"

	assert.Equal(t, FloatParam("float", emptyParameters, 0), float64(0))
	assert.Equal(t, FloatParam("float", parameters, 0), float64(1.2))
	assert.Equal(t, FloatParam("bad", parameters, 0), float64(0))
}

func TestIntegerParam(t *testing.T) {
	emptyParameters := make(map[string]string)

	parameters := make(map[string]string)
	parameters["integer"] = "11"
	parameters["bad"] = "bad"

	assert.Equal(t, IntegerParam("integer", emptyParameters, 0), 0)
	assert.Equal(t, IntegerParam("integer", parameters, 0), 11)
	assert.Equal(t, IntegerParam("bad", parameters, 0), 0)
}

func TestStringParam(t *testing.T) {
	emptyParameters := make(map[string]string)

	parameters := make(map[string]string)
	parameters["string"] = "value"

	assert.Equal(t, StringParam("string", emptyParameters, ""), "")
	assert.Equal(t, StringParam("string", parameters, ""), "value")
}
