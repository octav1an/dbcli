package data

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

func runSql(db *sql.DB, query string) error {
	q := strings.ToLower(strings.TrimSpace(query))
	// TODO: a delete query can start with "with"
	if strings.HasPrefix(q, "select") || strings.HasPrefix(q, "pragma") || strings.HasPrefix(q, "with") {
		cols, rows, err := runQuery(db, query)
		if err != nil {
			return fmt.Errorf("sql query error: %w", err)
		}
		printQueryRows(cols, rows)
	} else {
		res, err := runExec(db, query)
		if err != nil {
			return fmt.Errorf("sql exec error: %w", err)
		}
		printExecResult(res)
	}
	return nil
}

func runQuery(db *sql.DB, query string) ([]string, [][]string, error) {
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

func runExec(db *sql.DB, query string) (int64, error) {
	res, err := db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("error running sql string")
	}

	return res.RowsAffected()
}
