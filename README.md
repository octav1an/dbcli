# dbcli

CLI tool to interact with sqlite database

## Features

- List all tables in a SQLite database.
- Query tables with optional filters.
- Run raw SQL queries.
- Lightweight and easy to use.

## Usage

### List Tables

```bash
go run main.go -d <dbfile> list
```

### Query Records

```bash
# Select all records from a table
go run main.go -d <dbfile> select -t <table>

# Select a single column
go run main.go -d <dbfile> select -t <table> -c <column>

# Select a row range (e.g., 6th to 10th rows)
go run main.go -d <dbfile> select -t <table> -r "5:10"

# Select from the 10th row to the end
go run main.go -d <dbfile> select -t <table> -r "10:"

# Select the last 10 rows
go run main.go -d <dbfile> select -t <table> -r ":-10"
```

### Run Raw SQL

```bash
go run main.go -d <dbfile> sql "SELECT \* FROM <table>;"
```

## Build

```
go build -o bin/dbcli
```

## Test

```
go test ./...
```
