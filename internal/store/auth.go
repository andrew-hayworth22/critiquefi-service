package store

import (
	"context"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
)

type AuthStore interface {
	CreateUser(ctx context.Context, u types.User) (id int64, err error)
	GetUserByID(ctx context.Context, id int64) (*types.User, error)
	GetUserByEmail(ctx context.Context, email string) (user *types.User, err error)

	CreateRefreshToken(ctx context.Context, refreshToken *types.RefreshToken) (err error)
	GetRefreshToken(ctx context.Context, token string) (*types.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}
