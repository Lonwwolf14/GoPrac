-- name: CreateTask :one
INSERT INTO Tasks (id, title, detail, current_status) 
VALUES ($1, $2, $3, $4) 
RETURNING *;
