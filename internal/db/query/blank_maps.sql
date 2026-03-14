-- name: CreateBlankMap :one
INSERT INTO blank_maps (name, description, icon, created_by, updated_by)
VALUES ($1, $2, $3, $4, $4)
RETURNING *;

-- name: GetAllBlankMaps :many
SELECT 
  bm.*,
  COUNT(p.id) AS pin_count
FROM blank_maps bm
LEFT JOIN pins p ON p.blank_map_id = bm.id
GROUP BY bm.id
ORDER BY bm.name;

-- name: GetBlankMapByID :one
SELECT 
  bm.*,
  COUNT(p.id) AS pin_count
FROM blank_maps bm
LEFT JOIN pins p ON p.blank_map_id = bm.id
WHERE bm.id = $1
GROUP BY bm.id;

-- name: UpdateBlankMap :one
UPDATE blank_maps
SET name = $2, description = $3, icon = $4, updated_at = NOW(), updated_by = $5
WHERE id = $1
RETURNING *;

-- name: DeleteBlankMap :exec
DELETE FROM blank_maps
WHERE id = $1;

-- name: GetNoOfPins :one
SELECT COUNT(*) FROM pins WHERE blank_map_id = $1;
