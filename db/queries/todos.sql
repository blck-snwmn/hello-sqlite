-- name: CreateTodo :one
INSERT INTO todos (title, description, is_done) VALUES (?, ?, ?) RETURNING *;

-- name: ListTodos :many
SELECT * FROM todos;

-- name: GetTodo :one
SELECT * FROM todos WHERE id = ?;

-- name: UpdateTodo :exec
UPDATE todos SET title = ?, description = ?, is_done = ? WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;

-- name: UpdateTodoIsDone :exec
UPDATE todos SET is_done = ? WHERE id = ?;