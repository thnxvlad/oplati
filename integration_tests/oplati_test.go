package integration_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/thnxvlad/oplati/internal/service/oplati"
	postgresOplatiStorage "github.com/thnxvlad/oplati/internal/storages/postgres/oplati"
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

		if err != nil {
			log.Fatalf("failed to run: %s", err)
		}

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
	`
	_, err := pool.Exec(ctx, schema)
	if err != nil {
		log.Fatalf("failed to apply schema: %v", err)
	}
}
func TestCreateUser_Integration(t *testing.T) {
	storage := postgresOplatiStorage.New(testPool)
	service := oplati.New(storage)
	id := uuid.New()

	err := service.CreateUser(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	user, err := service.GetUser(context.Background(), id)
	if err != nil {
		t.Fatal("user does not exist")
	}

	if user.Balance != 0 {
		t.Fatalf("balance %d", user.Balance)
	}
}

func TestDeposit_Integration(t *testing.T) {
	storage := postgresOplatiStorage.New(testPool)
	service := oplati.New(storage)
	id := uuid.New()

	err := service.CreateUser(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	err = service.Deposit(context.Background(), id, 30)

	if err != nil {
		t.Fatal(err)
	}

	user, err := service.GetUser(context.Background(), id)
	if err != nil {
		t.Fatal("user does not exist")
	}
	
	if user.Balance != 30 {
		t.Fatalf("user balance must be 30, got %d", user.Balance)
	}
}
