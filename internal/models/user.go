package models

import (
	"net/mail"
)

type User struct {
	ID           int64
	Email        string
	DisplayName  string
	Name         string
	IsAdmin      bool
	PasswordHash string
	IsActive     bool
}

// NewUser represents the data needed by the storage layer to create a new user
type NewUser struct {
	Email        string
	DisplayName  string
	Name         string
	PasswordHash string
}

// NewUserRequest represents the data needed from the transport layer to create a new user
type NewUserRequest struct {
	Email           string
	DisplayName     string
	Name            string
	Password        string
	ConfirmPassword string
}

// UserFieldsTaken represents fields that are taken when creating a user
type UserFieldsTaken struct {
	EmailTaken       bool
	DisplayNameTaken bool
}

func (u NewUserRequest) Validate() error {
	ve := ValidationErrors{}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		ve.Add("email", "invalid email address")
	}
	if len(u.DisplayName) < 3 {
		ve.Add("display_name", "display name must be at least 3 characters long")
	}
	if len(u.Name) < 3 {
		ve.Add("name", "name must be at least 3 characters long")
	}
	if len(u.Password) < 8 {
		ve.Add("password", "password must be at least 8 characters long")
	}
	if u.Password != u.ConfirmPassword {
		ve.Add("confirm_password", "passwords do not match")
	}

	if ve.Any() {
		return ve
	}
	return nil
}
