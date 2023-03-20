package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	hellosqlite "github.com/blck-snwmn/hello-sqlite"
	"github.com/blck-snwmn/hello-sqlite/db"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()

	conn, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	q := db.New(conn)

	action, err := hellosqlite.NewAction(ctx, os.Args, q)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := action.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
