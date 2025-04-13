package storage

import (
	"context"
	"database/sql"
	_ "embed"
	"os"

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

func SetupTestDB() {
	dbFile := "test.db"
	os.Remove(dbFile)

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		panic(err)
	}

	_, err = conn.ExecContext(context.Background(), ddl)
	if err != nil {
		panic(err)
	}

	Conn = conn
	Q = New(conn)
}

func TeardownTestDB() {
	if Conn != nil {
		Conn.Close()
	}
	os.Remove("test.db")
}
