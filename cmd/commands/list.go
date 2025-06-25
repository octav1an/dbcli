package commands

import (
	"database/sql"
	"dbcli/internal/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var CmdListTables = &cobra.Command{
	Use:   "list",
	Short: "List all the tables names in the db",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateDb(config.DBPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", config.DBPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening database: %v\n", err)
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error querying tables: %v\n", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				fmt.Printf("error scanning the row %v\n", err)
				return
			}
			fmt.Println(tableName)
		}

		if err := rows.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error while iterating rows: %v\n", err)
			return
		}
	},
}
