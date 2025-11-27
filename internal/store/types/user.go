package types

import "time"

// User represents a user in the system
type User struct {
	ID           int64
	Email        string
	DisplayName  string
	Name         string
	PasswordHash string
	IsAdmin      bool
	LastLogin    NullableTime
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    NullableTime
}

// UserUpdate represents a user update payload
type UserUpdate struct {
	Email        *string
	DisplayName  *string
	Name         *string
	PasswordHash *string
	IsAdmin      *bool
	IsActive     *bool
}
