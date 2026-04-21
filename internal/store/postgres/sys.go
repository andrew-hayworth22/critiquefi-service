package postgres

import (
	"context"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/jmoiron/sqlx"
)

// SysStore is the postgres implementation of the system store
type SysStore struct {
	db *sqlx.DB
}

// NewSysStore creates a new system store
func NewSysStore(db *sqlx.DB) *SysStore {
	return &SysStore{db: db}
}

// Ping tests the database connection
func (s *SysStore) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return store.ErrInternal
	}
	return nil
}
