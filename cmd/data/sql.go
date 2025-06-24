package data

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	CmdData.AddCommand(cmdSqlSelect)
}

var cmdSqlSelect = &cobra.Command{
	Use:   "sql",
	Short: "Use sql query strings to interact sqlite db file",
	Long:  "This command takes a sql query string as a positional argument and prints the result",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires the query string")
		}
		// TODO: validate the query string!!
		if len(args[0]) > 0 {
			return nil
		}
		return fmt.Errorf("invalid query string %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		vlog("Database: %s", dbPath)
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		queryString := args[0]
		vlog(queryString)

		// TODO: Make a backup
		if err := runSql(db, queryString); err != nil {
			fmt.Fprint(os.Stderr, "error: %w", err)
			os.Exit(1)
		}
	},
}
