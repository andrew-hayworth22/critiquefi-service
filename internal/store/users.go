package store

import (
	"context"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
	"time"
)

// User represents a user in the system
type User struct {
	ID           int64
	Email        string
	DisplayName  string
	Name         string
	PasswordHash string
	IsAdmin      bool
	LastLogin    types.NullableTime
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    types.NullableTime
}

// UserUpdate represents the fields that can be updated in a User. Nil fields are not updated
type UserUpdate struct {
	Email        *string
	DisplayName  *string
	Name         *string
	PasswordHash *string
	IsAdmin      *bool
	IsActive     *bool
}

// UserStore defines the functionality needed to store users
type UserStore interface {
	Create(ctx context.Context, u User) (int64, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id int64, u UserUpdate) error
}
