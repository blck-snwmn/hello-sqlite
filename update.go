package hellosqlite

import (
	"context"
	"fmt"

	"github.com/blck-snwmn/hello-sqlite/db"
)

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
