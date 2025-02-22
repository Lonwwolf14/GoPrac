-- name: CreateTask :one
INSERT INTO tasks (title, detail, current_status) 
VALUES ($1, $2, $3) 
RETURNING *;

-- name: ListTasks :many
SELECT id, title, detail, current_status FROM tasks;


-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;

-- name: UpdateTask :exec
UPDATE tasks SET current_status = TRUE WHERE id = $1;

