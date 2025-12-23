package types

import "time"

type RefreshToken struct {
	Token     string
	UserID    int64
	UserAgent string
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time
}
