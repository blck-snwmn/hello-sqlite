package hellosqlite

import (
	"context"

	"github.com/blck-snwmn/hello-sqlite/db"
)

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
