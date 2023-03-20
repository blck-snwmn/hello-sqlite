package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/blck-snwmn/hello-sqlite/db"
	_ "github.com/mattn/go-sqlite3"
)

func addTodo(q *db.Queries, title, desc string) error {
	_, err := q.CreateTodo(context.Background(), db.CreateTodoParams{
		Title:       title,
		Description: desc,
		IsDone:      false,
	})
	if err != nil {
		return err
	}
	fmt.Println("Todo added")
	return nil
}

func listTodos(q *db.Queries) error {
	todos, err := q.ListTodos(context.Background())
	if err != nil {
		return err
	}
	for _, todo := range todos {
		fmt.Printf("%d: %s (%s) - Done: %v\n", todo.ID, todo.Title, todo.Description, todo.IsDone)
	}
	return nil
}

func getTodo(q *db.Queries, id int64) error {
	todo, err := q.GetTodo(context.Background(), id)
	if err != nil {
		return err
	}
	fmt.Printf("%d: %s (%s) - Done: %v\n", todo.ID, todo.Title, todo.Description, todo.IsDone)
	return nil
}

func main() {
	conn, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	q := db.New(conn)

	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 4 {
			usage()
			return
		}
		title := os.Args[2]
		desc := os.Args[3]

		if err := addTodo(q, title, desc); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Todo added")
	case "list":
		if err := listTodos(q); err != nil {
			log.Fatal(err)
		}
	case "get":
		if len(os.Args) < 3 {
			usage()
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		if err := getTodo(q, int64(id)); err != nil {
			log.Fatal(err)
		}
	case "update":
		if len(os.Args) < 6 {
			usage()
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		title := os.Args[3]
		desc := os.Args[4]
		isDone, err := strconv.ParseBool(os.Args[5])
		if err != nil {
			log.Fatal(err)
		}
		err = q.UpdateTodo(context.Background(), db.UpdateTodoParams{
			Title:       title,
			Description: desc,
			IsDone:      isDone,
			ID:          int64(id),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Todo updated")
	case "delete":
		if len(os.Args) < 3 {
			usage()
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		err = q.DeleteTodo(context.Background(), int64(id))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Todo deleted")
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println(" add <title> <description>")
	fmt.Println(" list")
	fmt.Println(" get <id>")
	fmt.Println(" update <id> <title> <description> <is_done>")
	fmt.Println(" delete <id>")
}
