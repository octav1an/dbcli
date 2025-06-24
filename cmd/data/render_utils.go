package data

import (
	"fmt"
	"strings"
)

func printQueryRows(columns []string, rows [][]string) {
	fmt.Printf("Columns: %s\n", strings.Join(columns, ", "))
	for idx, row := range rows {
		fmt.Printf("%d: %s\n", idx+1, strings.Join(row, ", "))
	}
}

func printExecResult(rowsAffected int64) {
	fmt.Printf("Rows affected: %d\n", rowsAffected)
}
