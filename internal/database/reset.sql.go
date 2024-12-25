// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: reset.sql

package database

import (
	"context"
)

const deleteUsers = `-- name: DeleteUsers :many
DELETE FROM users
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

func (q *Queries) DeleteUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, deleteUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Email,
			&i.HashedPassword,
			&i.IsChirpyRed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
