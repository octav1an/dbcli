package commands

import (
	"dbcli/internal/config"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ValidateDb(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("database file does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("could not access database file: %v", err)
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}
	return nil
}

func validateColumnInput(input string) error {
	if input == "" {
		return nil // will select all columns
	}

	validColumnPattern := regexp.MustCompile(`^\w*(, \w*)*$`)
	if !validColumnPattern.MatchString(input) {
		return errors.New("error: invalid column format: must be comma-separated words without spaces or special characters")
	}
	return nil
}

func validateRangeInput(input string) error {
	if input == "" {
		return nil // will select all rows
	}

	validRangePattern := regexp.MustCompile(`^\d*:-?\d*$`)
	if !validRangePattern.MatchString(input) {
		return errors.New("error: invalid range format")
	}
	return nil
}

// Returns the start and end value of a range input. If start or end is omitted, nil is returned
func parseRangeInput(input string) (*int, *int, error) {
	if input == "" {
		return nil, nil, nil
	}

	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return nil, nil, errors.New("error: invalid range format")
	}

	parsePart := func(s string) (*int, error) {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, nil
		}
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", s, err)
		}
		return &val, nil
	}

	start, err := parsePart(parts[0])
	if err != nil {
		return nil, nil, err
	}
	// Start value cannot be negative
	if start != nil && *start < 0 {
		return nil, nil, fmt.Errorf("start value cannot be negative: %d", *start)
	}

	end, err := parsePart(parts[1])
	if err != nil {
		return nil, nil, err
	}

	return start, end, nil
}

func vlog(format string, a ...interface{}) {
	if config.Verbose {
		fmt.Printf(format+"\n", a...)
	}
}

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
