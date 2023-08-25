package slacker

import (
	"regexp"

	"github.com/shomali11/proper"
)

func (s *Slacker) regexMatch(regex, text string) (*proper.Properties, bool) {

	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, false
	}

	values := re.FindStringSubmatch(text)
	if len(values) == 0 {
		return nil, false
	}
	valueIndex := 0
	keys := re.SubexpNames()
	parameters := make(map[string]string)
	for i := 1; i < len(keys) && valueIndex < len(values); i++ {
		if len(values[i]) == 0 {
			continue
		}
		parameters[keys[i]] = values[i]
		valueIndex++
	}

	return proper.NewProperties(parameters), re.MatchString(text)
}
