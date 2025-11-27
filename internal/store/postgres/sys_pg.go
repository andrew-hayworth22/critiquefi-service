package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SysPG struct {
	db *pgxpool.Pool
}

func NewSysPG(db *pgxpool.Pool) *SysPG {
	return &SysPG{db: db}
}

func (s *SysPG) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}
