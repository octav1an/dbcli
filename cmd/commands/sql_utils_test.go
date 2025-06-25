//go:build !integration

package commands

import (
	"database/sql"
	"reflect"
	"slices"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT);`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	_, err = db.Exec(`INSERT INTO test (name) VALUES ('Alice'), ('Bob');`)
	if err != nil {
		t.Fatalf("failed to insert data: %v", err)
	}
	return db
}

func Test_runQuery_Select(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cols, rows, err := runQuery(db, "SELECT id, name FROM test ORDER BY id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantCols := []string{"id", "name"}
	if !slices.Equal(cols, wantCols) {
		t.Errorf("cols are not equal")
	}

	wantRows := [][]string{
		{"1", "Alice"},
		{"2", "Bob"},
	}
	if !reflect.DeepEqual(rows, wantRows) {
		t.Errorf("rows = %v, want %v", rows, wantRows)
	}
}

func Test_runQuery_NoRows(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cols, rows, err := runQuery(db, "SELECT id, name FROM test WHERE id > 100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantCols := []string{"id", "name"}
	if !reflect.DeepEqual(cols, wantCols) {
		t.Errorf("cols = %v, want %v", cols, wantCols)
	}
	if len(rows) != 0 {
		t.Errorf("expected no rows, got %v", rows)
	}
}

func Test_runQuery_InvalidQuery(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, _, err := runQuery(db, "SELECT not_a_column FROM test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func Test_runExec_Insert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	affected, err := runExec(db, "INSERT INTO test (name) VALUES ('Charlie')")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Errorf("affected = %d, want 1", affected)
	}
	_, rows, err := runQuery(db, "SELECT * FROM test ORDER BY rowid DESC LIMIT 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantRows := [][]string{
		{"3", "Charlie"},
	}
	if !reflect.DeepEqual(rows, wantRows) {
		t.Errorf("cannot find inserted row")
	}
}

func Test_runExec_InvalidQuery(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := runExec(db, "INSERT INTO not_a_table (name) VALUES ('X')")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
