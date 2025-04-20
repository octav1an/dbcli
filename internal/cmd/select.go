package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var dbPath string
var tableName string
var columnName string

func init() {
	cmdSelect.PersistentFlags().StringVarP(&dbPath, "db", "d", "", "Path to SQLite DB file")
	cmdSelect.MarkPersistentFlagRequired("db")
	cmdSelect.PersistentFlags().StringVarP(&tableName, "table", "t", "", "Table name")
	cmdSelect.MarkPersistentFlagRequired("table")
	cmdSelect.Flags().StringVarP(&columnName, "column", "c", "",
		`Column name(s) to select (comma-separated for multiple).
Examples:
  --column file
  --column "file,time"`)

	cmdSelect.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := validateDb(dbPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := validateColumnInput(columnName); err != nil {
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
		vlog("Database: %s", dbPath)

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		query := queryBuilder(tableName, columnName)
		vlog("Query string: %s", query)

		cols, rows, err := getTableData(db, query)
		if err != nil {
			fmt.Printf("error getting data %v", err)
		}

		outputInConsole(cols, rows)
	},
}

func queryBuilder(tableName string, columns string) string {
	selector := "*"
	if columns != "" {
		// TODO: Validate column or multiple columns
		selector = columns
	}

	return fmt.Sprintf("SELECT %s FROM %s", selector, tableName)
}

func getTableData(db *sql.DB, query string) ([]string, [][]string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var rowsData [][]string
	colsCount := len(cols)

	for rows.Next() {
		values := make([]interface{}, colsCount)
		for i := range values {
			values[i] = new(interface{})
		}

		err := rows.Scan(values...)
		if err != nil {
			return nil, nil, fmt.Errorf("error scanning row, %v\n", err)
		}

		var row []string
		for _, val := range values {
			val_s := fmt.Sprintf("%v", *(val.(*interface{})))
			row = append(row, val_s)
		}
		rowsData = append(rowsData, row)
	}
	return cols, rowsData, nil
}

func outputInConsole(columns []string, rows [][]string) {
	fmt.Printf("Columns: %s", strings.Join(columns, ", "))
	fmt.Println()
	for _, row := range rows {
		fmt.Printf("%s", strings.Join(row, ", "))
		fmt.Println()
	}
}
