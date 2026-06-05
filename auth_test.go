package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/thnxvlad/oplati/internal/service/auth"
)

type mockDB struct {
	users map[string]struct {
		id   string
		hash string
	}
}

func (m *mockDB) GetUserByLogin(ctx context.Context, login string) (string, string, error) {
	u, ok := m.users[login]
	if !ok {
		return "", "", errors.New("not found")
	}
	return u.id, u.hash, nil
}

func (m *mockDB) SignUp(ctx context.Context, login, password, userID string) error {
	m.users[login] = struct {
		id   string
		hash string
	}{id: userID, hash: password}
	return nil
}

type mockOplati struct {
	calledWithID uuid.UUID
}

func (m *mockOplati) CreateUser(ctx context.Context, userID uuid.UUID) error {
	m.calledWithID = userID // Сохраняем ID, чтобы проверить его в тесте
	return nil
}

func TestSignUp(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_sercet")

	tests := []struct {
		name     string
		login    string
		password string
		wantErr  bool
	}{
		{
			name:     "mishanya_test",
			login:    "mik33",
			password: "mostpowerfulpswrd",
			wantErr:  false,
		},
		{
			name:     "andre_test",
			login:    "andr",
			password: "mostp649884648owerfulpswrd",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db := &mockDB{users: make(map[string]struct {
				id   string
				hash string
			})}
			oplati := &mockOplati{}
			service := auth.New(db, oplati)

			token, err := service.SignUp(context.Background(), tt.login, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wanted error = %v", err, tt.wantErr)
				return
			}

			if token == "" {
				t.Error("expected token, get empty string")
			}

			userInDB, ok := db.users[tt.login]
			if !ok {
				t.Error("user is not saved in db")
			}

			if userInDB.hash == tt.password {
				t.Error("password is not hashed")
			}

			if oplati.calledWithID == uuid.Nil {
				t.Error("OplatiService.CreateUser was not called")
			}
		})
	}
}

func TestGetUserByLogin(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_sercet")

	tests := []struct {
		name     string
		login    string
		password string
		wantErr  bool
	}{
		{
			name:     "mishanya_test",
			login:    "mik33",
			password: "mostpowerfulpswrd",
			wantErr:  false,
		},
		{
			name:     "andre_test",
			login:    "andr",
			password: "mostp649884648owerfulpswrd",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db := &mockDB{users: make(map[string]struct {
				id   string
				hash string
			})}
			oplati := &mockOplati{}
			service := auth.New(db, oplati)

			token, _ := service.SignUp(context.Background(), tt.login, tt.password)

			testID, err := service.GetAccountIdFromToken(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wanted error = %v", err, tt.wantErr)
				return
			}

			if testID != db.users[tt.login].id {
				t.Error("test id not equal real id")
			}
		})
	}
}
