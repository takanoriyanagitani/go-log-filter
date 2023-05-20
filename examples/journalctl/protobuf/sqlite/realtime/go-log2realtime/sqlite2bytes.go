package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func filename2sqlite(filename string) (*sql.DB, error) { return sql.Open("sqlite3", filename) }
