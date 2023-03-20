package hellosqlite

import (
	"context"
	"fmt"
	"strconv"

	"github.com/blck-snwmn/hello-sqlite/db"
)

type Todo db.Todo

func (t Todo) Print() {
	fmt.Printf("%d: %s (%s) - Done: %v\n", t.ID, t.Title, t.Description, t.IsDone)
}

type Action interface {
	Run(ctx context.Context) error
}

func NewAction(ctx context.Context, args []string, q *db.Queries) (Action, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("Usage: <command> <args>")
	}

	command := args[1]
	commandArgs := args[2:]

	switch command {
	case "add":
		if len(commandArgs) < 2 {
			return nil, fmt.Errorf("Usage: add <title> <description>")
		}
		title := commandArgs[0]
		desc := commandArgs[1]
		return &AddAction{q: q, title: title, desc: desc}, nil
	case "list":
		return &ListAction{q: q}, nil
	case "get":
		if len(commandArgs) < 1 {
			return nil, fmt.Errorf("Usage: get <id>")
		}
		id, err := strconv.Atoi(commandArgs[0])
		if err != nil {
			return nil, err
		}
		return &GetAction{q: q, id: int64(id)}, nil
	case "update":
		if len(commandArgs) < 4 {
			return nil, fmt.Errorf("Usage: update <id> <title> <description> <is_done>")
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
			return nil, fmt.Errorf("Usage: delete <id>")
		}
		id, err := strconv.Atoi(commandArgs[0])
		if err != nil {
			return nil, err
		}
		return &DeleteAction{q: q, id: int64(id)}, nil
	case "done":
		if len(commandArgs) < 1 {
			return nil, fmt.Errorf("Usage: done <id>")
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
