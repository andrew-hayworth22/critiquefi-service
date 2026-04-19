package postgres_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

func TestSysStore_Ping(t *testing.T) {
	t.Run("ping", func(t *testing.T) {
		s := postgres.NewSysStore(testDB)

		err := s.Ping(context.Background())
		testutil.CheckErr(err, nil, t)
	})
}
