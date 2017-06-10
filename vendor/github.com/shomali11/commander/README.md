# commander [![Go Report Card](https://goreportcard.com/badge/github.com/shomali11/commander)](https://goreportcard.com/report/github.com/shomali11/commander) [![GoDoc](https://godoc.org/github.com/shomali11/commander?status.svg)](https://godoc.org/github.com/shomali11/commander) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Command evaluator and parser

# Examples

```go
properties, isMatch = NewCommand("ping").Match("ping")
assert.True(t, isMatch)
assert.NotNil(t, properties)

properties, isMatch = NewCommand("repeat <word> <number>").Match("repeat hey 5")
assert.True(t, isMatch)
assert.Equal(t, properties.StringParam("word", ""), "hey")
assert.Equal(t, properties.IntegerParam("number", 0), 5)
```