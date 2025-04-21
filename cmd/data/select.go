package data

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var tableName string
var columnName string
var selectRange string

func init() {
	cmdSelect.PersistentFlags().StringVarP(&tableName, "table", "t", "", "Table name")
	cmdSelect.MarkPersistentFlagRequired("table")
	cmdSelect.Flags().StringVarP(&columnName, "column", "c", "",
		`Column name(s) to select (comma-separated for multiple). Examples:
  --column file
  --column "file,time"`)
	cmdSelect.Flags().StringVarP(&selectRange, "range", "r", "",
		`Selection range. Examples:
  --range "5:10" - Select a range from 5th to 10th
  --range ":10" - Select first 10 entries
  --range "10:" - Select a range from 10th to the last
  --range ":-10" - Select last 10 entries`)

	cmdSelect.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := validateColumnInput(columnName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if selectRange != "" { // range is optional
			if err := validateRangeInput(selectRange); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}

		return nil
	}

	CmdData.AddCommand(cmdSelect)
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

		start, end, err := parseRangeInput(selectRange)
		if err != nil {
			fmt.Printf("error parsing range %v", err)
			return
		}

		query := queryBuilder(tableName, columnName, start, end)
		vlog("Query string: %s", query)

		cols, rows, err := getTableData(db, query)
		if err != nil {
			fmt.Printf("error getting data %v", err)
			return
		}

		outputToConsole(cols, rows)
	},
}

func queryBuilder(tableName string, columns string, start *int, end *int) string {
	selector := "*"
	if columns != "" {
		selector = columns
	}
	query := fmt.Sprintf("SELECT %s FROM %s", selector, tableName)

	// Case 1: select first 'n' entries
	if end != nil && *end >= 0 && start == nil {
		query = fmt.Sprintf("%s LIMIT %d", query, *end)
		return query
	}

	// Case 2: select from 'n' to 'end'
	if start != nil && end == nil {
		query = fmt.Sprintf("%s LIMIT -1 OFFSET %d", query, *start)
		return query
	}

	// Case 3: select last 'n' entries (in reverse order)
	if start == nil && end != nil && *end < 0 {
		query = fmt.Sprintf("%s ORDER BY ROWID DESC LIMIT %d", query, intAbs(*end))
		return query
	}

	// Case 4: select range (start to end)
	if start != nil && end != nil && *end > 0 {
		query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, *end-*start, *start)
		return query
	}

	// TODO: handle case where the start > end

	return query
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
			return nil, nil, errors.New("error scanning row")
		}

		var row []string
		for _, val := range values {
			val_s := fmt.Sprintf("%v", *(val.(*interface{})))
			row = append(row, val_s)
		}
		rowsData = append(rowsData, row)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error while iterating rows: %v\n", err)
		return nil, nil, err
	}

	return cols, rowsData, nil
}

func outputToConsole(columns []string, rows [][]string) {
	fmt.Printf("Columns: %s\n", strings.Join(columns, ", "))
	for idx, row := range rows {
		fmt.Printf("%d: %s\n", idx+1, strings.Join(row, ", "))
	}
}
