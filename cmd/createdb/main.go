package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/blck-snwmn/hello-sqlite/db"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conn, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	// create tables
	if _, err := conn.ExecContext(context.Background(), db.DDL); err != nil {
		log.Fatal(err)
	}
}
