package data

import (
	"database/sql"
	"fmt"
	"strings"
)

func printQueryRows(columns []string, rows [][]string) {
	fmt.Printf("Columns: %s\n", strings.Join(columns, ", "))
	for idx, row := range rows {
		fmt.Printf("%d: %s\n", idx+1, strings.Join(row, ", "))
	}
}

func printExecResult(res sql.Result) {
	affected, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("error getting affected rows")
	}

	fmt.Printf("Rows affected: %d\n", affected)
}
