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

const (
	errAddUsage    = "Usage: add <title> <description>"
	errGetUsage    = "Usage: get <id>"
	errUpdateUsage = "Usage: update <id> <title> <description> <is_done>"
	errDeleteUsage = "Usage: delete <id>"
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

type AddAction struct {
	q     *db.Queries
	title string
	desc  string
}

func (a *AddAction) Run(ctx context.Context) error {
	return addTodo(ctx, a.q, a.title, a.desc)
}

type ListAction struct {
	q *db.Queries
}

func (a *ListAction) Run(ctx context.Context) error {
	return listTodos(ctx, a.q)
}

type GetAction struct {
	q  *db.Queries
	id int64
}

func (a *GetAction) Run(ctx context.Context) error {
	return getTodo(ctx, a.q, a.id)
}

type UpdateAction struct {
	q      *db.Queries
	id     int64
	title  string
	desc   string
	isDone bool
}

func (a *UpdateAction) Run(ctx context.Context) error {
	return updateTodo(ctx, a.q, a.id, a.title, a.desc, a.isDone)
}

type DeleteAction struct {
	q  *db.Queries
	id int64
}

func (a *DeleteAction) Run(ctx context.Context) error {
	return deleteTodo(ctx, a.q, a.id)
}

type Action interface {
	Run(ctx context.Context) error
}

func newAction(ctx context.Context, args []string, q *db.Queries) (Action, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("Usage: <command> <args>")
	}

	command := args[1]
	commandArgs := args[2:]

	switch command {
	case "add":
		if len(commandArgs) < 2 {
			return nil, fmt.Errorf(errAddUsage)
		}
		title := commandArgs[0]
		desc := commandArgs[1]
		return &AddAction{q: q, title: title, desc: desc}, nil
	case "list":
		return &ListAction{q: q}, nil
	case "get":
		if len(commandArgs) < 1 {
			return nil, fmt.Errorf(errGetUsage)
		}
		id, err := strconv.Atoi(commandArgs[0])
		if err != nil {
			return nil, err
		}
		return &GetAction{q: q, id: int64(id)}, nil
	case "update":
		if len(commandArgs) < 4 {
			return nil, fmt.Errorf(errUpdateUsage)
		}
		id, err := strconv.Atoi(commandArgs[0])
		if err != nil {
			return nil, err
		}
		isDone, err := strconv.ParseBool(commandArgs[3])
		if err != nil {
			return nil, err
		}
		return &UpdateAction{q: q, id: int64(id), title: commandArgs[1], desc: commandArgs[2], isDone: isDone}, nil
	case "delete":
		if len(commandArgs) < 1 {
			return nil, fmt.Errorf(errDeleteUsage)
		}
		id, err := strconv.Atoi(commandArgs[0])
		if err != nil {
			return nil, err
		}
		return &DeleteAction{q: q, id: int64(id)}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
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

	action, err := newAction(ctx, os.Args, q)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := action.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println(" " + errAddUsage)
	fmt.Println(" list")
	fmt.Println(" " + errGetUsage)
	fmt.Println(" " + errUpdateUsage)
	fmt.Println(" " + errDeleteUsage)
}
