package utils

import (
	"fmt"
)

// ParseSourceAndDestination parses the command arguments to extract source and destination.
func ParseSourceAndDestination(args []string) (string, string, error) {
	if len(args) < 2 {
		return "", "", fmt.Errorf("not enough arguments")
	}

	if len(args) == 2 {
		return args[0], args[1], nil
	}

	if len(args) == 3 && args[1] == "to" {
		return args[0], args[2], nil
	}

	return "", "", fmt.Errorf("invalid command format")
}
