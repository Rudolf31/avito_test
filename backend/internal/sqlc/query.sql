-- name: CreateTeam :one
INSERT INTO team (team_name) VALUES ($1)
RETURNING *;

-- name: CreateUser :one
INSERT INTO "user" (id, username, team_id, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM "user" WHERE id = $1;

-- name: SetIsActive :exec
UPDATE "user" SET is_active = $2 WHERE id = $1;

-- name: GetTeamNameByUserID :one
SELECT t.team_name
FROM team t
JOIN "user" u ON u.team_id = t.id
WHERE u.id = $1;

-- name: GetUserWithTeamNameByID :one
SELECT u.id, u.username, t.team_name, u.is_active
FROM "user" u
JOIN team t ON u.team_id = t.id
WHERE u.id = $1;

-- name: GetTeamByID :one
SELECT * FROM team WHERE id = $1;

-- name: GetTeamByName :one
SELECT * FROM team WHERE team_name = $1;

-- name: UpdateUser :one
UPDATE "user"
SET username = $2, team_id = $3, is_active = $4
WHERE id = $1
RETURNING *;

-- name: UpsertUser :one
INSERT INTO "user" (id, username, team_id, is_active)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE SET
    username = EXCLUDED.username,
    team_id = EXCLUDED.team_id,
    is_active = EXCLUDED.is_active
RETURNING *;

-- name: GetUsersByTeamID :many
SELECT * FROM "user" WHERE team_id = $1;

-- name: CreatePullRequest :one
INSERT INTO pull_request (id, pull_request_name, author_id, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPullRequestByID :one
SELECT * FROM pull_request WHERE id = $1;

-- name: CreateReview :one
INSERT INTO review (id, user_id, pull_request_id, reviewed)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetReviewByPRAndUser :one
SELECT * FROM review WHERE pull_request_id = $1 AND user_id = $2;

-- name: GetReviewsByPullRequestID :many
SELECT * FROM review WHERE pull_request_id = $1;

-- name: UpdateReviewUser :one
UPDATE review
SET user_id = $2
WHERE id = $1
RETURNING *;

-- name: UpdatePullRequestStatus :one
UPDATE pull_request
SET status = $2
    , merged_at = $3
WHERE id = $1
RETURNING *;

-- name: GetReviewsByUserID :many
SELECT * FROM review WHERE user_id = $1;