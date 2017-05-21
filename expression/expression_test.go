package expression

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	isMatch, parameters := Match("", "ping")
	assert.False(t, isMatch)
	assert.Empty(t, parameters)

	isMatch, parameters = Match("", "")
	assert.False(t, isMatch)
	assert.Empty(t, parameters)

	isMatch, parameters = Match("ping", "ping")
	assert.True(t, isMatch)
	assert.Empty(t, parameters)

	isMatch, parameters = Match("ping", "pong")
	assert.False(t, isMatch)
	assert.Empty(t, parameters)

	isMatch, parameters = Match("echo <word>", "echo")
	assert.True(t, isMatch)
	assert.Empty(t, parameters)

	isMatch, parameters = Match("echo <word>", "echo hey")
	assert.True(t, isMatch)
	assert.Equal(t, parameters["word"], "hey")

	isMatch, parameters = Match("repeat <word> <number>", "repeat hey 5")
	assert.True(t, isMatch)
	assert.Equal(t, parameters["word"], "hey")
	assert.Equal(t, parameters["number"], "5")
}

func TestIsParameter(t *testing.T) {
	assert.True(t, IsParameter("<value>"))
	assert.True(t, IsParameter("<123>"))
	assert.True(t, IsParameter("<value123>"))
	assert.False(t, IsParameter("value>"))
	assert.False(t, IsParameter("<value"))
	assert.False(t, IsParameter("value"))
}
