package store

import "context"

type SysStore interface {
	Ping(ctx context.Context) error
}
