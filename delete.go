package hellosqlite

import (
	"context"
	"fmt"

	"github.com/blck-snwmn/hello-sqlite/db"
)

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
