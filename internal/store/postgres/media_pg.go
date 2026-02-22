package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type mediaPG struct {
	db *pgxpool.Pool
}

func NewMediaPG(db *pgxpool.Pool) *mediaPG {
	return &mediaPG{db: db}
}

func (m *mediaPG) GetMedia(ctx context.Context) ([]*types.Media, error) {
	const q = `
SELECT
	id,
	media_type,
	title,
	release_date,
	year,
	description,
	external_references,
	created_at,
	created_by,
	updated_at,
	updated_by
FROM media;
`

	rows, err := m.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	medias := make([]*types.Media, 0)
	for rows.Next() {
		var (
			md           types.Media
			externalRefs json.RawMessage
		)

		if err := rows.Scan(&md.ID, &md.MediaType, &md.Title, &md.ReleaseDate, &md.Year, &md.Description, &externalRefs, &md.CreatedAt, &md.CreatedBy, &md.UpdatedAt, &md.UpdatedBy); err != nil {
			return nil, err
		}

		if len(externalRefs) == 0 {
			md.ExternalReferences = "{}"
		} else {
			md.ExternalReferences = string(externalRefs)
		}

		medias = append(medias, &md)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return medias, nil
}

func (m *mediaPG) GetFilmById(ctx context.Context, id int64) (*types.Film, error) {
	const q = `
SELECT
	id,
	media_type,
	title,
	release_date,
	year,
	description,
	external_references,
	created_at,
	updated_at,
	film_type,
	runtime_minutes
FROM media
INNER JOIN films
	ON media.id = films.media_id
WHERE id = $1`

	row := m.db.QueryRow(ctx, q, id)
	var film types.Film
	if err := row.Scan(&film.ID, &film.MediaType, &film.Title, &film.ReleaseDate, &film.Year, &film.Description, &film.ExternalReferences, &film.CreatedAt, &film.CreatedBy, &film.UpdatedAt, &film.UpdatedBy, &film.FilmType, &film.RuntimeMinutes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &film, nil
}

func (m *mediaPG) CreateFilm(ctx context.Context, film *types.Film) (int64, error) {
	const q = `
INSERT INTO media (
	media_type,
	title,
	release_date,
	year,
	description,
	external_references,
	created_at,
    created_by,
	updated_at,
	updated_by
) VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7, NULL, NULL)
RETURNING id`

	var id int64
	err := m.db.QueryRow(ctx, q,
		film.MediaType,
		film.Title,
		film.ReleaseDate,
		film.Year,
		film.Description,
		film.ExternalReferences,
		film.CreatedBy,
	).Scan(&id)
	if err != nil {
		return id, err
	}

	const q2 = `
INSERT INTO film (
	media_id,
    film_type,
	runtime_minutes
) VALUES ($1, $2, $3)`

	_, err = m.db.Exec(ctx, q2, id, film.FilmType, film.RuntimeMinutes)
	return id, err
}
