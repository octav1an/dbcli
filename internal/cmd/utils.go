package cmd

import (
"errors"
	"fmt"
	"os"
"regexp"
)

func validateDb(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("dabase file does not exist: %s", path)
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

	var validColumnPattern = regexp.MustCompile(`^[\w\d_]*(, [\w\d_]*)*$`)
	if !validColumnPattern.MatchString(input) {
		return errors.New("invalid column format: must be comma-separated words without spaces or special characters")
	}
	return nil
}

func vlog(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format+"\n", a...)
	}
}

