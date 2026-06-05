package main

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thnxvlad/oplati/internal/service/auth"
	pgStorage "github.com/thnxvlad/oplati/internal/storages/postgres/auth"
)

var databaseURL string
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()
	databaseURL = "postgres://oplati:oplati@localhost:5432/oplati?sslmode=disable"
	testPool, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
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
