package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/jmoiron/sqlx"
)

// AuthStore is the postgres implementation of the authbus store
type AuthStore struct {
	db *sqlx.DB
}

// NewAuthStore creates a new authbus store
func NewAuthStore(db *sqlx.DB) *AuthStore {
	return &AuthStore{db: db}
}

// userRow represents a user row in the database
type userRow struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	DisplayName  string    `db:"display_name"`
	Name         string    `db:"name"`
	PasswordHash string    `db:"password_hash"`
	IsAdmin      bool      `db:"is_admin"`
	LastLogin    time.Time `db:"last_login"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (usr userRow) toModel() models.User {
	return models.User{
		ID:           usr.ID,
		Email:        usr.Email,
		DisplayName:  usr.DisplayName,
		Name:         usr.Name,
		IsAdmin:      usr.IsAdmin,
		PasswordHash: usr.PasswordHash,
		IsActive:     usr.IsActive,
	}
}

// refreshTokenRow represents a refresh token row in the database
type refreshTokenRow struct {
	TokenHash  string     `db:"token_hash"`
	UserID     int64      `db:"user_id"`
	UserAgent  string     `db:"user_agent"`
	LastUsedAt *time.Time `db:"last_used_at"`
	ExpiresAt  time.Time  `db:"expires_at"`
	CreatedAt  time.Time  `db:"created_at"`
}

func (rt *refreshTokenRow) toModel() models.RefreshToken {
	return models.RefreshToken{
		TokenHash: rt.TokenHash,
		UserID:    rt.UserID,
		UserAgent: rt.UserAgent,
		ExpiresAt: rt.ExpiresAt,
		CreatedAt: rt.CreatedAt,
	}
}

// userFieldsTaken represents the result of a check for taken fields
type userFieldsTaken struct {
	EmailTaken       bool `db:"email_taken"`
	DisplayNameTaken bool `db:"display_name_taken"`
}

func (uft *userFieldsTaken) toModel() models.UserFieldsTaken {
	return models.UserFieldsTaken{
		EmailTaken:       uft.EmailTaken,
		DisplayNameTaken: uft.DisplayNameTaken,
	}
}

// CreateUser creates a user record and returns the new user's ID
func (s *AuthStore) CreateUser(ctx context.Context, user models.NewUser) (int64, error) {
	const q = `
		INSERT INTO users (
			email,
			display_name,
			name,
			password_hash
		) VALUES (:email, :display_name, :name, :password_hash)
		RETURNING id;
	`

	rows, err := s.db.NamedQueryContext(ctx, q, map[string]any{
		"email":         user.Email,
		"display_name":  user.DisplayName,
		"name":          user.Name,
		"password_hash": user.PasswordHash,
	})
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", mapError(err))
	}
	defer closeRows(rows)
	if !rows.Next() {
		return 0, store.ErrNotFound
	}

	var id int64
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("error scanning user ID: %w", mapError(err))
	}
	return id, nil
}

// GetUserByID fetches a user by their ID
func (s *AuthStore) GetUserByID(ctx context.Context, id int64) (models.User, error) {
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
		WHERE id = :id
	`

	rows, err := s.db.NamedQueryContext(ctx, q, map[string]any{
		"id": id,
	})
	if err != nil {
		return models.User{}, store.ErrNotFound
	}
	defer closeRows(rows)
	if !rows.Next() {
		return models.User{}, store.ErrNotFound
	}

	var row userRow
	if err := rows.StructScan(&row); err != nil {
		return models.User{}, mapError(err)
	}

	return row.toModel(), nil
}

// GetUserByEmail fetches a user by their email
func (s *AuthStore) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
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
		WHERE email = :email
	`

	rows, err := s.db.NamedQueryContext(ctx, q, map[string]any{
		"email": email,
	})
	if err != nil {
		return models.User{}, mapError(err)
	}
	defer closeRows(rows)
	if !rows.Next() {
		return models.User{}, store.ErrNotFound
	}

	var row userRow
	if err := rows.StructScan(&row); err != nil {
		return models.User{}, mapError(err)
	}

	return row.toModel(), nil
}

// CheckTakenUserFields checks if another user takes the given fields
func (s *AuthStore) CheckTakenUserFields(ctx context.Context, request models.NewUserRequest) (models.UserFieldsTaken, error) {
	const q = `
		SELECT
			EXISTS (SELECT 1 FROM users WHERE email = :email) AS email_taken,
			EXISTS (SELECT 1 FROM users WHERE display_name = :display_name) AS display_name_taken
	`

	rows, err := s.db.NamedQueryContext(ctx, q, map[string]any{
		"email":        request.Email,
		"display_name": request.DisplayName,
	})
	if err != nil {
		return models.UserFieldsTaken{}, mapError(err)
	}
	defer closeRows(rows)

	if !rows.Next() {
		return models.UserFieldsTaken{}, store.ErrNotFound
	}

	var uft userFieldsTaken
	if err := rows.StructScan(&uft); err != nil {
		return models.UserFieldsTaken{}, mapError(err)
	}

	return uft.toModel(), nil
}

// SetUserLastLogin sets the last login time for a user to now
func (s *AuthStore) SetUserLastLogin(ctx context.Context, id int64) error {
	const q = `
		UPDATE users
		SET last_login = NOW()
		WHERE id = :id
	`

	_, err := s.db.NamedExecContext(ctx, q, map[string]any{
		"id": id,
	})
	return mapError(err)
}

// CreateRefreshToken creates a new refresh token
func (s *AuthStore) CreateRefreshToken(ctx context.Context, refreshToken models.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (
			token_hash,
			user_id,
			user_agent,
			expires_at
		) VALUES (:token_hash, :user_id, :user_agent, :expires_at)
	`

	_, err := s.db.NamedExecContext(ctx, q, map[string]any{
		"token_hash": refreshToken.TokenHash,
		"user_id":    refreshToken.UserID,
		"user_agent": refreshToken.UserAgent,
		"expires_at": refreshToken.ExpiresAt,
	})
	return mapError(err)
}

// GetRefreshToken fetches a refresh token by its hash
func (s *AuthStore) GetRefreshToken(ctx context.Context, tokenHash string) (models.RefreshToken, error) {
	const q = `
		SELECT
			token_hash,
			user_id,
			user_agent,
			expires_at,
			created_at
		FROM refresh_tokens
		WHERE token_hash = :token_hash
	`

	rows, err := s.db.NamedQueryContext(ctx, q, map[string]any{
		"token_hash": tokenHash,
	})
	if err != nil {
		return models.RefreshToken{}, mapError(err)
	}
	defer closeRows(rows)
	if !rows.Next() {
		return models.RefreshToken{}, store.ErrNotFound
	}

	var row refreshTokenRow
	if err := rows.StructScan(&row); err != nil {
		return models.RefreshToken{}, mapError(err)
	}
	return row.toModel(), nil
}

// DeleteRefreshToken deletes a refresh token by its hash
func (s *AuthStore) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	const q = `
		DELETE FROM
			refresh_tokens
		WHERE token_hash = :token_hash
	`

	rows, err := s.db.NamedQueryContext(ctx, q, map[string]any{
		"token_hash": tokenHash,
	})
	if err != nil {
		return mapError(err)
	}
	defer closeRows(rows)
	return nil
}
