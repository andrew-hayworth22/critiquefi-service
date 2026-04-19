package postgres_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
)

func TestSysStore_Ping(t *testing.T) {
	t.Run("ping", func(t *testing.T) {
		s := postgres.NewSysStore(testDB)

		err := s.Ping(context.Background())
		checkErr(err, nil, t)
	})
}
