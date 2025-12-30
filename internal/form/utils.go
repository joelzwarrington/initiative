package form

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

// GetInt attempts to get an integer value from a huh form field.
// Returns the parsed integer and true if successful, or 0 and false if parsing failed.
func GetInt(form *huh.Form, key string) (int, bool) {
	rawValue := form.GetString(key)
	if rawValue == "" {
		return 0, false
	}

	parsed, err := strconv.Atoi(strings.TrimSpace(rawValue))
	if err != nil {
		return 0, false
	}

	return parsed, true
}

// GetIntWithDefault attempts to get an integer value from a huh form field.
// Returns the parsed integer if successful, or the default value if parsing failed.
func GetIntWithDefault(form *huh.Form, key string, defaultValue int) int {
	if value, ok := GetInt(form, key); ok {
		return value
	}
	return defaultValue
}
