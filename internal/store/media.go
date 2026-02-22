package store

import (
	"context"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
)

type MediaStore interface {
	GetMedia(ctx context.Context) ([]*types.Media, error)

	GetFilmById(ctx context.Context, id int64) (*types.Film, error)
	CreateFilm(ctx context.Context, film *types.Film) (id int64, err error)

	/*GetBookById(ctx context.Context, id int64) (*types.Book, error)
	CreateBook(ctx context.Context, book *types.Book) (id int64, err error)

	GetMusicById(ctx context.Context, id int64) (*types.Music, error)
	CreateMusic(ctx context.Context, music *types.Music) (id int64, err error)

	GetShowById(ctx context.Context, id int64) (*types.Show, error)
	CreateShow(ctx context.Context, show *types.Show) (id int64, err error)
	*/
}
