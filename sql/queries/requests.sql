-- name: GetAllRequests :many
SELECT id, name, url, method, headers, body, COALESCE(collection_id, 1) as collection_id
FROM requests 
ORDER BY name;

-- name: GetRequestByID :one
SELECT id, name, url, method, headers, body, COALESCE(collection_id, 1) as collection_id
FROM requests 
WHERE id = $1;

-- name: GetRequestsByCollectionID :many
SELECT id, name, url, method, headers, body, COALESCE(collection_id, 1) as collection_id
FROM requests 
WHERE collection_id = $1
ORDER BY name;

-- name: CreateRequest :one
INSERT INTO requests (name, url, method, headers, body, collection_id) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, url, method, headers, body, collection_id;

-- name: UpdateRequest :exec
UPDATE requests 
SET name = $1, url = $2, method = $3, headers = $4, body = $5, collection_id = $6 
WHERE id = $7;

-- name: DeleteRequest :exec
DELETE FROM requests 
WHERE id = $1;
