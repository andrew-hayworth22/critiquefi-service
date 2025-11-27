package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres/tools"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// userPg the postgres implementation of the user store
type userPG struct {
	db *pgxpool.Pool
}

// Create creates a user record
func (u *userPG) Create(ctx context.Context, usr types.User) (int64, error) {
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

// GetByID fetches a user by their ID
func (u *userPG) GetByID(ctx context.Context, id int64) (*types.User, error) {
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

	row := u.db.QueryRow(ctx, q, id)
	var usr types.User
	if err := row.Scan(&usr.ID, &usr.Email, &usr.DisplayName, &usr.Name, &usr.PasswordHash, &usr.IsAdmin, &usr.IsActive, &usr.CreatedAt, &usr.UpdatedAt, &usr.LastLogin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &usr, nil
}

// GetByEmail fetches a user record by their email
func (u *userPG) GetByEmail(ctx context.Context, email string) (*types.User, error) {
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

// Update updates a user record
func (u *userPG) Update(ctx context.Context, id int64, usr types.UserUpdate) error {
	sb := tools.NewSetBuilder()
	sb.SetIf(usr.Email != nil, "email", *usr.Email)
	sb.SetIf(usr.DisplayName != nil, "display_name", *usr.DisplayName)
	sb.SetIf(usr.Name != nil, "name", *usr.Name)
	sb.SetIf(usr.PasswordHash != nil, "password_hash", *usr.PasswordHash)
	sb.SetIf(usr.IsAdmin != nil, "is_admin", *usr.IsAdmin)
	sb.SetIf(usr.IsActive != nil, "is_active", *usr.IsActive)
	if sb.Empty() {
		return nil
	}

	q := `UPDATE users SET ` + sb.BuildSet() + `WHERE id = $1`
	args := append(sb.Args(), id)
	_, err := u.db.Exec(ctx, q, args...)
	return err
}
