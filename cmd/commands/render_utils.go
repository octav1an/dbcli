package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aquasecurity/table"
)

func printQueryRows(columns []string, rows [][]string) {
	t := table.New(os.Stdout)
	fullHeader := prependSlice(columns, "idx")
	t.SetHeaders(fullHeader...)

	for idx, row := range rows {
		fullRow := prependSlice(row, strconv.Itoa(idx))
		t.AddRow(fullRow...)
	}
	t.SetFooters(fullHeader...)

	t.Render()
}

func prependSlice(s []string, el string) []string {
	s = append(s, "")
	copy(s[1:], s)
	s[0] = el
	return s
}

func printExecResult(rowsAffected int64) {
	fmt.Printf("Rows affected: %d\n", rowsAffected)
}
