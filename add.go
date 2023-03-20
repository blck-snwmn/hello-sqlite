package hellosqlite

import (
	"context"
	"fmt"

	"github.com/blck-snwmn/hello-sqlite/db"
)

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
