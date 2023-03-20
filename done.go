package hellosqlite

import (
	"context"
	"fmt"

	"github.com/blck-snwmn/hello-sqlite/db"
)

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
