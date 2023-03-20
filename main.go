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

func parseArgs(args []string) (string, []string) {
	return args[1], args[2:]
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

func newAction(ctx context.Context, command string, args []string, q *db.Queries) (Action, error) {
	switch command {
	case "add":
		if len(args) < 2 {
			return nil, fmt.Errorf("Usage: add <title> <description>")
		}
		title := args[0]
		desc := args[1]
		return &AddAction{q: q, title: title, desc: desc}, nil
	case "list":
		return &ListAction{q: q}, nil
	case "get":
		if len(args) < 1 {
			return nil, fmt.Errorf("Usage: get <id>")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}
		return &GetAction{q: q, id: int64(id)}, nil
	case "update":
		if len(args) < 4 {
			return nil, fmt.Errorf("Usage: update <id> <title> <description> <is_done>")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}
		isDone, err := strconv.ParseBool(args[3])
		if err != nil {
			return nil, err
		}
		return &UpdateAction{q: q, id: int64(id), title: args[1], desc: args[2], isDone: isDone}, nil
	case "delete":
		if len(args) < 1 {
			return nil, fmt.Errorf("Usage: delete <id>")
		}
		id, err := strconv.Atoi(args[0])
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

	command, args := parseArgs(os.Args)
	action, err := newAction(ctx, command, args, q)
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
	fmt.Println(" add <title> <description>")
	fmt.Println(" list")
	fmt.Println(" get <id>")
	fmt.Println(" update <id> <title> <description> <is_done>")
	fmt.Println(" delete <id>")
}
