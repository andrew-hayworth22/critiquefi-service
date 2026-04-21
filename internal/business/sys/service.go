// Package sys provides business logic for system-related operations.
package sys

import "context"

// Store defines logic for system-related storage operations.
type Store interface {
	Ping(ctx context.Context) error
}

// Bus defines the system-related business logic.
type Bus struct {
	store Store
}

// New creates a new system-related business logic package.
func New(store Store) *Bus {
	return &Bus{store: store}
}

// Ping tests the database connection.
func (b *Bus) Ping(ctx context.Context) error {
	return b.store.Ping(ctx)
}
