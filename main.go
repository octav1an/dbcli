package main

import (
	"go-sqlite-cli/internal/cmd"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
