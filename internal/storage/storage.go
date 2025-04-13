package storage

import (
	"context"
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

var Conn *sql.DB
var Q *Queries

func Init() error {
	ctx := context.Background()
	var err error
	Conn, err = sql.Open("sqlite", "./db.sqlite")
	if err != nil {
		return err
	}

	if _, err := Conn.ExecContext(ctx, ddl); err != nil {
		return err
	}

	Q = New(Conn)
	return nil
}
