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

type Todo db.Todo

func (t Todo) Print() {
	fmt.Printf("%d: %s (%s) - Done: %v\n", t.ID, t.Title, t.Description, t.IsDone)
}

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
		Todo(todo).Print()
	}
	return nil
}

func getTodo(ctx context.Context, q *db.Queries, id int64) error {
	todo, err := q.GetTodo(ctx, id)
	if err != nil {
		return err
	}
	Todo(todo).Print()
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

type AddAction struct {
	q    *db.Queries
	args []string
}

func (a *AddAction) Run(ctx context.Context) error {
	title := a.args[0]
	desc := a.args[1]

	return addTodo(ctx, a.q, title, desc)
}

type ListAction struct {
	q *db.Queries
}

func (a *ListAction) Run(ctx context.Context) error {
	return listTodos(ctx, a.q)
}

type GetAction struct {
	q    *db.Queries
	args []string
}

func (a *GetAction) Run(ctx context.Context) error {
	id, err := strconv.Atoi(a.args[0])
	if err != nil {
		return err
	}

	return getTodo(ctx, a.q, int64(id))
}

type UpdateAction struct {
	q    *db.Queries
	args []string
}

func (a *UpdateAction) Run(ctx context.Context) error {
	id, err := strconv.Atoi(a.args[0])
	if err != nil {
		return err
	}
	title := a.args[1]
	desc := a.args[2]
	isDone, err := strconv.ParseBool(a.args[3])
	if err != nil {
		return err
	}

	return updateTodo(ctx, a.q, int64(id), title, desc, isDone)
}

type DeleteAction struct {
	q    *db.Queries
	args []string
}

func (a *DeleteAction) Run(ctx context.Context) error {
	id, err := strconv.Atoi(a.args[0])
	if err != nil {
		return err
	}

	return deleteTodo(ctx, a.q, int64(id))
}

type Action interface {
	Run(ctx context.Context) error
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

	var action Action
	switch command {
	case "add":
		action = &AddAction{q: q, args: args}
	case "list":
		action = &ListAction{q: q}
	case "get":
		action = &GetAction{q: q, args: args}
	case "update":
		action = &UpdateAction{q: q, args: args}
	case "delete":
		action = &DeleteAction{q: q, args: args}
	default:
		usage()
		return
	}
	if err := action.Run(ctx); err != nil {
		log.Fatal(err)
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
