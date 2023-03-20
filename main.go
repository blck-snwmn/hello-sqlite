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
	errDoneUsage   = "Usage: done <id>"
)

type Todo db.Todo

func (t Todo) Print() {
	fmt.Printf("%d: %s (%s) - Done: %v\n", t.ID, t.Title, t.Description, t.IsDone)
}

type AddAction struct {
	q     *db.Queries
	title string
	desc  string
}

func (a *AddAction) Run(ctx context.Context) error {
	_, err := a.q.CreateTodo(ctx, db.CreateTodoParams{
		Title:       a.title,
		Description: a.desc,
		IsDone:      false,
	})
	if err != nil {
		return err
	}
	fmt.Println("Todo added")
	return nil
}

type ListAction struct {
	q *db.Queries
}

func (a *ListAction) Run(ctx context.Context) error {
	todos, err := a.q.ListTodos(ctx)
	if err != nil {
		return err
	}
	for _, todo := range todos {
		Todo(todo).Print()
	}
	return nil
}

type GetAction struct {
	q  *db.Queries
	id int64
}

func (a *GetAction) Run(ctx context.Context) error {
	todo, err := a.q.GetTodo(ctx, a.id)
	if err != nil {
		return err
	}
	Todo(todo).Print()
	return nil
}

type UpdateAction struct {
	q      *db.Queries
	id     int64
	title  string
	desc   string
	isDone bool
}

func (a *UpdateAction) Run(ctx context.Context) error {
	err := a.q.UpdateTodo(ctx, db.UpdateTodoParams{
		Title:       a.title,
		Description: a.desc,
		IsDone:      a.isDone,
		ID:          a.id,
	})
	if err != nil {
		return err
	}
	fmt.Println("Todo updated")
	return nil
}

type DeleteAction struct {
	q  *db.Queries
	id int64
}

func (a *DeleteAction) Run(ctx context.Context) error {
	err := a.q.DeleteTodo(ctx, a.id)
	if err != nil {
		return err
	}
	fmt.Println("Todo deleted")
	return nil
}

type DoneAction struct {
	q  *db.Queries
	id int64
}

func (a *DoneAction) Run(ctx context.Context) error {
	err := a.q.UpdateTodoIsDone(ctx, db.UpdateTodoIsDoneParams{
		ID:     a.id,
		IsDone: true,
	})
	if err != nil {
		return err
	}
	fmt.Println("Todo marked as done")
	return nil
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
	case "done":
		if len(commandArgs) < 1 {
			return nil, fmt.Errorf(errDoneUsage)
		}
		id, err := strconv.Atoi(commandArgs[0])
		if err != nil {
			return nil, err
		}
		return &DoneAction{q: q, id: int64(id)}, nil
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
	fmt.Println(" " + errDoneUsage)
}
