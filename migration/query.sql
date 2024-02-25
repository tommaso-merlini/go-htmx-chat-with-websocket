-- name: GetUserByID :one
SELECT * FROM "users"
WHERE id = $1;

-- name: GetUserByAuthID :one
SELECT * FROM "users"
WHERE Authid = $1;

-- name: GetMessages :many
SELECT * FROM "messages";

-- name: CreateUser :one
INSERT INTO "users" (name, Authid, email) VALUES ($1, $2, $3) RETURNING id;

-- name: CreateMessage :one
INSERT INTO "messages" (from_id, from_authid, from_name, message) VALUES ($1, $2, $3, $4) RETURNING id;
