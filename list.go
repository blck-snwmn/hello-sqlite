package hellosqlite

import (
	"context"

	"github.com/blck-snwmn/hello-sqlite/db"
)

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
