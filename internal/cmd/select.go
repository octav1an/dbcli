package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dbPath string

func init() {
	cmdSelect.PersistentFlags().StringVarP(&dbPath, "db", "d", "", "Path to SQLite DB file")
	cmdSelect.MarkPersistentFlagRequired("db")

	cmdSelect.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := validateDb(dbPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		return nil
	}
}

var cmdSelect = &cobra.Command{
	Use:   "select",
	Short: "Query the sql db file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Select from ", dbPath)

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		columns, err := getColumns(db, "file_log")
		if err != nil {
			fmt.Printf("error %v", err)
		}
		// for _, column := range columns {
		// 	fmt.Printf(" %s ", column)
		// }

		rows, err := db.Query("SELECT * FROM file_log")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		values := make([]interface{}, len(columns))
		for rows.Next() {
			for i := range values {
				values[i] = new(interface{})
			}

			err = rows.Scan(values...)
			if err != nil {
				fmt.Printf("error, %v\n", err)
			}

			for i, column := range columns {
				fmt.Printf(" %s:%v ", column, *(values[i].(*interface{})))
			}
			fmt.Print("\n")
		}
	},
}

func getColumns(db *sql.DB, tableName string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s);", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name, ctype, notnull, pk string
		var dflt_value sql.NullString
		err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
		if err != nil {
			return nil, err
		}

		columns = append(columns, name)
	}

	return columns, nil
}
