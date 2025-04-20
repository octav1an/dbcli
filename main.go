package main

import (
	"dbcli/cmd"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
