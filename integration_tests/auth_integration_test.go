package auth

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/thnxvlad/oplati/internal/service/auth"
	pgStorage "github.com/thnxvlad/oplati/internal/storages/postgres/auth"
)

var databaseURL string
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()
	databaseURL = os.Getenv("DATABASE_URL")
	var pgContainer *postgres.PostgresContainer

	if databaseURL == "" {
		pgContainer, err = postgres.Run(ctx,
			"postgres:16-alpine",
			postgres.WithDatabase("oplati"),
			postgres.WithUsername("oplati"),
			postgres.WithPassword("oplati"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("db is ready!").
					WithOccurrence(2).
					WithStartupTimeout(30*time.Second)),
		)

		databaseURL, err = pgContainer.ConnectionString(ctx, "sslmode-disable")
		if err != nil {
			log.Fatalf("failed to get connection string: %s", err)
		}
	}
	testPool, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		panic(err)
	}
	ensureSchema(ctx, testPool)

	code := m.Run()

	testPool.Close()
	if pgContainer != nil {
		pgContainer.Terminate(ctx)
	}
	os.Exit(code)
}

func ensureSchema(ctx context.Context, pool *pgxpool.Pool) {
	schema := `
	CREATE TABLE users (
    id      UUID    PRIMARY KEY,
    balance INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT users_balance_non_negative CHECK (balance >= 0));

	CREATE TABLE accounts (
    login         TEXT PRIMARY KEY,
    password_hash TEXT NOT NULL,
    user_id       UUID NOT NULL);
	`

	_, err := pool.Exec(ctx, schema)
	if err != nil {
		log.Fatalf("failed to apply schema: %v", err)
	}
}

func setupTest(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "TRUNCATE TABLE accounts, users CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func TestSignUp_Integration(t *testing.T) {
	setupTest(t)

	storage := pgStorage.New(testPool)
	service := auth.New(storage, &mockOplati{})

	ctx := context.Background()
	login := "mishanya_test"
	password := "mishanya's pswrd"

	token, err := service.SignUp(ctx, login, password)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("token is empty")
	}

	_, err = service.SignUp(ctx, login, password)
	if err == nil {
		t.Error("expected error cause it's a duplicate")
	}
}

func TestGetUserByLogin_Integration(t *testing.T) {
	setupTest(t)

	storage := pgStorage.New(testPool)
	service := auth.New(storage, &mockOplati{})

	ctx := context.Background()
	login := "mishanya_test"
	password := "mishanya's pswrd"
	token, _ := service.SignUp(ctx, login, password)

	testID, err := service.GetAccountIdFromToken(token)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	realID, _, err := storage.GetUserByLogin(ctx, login)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if realID != testID {
		t.Error("test id is not equal real")
	}
}
