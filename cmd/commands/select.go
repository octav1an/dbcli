package commands

import (
	"database/sql"
	"dbcli/internal/config"
	"fmt"

	"github.com/spf13/cobra"
)

var tableName string
var columnName string
var selectRange string

func init() {
	CmdSelect.PersistentFlags().StringVarP(&tableName, "table", "t", "", "Table name")
	CmdSelect.MarkPersistentFlagRequired("table")
	CmdSelect.Flags().StringVarP(&columnName, "column", "c", "",
		`Column name(s) to select (comma-separated for multiple). Examples:
  --column file
  --column "file,time"`)
	CmdSelect.Flags().StringVarP(&selectRange, "range", "r", "",
		`Selection range. Examples:
  --range "5:10" - Select a range from 5th to 10th
  --range ":10" - Select first 10 entries
  --range "10:" - Select a range from 10th to the last
  --range ":-10" - Select last 10 entries`)

	CmdSelect.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := ValidateDb(config.DBPath); err != nil {
			return fmt.Errorf("invalid database: %w", err)
		}
		if err := validateColumnInput(columnName); err != nil {
			return fmt.Errorf("invalid column: %w", err)
		}

		if selectRange != "" { // range is optional
			if err := validateRangeInput(selectRange); err != nil {
				return fmt.Errorf("invalid range: %w", err)
			}
		}

		return nil
	}
}

var CmdSelect = &cobra.Command{
	Use:   "select",
	Short: "Query the sql db file",
	Run: func(cmd *cobra.Command, args []string) {
		vlog("Database: %s", config.DBPath)

		db, err := sql.Open("sqlite3", config.DBPath)
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

		cols, rows, err := runQuery(db, query)
		if err != nil {
			fmt.Printf("error getting data %v", err)
			return
		}

		printQueryRows(cols, rows)
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
