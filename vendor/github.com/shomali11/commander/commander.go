package commander

import (
	"github.com/shomali11/proper"
	"regexp"
	"strings"
)

const (
	empty            = ""
	space            = " "
	ignoreCase       = "(?i)"
	parameterPattern = "<\\S+>"
	spacePattern     = "\\s*"
	wordPattern      = "(\\S+)?"
)

var parameterRegex *regexp.Regexp

func init() {
	parameterRegex = regexp.MustCompile(parameterPattern)
}

// NewCommand creates a new Command object from the format passed in
func NewCommand(format string) *Command {
	expression := compile(format)
	return &Command{format: format, expression: expression}
}

// Command represents the Command object
type Command struct {
	format     string
	expression *regexp.Regexp
}

// Match takes in the command and the text received, attempts to find the pattern and extract the parameters
func (c *Command) Match(text string) (*proper.Properties, bool) {
	if c.expression == nil {
		return nil, false
	}

	result := strings.TrimSpace(c.expression.FindString(text))
	if len(result) == 0 {
		return nil, false
	}

	parameters := make(map[string]string)
	commandTokens := strings.Split(c.format, space)
	resultTokens := strings.Split(result, space)

	for i, resultToken := range resultTokens {
		commandToken := commandTokens[i]
		if !IsParameter(commandToken) {
			continue
		}

		parameters[commandToken[1:len(commandToken)-1]] = resultToken
	}
	return proper.NewProperties(parameters), true
}

// IsParameter determines whether a string value satisfies the parameter pattern
func IsParameter(text string) bool {
	return parameterRegex.MatchString(text)
}

func compile(commandFormat string) *regexp.Regexp {
	commandFormat = strings.TrimSpace(commandFormat)
	tokens := strings.Split(commandFormat, space)
	pattern := empty
	for _, token := range tokens {
		if len(token) == 0 {
			continue
		}

		if IsParameter(token) {
			pattern += wordPattern
		} else {
			pattern += token
		}
		pattern += spacePattern
	}

	if len(pattern) == 0 {
		return nil
	}
	return regexp.MustCompile(ignoreCase + pattern)
}
