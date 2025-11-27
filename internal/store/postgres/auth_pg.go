package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// authPg is the postgres implementation of the auth store
type AuthPG struct {
	db *pgxpool.Pool
}

// NewAuthPG creates a new auth store
func NewAuthPG(db *pgxpool.Pool) *AuthPG {
	return &AuthPG{db: db}
}

// Create creates a user record
func (u *AuthPG) CreateUser(ctx context.Context, usr types.User) (int64, error) {
	const q = `
INSERT INTO users (
	email,
	display_name,
	name,
	password_hash,
	is_admin,
	is_active
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`

	var id int64

	err := u.db.QueryRow(ctx, q,
		usr.Email,
		usr.DisplayName,
		usr.Name,
		usr.PasswordHash,
		usr.IsAdmin,
		usr.IsActive,
	).Scan(&id)

	return id, err
}

func (u *AuthPG) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	const q = `
SELECT
	id,
	email,
	display_name,
	name,
	password_hash,
	is_admin,
	is_active,
	created_at,
	updated_at,
	last_login
FROM users
WHERE email = $1`

	row := u.db.QueryRow(ctx, q, email)
	var usr types.User
	if err := row.Scan(&usr.ID, &usr.Email, &usr.DisplayName, &usr.Name, &usr.PasswordHash, &usr.IsAdmin, &usr.IsActive, &usr.CreatedAt, &usr.UpdatedAt, &usr.LastLogin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &usr, nil
}
