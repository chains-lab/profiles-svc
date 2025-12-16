-- name: CreateProfile :one
INSERT INTO profiles (account_id, username)
VALUES ($1, $2)
RETURNING *;

-- name: GetProfileByAccountID :one
SELECT * FROM profiles
WHERE account_id = $1;

-- name: GetProfileByUsername :one
SELECT * FROM profiles
WHERE username = $1;

-- name: UpdateProfileUsername :one
UPDATE profiles
SET username = $2, updated_at = NOW()
WHERE account_id = $1
RETURNING *;

-- name: UpdateProfileOfficial :one
UPDATE profiles
SET official = $2, updated_at = NOW()
WHERE account_id = $1
RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
SET
    pseudonym = CASE
        WHEN $2 IS NULL THEN pseudonym
        WHEN $2 = ''    THEN NULL
        ELSE $2
    END,
    description = CASE
        WHEN $3 IS NULL THEN description
        WHEN $3 = ''    THEN NULL
        ELSE $3
    END,
    avatar = CASE
        WHEN $4 IS NULL THEN avatar
        WHEN $4 = ''    THEN NULL
        ELSE $4
    END,
    updated_at = now()
WHERE account_id = $1
RETURNING *;

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE account_id = $1;

-- name: ListProfilesByUsername :many
SELECT * FROM profiles
WHERE
    (sqlc.arg(prefix)::text = '' OR username ILIKE ('%' || sqlc.arg(q)::text || '%'))
    LIMIT sqlc.arg(limit_)::int
OFFSET sqlc.arg(offset_)::int;

-- name: ListProfilesByPseudonym :many
SELECT * FROM profiles
WHERE
    (sqlc.arg(prefix)::text = '' OR pseudonym ILIKE ('%' || sqlc.arg(q)::text || '%'))
    LIMIT sqlc.arg(limit_)::int
OFFSET sqlc.arg(offset_)::int;