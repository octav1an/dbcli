package commands

import "testing"

func TestQueryBuilder2(t *testing.T) {
	var tests = []struct {
		table     string
		column    string
		start     *int
		end       *int
		wantQuery string
	}{
		{"test_table", "test_column", nil, nil, "SELECT test_column FROM test_table"},
		{"test_table", "", nil, nil, "SELECT * FROM test_table"},                                    // if column is omitted, query all columns
		{"test_table", "", nil, intPtr(10), "SELECT * FROM test_table LIMIT 10"},                    // select first 10 rows
		{"test_table", "", intPtr(10), nil, "SELECT * FROM test_table LIMIT -1 OFFSET 10"},          // select from n to end rows
		{"test_table", "", nil, intPtr(-5), "SELECT * FROM test_table ORDER BY ROWID DESC LIMIT 5"}, // select last 5 rows
		{"test_table", "", intPtr(6), intPtr(10), "SELECT * FROM test_table LIMIT 4 OFFSET 6"},      // select range from a to b
	}

	for _, tt := range tests {
		ans := queryBuilder(tt.table, tt.column, tt.start, tt.end)
		if ans != tt.wantQuery {
			t.Errorf("sql query: %s; want sql query: %s", ans, tt.wantQuery)
		}
	}
}
