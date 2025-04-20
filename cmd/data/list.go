package data

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listTables = &cobra.Command{
	Use:   "list",
	Short: "List all the tables names in the db",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", dbPath)
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

func init() {
	CmdData.AddCommand(listTables)
}
