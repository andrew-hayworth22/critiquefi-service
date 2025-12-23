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
func (a *AuthPG) CreateUser(ctx context.Context, usr types.User) (int64, error) {
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

	err := a.db.QueryRow(ctx, q,
		usr.Email,
		usr.DisplayName,
		usr.Name,
		usr.PasswordHash,
		usr.IsAdmin,
		usr.IsActive,
	).Scan(&id)

	return id, err
}

func (a *AuthPG) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
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

	row := a.db.QueryRow(ctx, q, email)
	var usr types.User
	if err := row.Scan(&usr.ID, &usr.Email, &usr.DisplayName, &usr.Name, &usr.PasswordHash, &usr.IsAdmin, &usr.IsActive, &usr.CreatedAt, &usr.UpdatedAt, &usr.LastLogin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &usr, nil
}

// GetUserByID fetches a user by their ID
func (a *AuthPG) GetUserByID(ctx context.Context, id int64) (*types.User, error) {
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
WHERE id = $1`

	row := a.db.QueryRow(ctx, q, id)
	var usr types.User
	if err := row.Scan(&usr.ID, &usr.Email, &usr.DisplayName, &usr.Name, &usr.PasswordHash, &usr.IsAdmin, &usr.IsActive, &usr.CreatedAt, &usr.UpdatedAt, &usr.LastLogin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &usr, nil
}

func (a *AuthPG) CreateRefreshToken(ctx context.Context, refreshToken *types.RefreshToken) error {
	const q = `
INSERT INTO refresh_tokens (
	token,
	user_id,
    user_agent,
	expires_at,
	created_at
) VALUES ($1, $2, $3, $4, $5)`

	_, err := a.db.Exec(ctx, q, refreshToken.Token, refreshToken.UserID, refreshToken.UserAgent, refreshToken.ExpiresAt, refreshToken.CreatedAt)
	return err
}

func (a *AuthPG) GetRefreshToken(ctx context.Context, token string) (*types.RefreshToken, error) {
	const q = `SELECT token, user_id, user_agent, revoked, expires_at, created_at FROM refresh_tokens WHERE token = $1`

	row := a.db.QueryRow(ctx, q, token)
	var rt types.RefreshToken
	if err := row.Scan(&rt.Token, &rt.UserID, &rt.UserAgent, &rt.Revoked, &rt.ExpiresAt, &rt.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (a *AuthPG) RevokeRefreshToken(ctx context.Context, token string) error {
	const q = `UPDATE refresh_tokens SET revoked = true WHERE token = $1`

	_, err := a.db.Exec(ctx, q, token)
	return err
}
