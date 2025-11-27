package store

import (
	"context"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
)

// UserUpdate represents the fields that can be updated in a User. Nil fields are not updated

// UserStore defines the functionality needed to store auth
type UserStore interface {
	Create(ctx context.Context, u types.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	Update(ctx context.Context, id int64, u types.UserUpdate) error
}
