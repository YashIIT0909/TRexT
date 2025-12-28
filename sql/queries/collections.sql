-- name: GetCollections :many
SELECT id, name, description 
FROM collections 
ORDER BY name;

-- name: GetCollectionByID :one
SELECT id, name, description 
FROM collections 
WHERE id = $1;

-- name: CreateCollection :one
INSERT INTO collections (name, description) 
VALUES ($1, $2)
RETURNING id, name, description;

-- name: UpdateCollection :exec
UPDATE collections 
SET name = $1, description = $2 
WHERE id = $3;

-- name: DeleteCollection :exec
DELETE FROM collections 
WHERE id = $1;
