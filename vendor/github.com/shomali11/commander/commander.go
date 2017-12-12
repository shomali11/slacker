package commander

import (
	"github.com/shomali11/proper"
	"regexp"
	"strings"
)

const (
	escapeCharacter      = "\\"
	ignoreCase           = "(?i)"
	parameterPattern     = "<\\S+>"
	spacePattern         = "\\s+"
	optionalSpacePattern = "\\s*"
	inputPattern         = "(\\S+)?"
	preCommandPattern    = "(\\s|^)"
	postCommandPattern   = "(\\s|$)"
)

var (
	regexCharacters = []string{"\\", "(", ")", "{", "}", "[", "]", "?", ".", "+", "|", "^", "$"}
)

// NewCommand creates a new Command object from the format passed in
func NewCommand(format string) *Command {
	tokens := tokenize(format)
	expression := compile(tokens)
	return &Command{tokens: tokens, expression: expression}
}

// Token represents the Token object
type Token struct {
	Word        string
	IsParameter bool
}

// Command represents the Command object
type Command struct {
	tokens     []*Token
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
	resultTokens := strings.Fields(result)
	for i, resultToken := range resultTokens {
		commandToken := c.tokens[i]
		if !commandToken.IsParameter {
			continue
		}

		parameters[commandToken.Word] = resultToken
	}
	return proper.NewProperties(parameters), true
}

// Tokenize returns Command info as tokens
func (c *Command) Tokenize() []*Token {
	return c.tokens
}

func escape(text string) string {
	for _, character := range regexCharacters {
		text = strings.Replace(text, character, escapeCharacter+character, -1)
	}
	return text
}

func tokenize(format string) []*Token {
	parameterRegex := regexp.MustCompile(parameterPattern)
	words := strings.Fields(format)
	tokens := make([]*Token, len(words))
	for i, word := range words {
		if parameterRegex.MatchString(word) {
			tokens[i] = &Token{Word: word[1 : len(word)-1], IsParameter: true}
		} else {
			tokens[i] = &Token{Word: word, IsParameter: false}
		}
	}
	return tokens
}

func compile(tokens []*Token) *regexp.Regexp {
	if len(tokens) == 0 {
		return nil
	}

	pattern := preCommandPattern
	if tokens[0].IsParameter {
		pattern += inputPattern
	} else {
		pattern += escape(tokens[0].Word)
	}

	for index := 1; index < len(tokens); index++ {
		previousToken := tokens[index-1]
		currentToken := tokens[index]

		if !previousToken.IsParameter && !currentToken.IsParameter {
			pattern += spacePattern + escape(currentToken.Word)
		} else if previousToken.IsParameter && currentToken.IsParameter {
			pattern += optionalSpacePattern + inputPattern
		} else if previousToken.IsParameter && !currentToken.IsParameter {
			pattern += optionalSpacePattern + escape(currentToken.Word)
		} else {
			pattern += optionalSpacePattern + inputPattern
		}
	}
	pattern += postCommandPattern
	return regexp.MustCompile(ignoreCase + pattern)
}
