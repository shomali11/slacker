package expression

import (
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

// Match takes in the command and the text received, attempts to find the pattern and extract the parameters
func Match(command string, text string) (bool, map[string]string) {
	parameters := make(map[string]string)
	pattern := extractPattern(command)
	if len(pattern) == 0 {
		return false, parameters
	}

	compiledExpression, err := regexp.Compile(pattern)
	if err != nil {
		return false, parameters
	}

	result := strings.TrimSpace(compiledExpression.FindString(text))
	if len(result) == 0 {
		return false, parameters
	}

	commandTokens := strings.Split(command, space)
	resultTokens := strings.Split(result, space)

	valueRegex := regexp.MustCompile(parameterPattern)
	for i, resultToken := range resultTokens {
		commandToken := commandTokens[i]
		isValue := valueRegex.MatchString(commandToken)
		if !isValue {
			continue
		}

		parameters[commandToken[1:len(commandToken)-1]] = resultToken
	}
	return true, parameters
}

// IsParameter determines whether a string value satisfies the parameter pattern
func IsParameter(text string) bool {
	valueRegex := regexp.MustCompile(parameterPattern)
	return valueRegex.MatchString(text)
}

func extractPattern(command string) string {
	command = strings.TrimSpace(command)
	tokens := strings.Split(command, space)
	if len(tokens) == 0 {
		return empty
	}

	pattern := empty
	for _, token := range tokens {
		isMatch := IsParameter(token)
		if isMatch {
			pattern += wordPattern
		} else {
			pattern += token
		}
		pattern += spacePattern
	}
	return ignoreCase + pattern
}
