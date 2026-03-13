-- name: CreatePin :one
INSERT INTO pins (name, blank_map_id, latitude, longitude, address, contact, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
RETURNING *;

-- name: GetAllPins :many
SELECT
    p.*,
    bm.name AS blank_map_name,
    bm.icon AS blank_map_icon
FROM pins p
LEFT JOIN blank_maps bm ON bm.id = p.blank_map_id
ORDER BY p.created_at DESC;

-- name: GetPinByID :one
SELECT
    p.*,
    bm.name AS blank_map_name,
    bm.icon AS blank_map_icon
FROM pins p
LEFT JOIN blank_maps bm ON bm.id = p.blank_map_id
WHERE p.id = $1;

-- name: GetPinsByBlankMapID :many
SELECT * FROM pins
WHERE blank_map_id = $1
ORDER BY created_at DESC;

-- name: GetPinsNearLocation :many
SELECT
    p.*,
    bm.name AS blank_map_name,
    bm.icon AS blank_map_icon,
    (
        6371 * acos(
            cos(radians($1)) * cos(radians(p.latitude)) *
            cos(radians(p.longitude) - radians($2)) +
            sin(radians($1)) * sin(radians(p.latitude))
        )
    ) AS distance_km
FROM pins p
LEFT JOIN blank_maps bm ON bm.id = p.blank_map_id
HAVING (
    6371 * acos(
        cos(radians($1)) * cos(radians(p.latitude)) *
        cos(radians(p.longitude) - radians($2)) +
        sin(radians($1)) * sin(radians(p.latitude))
    )
) < $3
ORDER BY distance_km;

-- name: UpdatePin :one
UPDATE pins
SET name = $2, address = $3, contact = $4, updated_at = NOW(), updated_by = $5
WHERE id = $1
RETURNING *;

-- name: DeletePin :exec
DELETE FROM pins
WHERE id = $1;
