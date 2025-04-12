package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"log"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

var Conn *sql.DB
var Q *Queries

func Init() {
	ctx := context.Background()
	var err error
	Conn, err = sql.Open("sqlite", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := Conn.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	Q = New(Conn)
}
