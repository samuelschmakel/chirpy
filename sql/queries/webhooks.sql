-- name: UpgradeToChirpyRed :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;