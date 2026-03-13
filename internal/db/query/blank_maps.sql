-- name: CreateBlankMap :one
INSERT INTO blank_maps (name, description, icon, created_by, updated_by)
VALUES ($1, $2, $3, $4, $4)
RETURNING *;

-- name: GetAllBlankMaps :many
SELECT * FROM blank_maps
ORDER BY name;

-- name: GetBlankMapByID :one
SELECT * FROM blank_maps
WHERE id = $1;

-- name: UpdateBlankMap :one
UPDATE blank_maps
SET name = $2, description = $3, icon = $4, updated_at = NOW(), updated_by = $5
WHERE id = $1
RETURNING *;

-- name: DeleteBlankMap :exec
DELETE FROM blank_maps
WHERE id = $1;
