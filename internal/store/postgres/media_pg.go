package postgres

import (
	"context"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type mediaPG struct {
	db *pgxpool.Pool
}

func (m *mediaPG) GetMedia(ctx context.Context) ([]*types.Media, error) {
	const q = `
SELECT
	id,
	media_type,
	title,
	release_date,
	
`
}
