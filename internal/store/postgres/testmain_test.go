package postgres_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sqlx.DB

var seededUsers []models.User

// TestMain runs before all tests and sets up a test database
func TestMain(m *testing.M) {

	// Start the postgres test container
	ctx := context.Background()
	container, err := postgresContainer.Run(ctx, "postgres:18-alpine",
		postgresContainer.WithDatabase("critiquefi_test"),
		postgresContainer.WithUsername("postgres"),
		postgresContainer.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgresContainer container: %v", err)
	}
	defer func(container *postgresContainer.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			log.Fatalf("failed to terminate postgresContainer container: %v", err)
		}
	}(container, ctx)

	// Connect to the test database and run migrations
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}
	mig, err := migrate.New("file://../../../migrations", connStr)
	if err != nil {
		log.Fatalf("error connecting to database for migrations: %v", err)
	}
	if err := mig.Up(); err != nil {
		log.Fatal(err)
	}

	// Create the test database connection
	testDB, err = postgres.NewDB(ctx, postgres.DBConfig{
		URL:             connStr,
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	})
	if err != nil {
		log.Fatalf("failed to connect to test db: %v", err)
	}

	// Run tests
	os.Exit(m.Run())
}

// resetTables cleans up and seeds the tables used in the tests
func resetTables(t *testing.T) {
	t.Helper()
	_, err := testDB.ExecContext(context.Background(), `
        TRUNCATE TABLE refresh_tokens, users RESTART IDENTITY CASCADE;

        INSERT INTO users (id, email, display_name, name, password_hash, is_admin, is_active) VALUES
        (1, 'user@critiquefi.com', 'test.user', 'Test User', 'password', false, true),
        (2, 'admin@critiquefi.com', 'admin.user', 'Test Admin', 'password', true, true),
        (3, 'deactivated@critiquefi.com', 'deactivated.user', 'Deactivated User', 'password', false, false);

		INSERT INTO refresh_tokens (token_hash, user_id, user_agent, expires_at) VALUES
		('test-token-hash', 1, 'test-user-agent', '2027-01-01 00:00:00');
		
        SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT MAX(id) FROM users));
    `)
	if err != nil {
		t.Fatalf("failed to seed tables: %v", err)
	}
}
