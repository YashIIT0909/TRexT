-- name: GetHistory :many
SELECT id, url, method, status_code, duration_ms, timestamp 
FROM history 
ORDER BY timestamp DESC 
LIMIT $1;

-- name: GetHistoryByID :one
SELECT id, url, method, status_code, duration_ms, timestamp 
FROM history 
WHERE id = $1;

-- name: AddToHistory :one
INSERT INTO history (url, method, status_code, duration_ms, timestamp) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id, url, method, status_code, duration_ms, timestamp;

-- name: DeleteHistoryEntry :exec
DELETE FROM history 
WHERE id = $1;

-- name: ClearHistory :exec
DELETE FROM history;
