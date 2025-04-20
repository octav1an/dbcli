package data

import (
	"dbcli/internal/config"
	"errors"
	"fmt"
	"os"
	"regexp"
)

func validateDb(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("database file does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("could not access database file: %v", err)
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}
	return nil
}

func validateColumnInput(input string) error {
	if input == "" {
		return nil // will select all columns
	}

	var validColumnPattern = regexp.MustCompile(`^[\w\d_]*(, [\w\d_]*)*$`)
	if !validColumnPattern.MatchString(input) {
		return errors.New("invalid column format: must be comma-separated words without spaces or special characters")
	}
	return nil
}

func vlog(format string, a ...interface{}) {
	if config.Verbose {
		fmt.Printf(format+"\n", a...)
	}
}

// func getColumns(db *sql.DB, tableName string) ([]string, error) {
// 	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s);", tableName))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var columns []string
// 	for rows.Next() {
// 		var cid int
// 		var name, ctype, notnull, pk string
// 		var dflt_value sql.NullString
// 		err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
// 		if err != nil {
// 			return nil, err
// 		}

// 		columns = append(columns, name)
// 	}

// 	return columns, nil
// }
