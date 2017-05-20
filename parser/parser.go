package parser

import "strconv"

func StringParam(key string, parameters map[string]string, defaultValue string) string {
	value, ok := parameters[key]
	if !ok {
		return defaultValue
	}
	return value
}

func BooleanParam(key string, parameters map[string]string, defaultValue bool) bool {
	value, ok := parameters[key]
	if !ok {
		return defaultValue
	}

	integerValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return integerValue
}

func IntegerParam(key string, parameters map[string]string, defaultValue int) int {
	value, ok := parameters[key]
	if !ok {
		return defaultValue
	}

	integerValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return integerValue
}

func FloatParam(key string, parameters map[string]string, defaultValue float64) float64 {
	value, ok := parameters[key]
	if !ok {
		return defaultValue
	}

	integerValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return integerValue
}
