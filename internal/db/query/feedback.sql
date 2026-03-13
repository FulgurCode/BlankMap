-- name: CreateFeedback :one
INSERT INTO feedback (pin_id, user_id, rating, review)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFeedbackByPinID :many
SELECT
    f.*,
    u.name AS user_name
FROM feedback f
JOIN users u ON u.id = f.user_id
WHERE f.pin_id = $1
ORDER BY f.created_at DESC;

-- name: GetFeedbackByID :one
SELECT * FROM feedback
WHERE id = $1;

-- name: UpdateFeedback :one
UPDATE feedback
SET rating = $2, review = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteFeedback :exec
DELETE FROM feedback
WHERE id = $1;
