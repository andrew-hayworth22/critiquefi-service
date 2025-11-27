package store

import (
	"context"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
)

type AuthStore interface {
	CreateUser(ctx context.Context, u types.User) (id int64, err error)
	GetUserByEmail(ctx context.Context, email string) (user *types.User, err error)
}
