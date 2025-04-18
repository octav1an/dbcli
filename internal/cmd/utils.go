package cmd

import (
	"fmt"
	"os"
)

func validateDb(path string) error {
	info, err := os.Stat(path)
	fmt.Println(path)
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
