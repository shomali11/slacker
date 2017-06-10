# properties [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/properties)](https://goreportcard.com/report/github.com/shomali11/properties) [![GoDoc](https://godoc.org/github.com/shomali11/properties?status.svg)](https://godoc.org/github.com/shomali11/properties) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A `map[string]string` decorator offering a collection of helpful functions to extract the values in different types

# Examples

```go
parameters := make(map[string]string)
parameters["boolean"] = "true"
parameters["float"] = "1.2"
parameters["integer"] = "11"
parameters["string"] = "value"
	
properties := NewProperties(parameters)
	
assert.Equal(t, properties.BooleanParam("boolean", false), true)
assert.Equal(t, properties.FloatParam("float", 0), float64(1.2))
assert.Equal(t, properties.IntegerParam("integer", 0), 11)
assert.Equal(t, properties.StringParam("string", ""), "value")
```