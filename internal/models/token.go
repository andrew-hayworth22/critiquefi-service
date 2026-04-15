package models

import "time"

type Claims struct {
	UserID  int64  `json:"uid"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type RefreshToken struct {
	TokenHash string
	UserID    int64
	UserAgent string
	ExpiresAt time.Time
	CreatedAt time.Time
}
