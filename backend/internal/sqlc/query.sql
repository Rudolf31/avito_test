-- name: CreateTeam :one
INSERT INTO team (team_name) VALUES ('Backend Team')
RETURNING *;