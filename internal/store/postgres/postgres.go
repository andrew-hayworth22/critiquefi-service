// Package postgres wraps a connection to a postgres database
package postgres

import (
	"context"
	"database/sql"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db *sql.DB
}

func Open(dsn string) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return &DB{db}, nil
}

// repo is the postgres implementation of the application's store
type repo struct {
	db    *sql.DB
	tx    *sql.Tx
	users userPG
}

func (d *DB) Repo() *repo {
	r := &repo{db: d.db}
	r.users.repo = r
	return r
}

func (r *repo) Users() store.UserStore {
	return &r.users
}

func (r *repo) BeginTx(ctx context.Context) (store.RepoTx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	repoTx := &repo{db: r.db, tx: tx}
	repoTx.users.repo = repoTx
	return repoTx, nil
}

func (r *repo) Commit(ctx context.Context) error {
	return r.tx.Commit()
}

func (r *repo) Rollback(ctx context.Context) error {
	return r.tx.Rollback()
}

func (r *repo) execer() interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *repo) queryer() interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}
