package store

import "context"

// Repo represents the functionality needed to store application data
type Repo interface {
	Users() UserStore

	BeginTx(ctx context.Context) (RepoTx, error)
}

// RepoTx represents the functionality needed for a database transaction across the application's stores
type RepoTx interface {
	Repo
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
