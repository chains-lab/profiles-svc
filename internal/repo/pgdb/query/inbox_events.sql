-- name: CreateInboxEvent :one
INSERT INTO inbox_events (
    topic, event_type, event_version, key, payload
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetPendingInboxEvents :many
SELECT * FROM inbox_events
WHERE status = 'pending' AND next_retry_at <= (now() AT TIME ZONE 'UTC')
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkInboxEventsProcessed :exec
UPDATE inbox_events
SET status = 'processed',
    processed_at = $2
WHERE id = ANY($1::uuid[]);

-- name: DelayInboxEvents :exec
UPDATE inbox_events
SET attempts = attempts + 1,
    next_retry_at = $2
WHERE id = ANY($1::uuid[]);