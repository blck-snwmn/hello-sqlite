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

func addTodo(ctx context.Context, q *db.Queries, title, desc string) error {
	_, err := q.CreateTodo(ctx, db.CreateTodoParams{
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

func listTodos(ctx context.Context, q *db.Queries) error {
	todos, err := q.ListTodos(ctx)
	if err != nil {
		return err
	}
	for _, todo := range todos {
		printTodo(todo)
	}
	return nil
}

func getTodo(ctx context.Context, q *db.Queries, id int64) error {
	todo, err := q.GetTodo(ctx, id)
	if err != nil {
		return err
	}
	printTodo(todo)
	return nil
}

func updateTodo(ctx context.Context, q *db.Queries, id int64, title, desc string, isDone bool) error {
	err := q.UpdateTodo(ctx, db.UpdateTodoParams{
		Title:       title,
		Description: desc,
		IsDone:      isDone,
		ID:          id,
	})
	if err != nil {
		return err
	}
	fmt.Println("Todo updated")
	return nil
}

func deleteTodo(ctx context.Context, q *db.Queries, id int64) error {
	err := q.DeleteTodo(ctx, id)
	if err != nil {
		return err
	}
	fmt.Println("Todo deleted")
	return nil
}

func printTodo(todo db.Todo) {
	fmt.Printf("%d: %s (%s) - Done: %v\n", todo.ID, todo.Title, todo.Description, todo.IsDone)
}

func validateArgs(args []string) error {
	command := args[1]

	switch command {
	case "add":
		if len(args) < 4 {
			return fmt.Errorf("Usage: add <title> <description>")
		}
	case "get", "delete":
		if len(args) < 3 {
			return fmt.Errorf("Usage: %s <id>", command)
		}
	case "update":
		if len(args) < 6 {
			return fmt.Errorf("Usage: update <id> <title> <description> <is_done>")
		}
	}

	return nil
}

func parseArgs(args []string) (string, []string) {
	return args[1], args[2:]
}

func main() {
	ctx := context.Background()

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

	command, args := parseArgs(os.Args)
	if err := validateArgs(os.Args); err != nil {
		fmt.Println(err)
		return
	}

	switch command {
	case "add":
		title := args[0]
		desc := args[1]

		if err := addTodo(ctx, q, title, desc); err != nil {
			log.Fatal(err)
		}
	case "list":
		if err := listTodos(ctx, q); err != nil {
			log.Fatal(err)
		}
	case "get":
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		if err := getTodo(ctx, q, int64(id)); err != nil {
			log.Fatal(err)
		}
	case "update":
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		title := args[1]
		desc := args[2]
		isDone, err := strconv.ParseBool(args[3])
		if err != nil {
			log.Fatal(err)
		}
		if err := updateTodo(ctx, q, int64(id), title, desc, isDone); err != nil {
			log.Fatal(err)
		}
	case "delete":
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		if err := deleteTodo(ctx, q, int64(id)); err != nil {
			log.Fatal(err)
		}
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
