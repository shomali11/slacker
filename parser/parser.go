package parser

import "strconv"

// StringParam attempts to look up a string value by key. If not found, return the default string value
func StringParam(key string, parameters map[string]string, defaultValue string) string {
	value, ok := parameters[key]
	if !ok {
		return defaultValue
	}
	return value
}

// BooleanParam attempts to look up a boolean value by key. If not found, return the default boolean value
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

// IntegerParam attempts to look up a integer value by key. If not found, return the default integer value
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

// FloatParam attempts to look up a float value by key. If not found, return the default float value
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
