package integration_tests

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/thnxvlad/oplati/internal/service/auth"
	inmStorage "github.com/thnxvlad/oplati/internal/storages/inmemory/auth"
)

type mockOplati struct {
	calledWithID uuid.UUID
}

func (m *mockOplati) CreateUser(ctx context.Context, userID uuid.UUID) error {
	m.calledWithID = userID
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
			db := inmStorage.New()
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

			_, hashedpassword, err := db.GetUserByLogin(context.Background(), tt.login)
			if err != nil {
				t.Error("user is not saved in db")
			}

			if hashedpassword == tt.password {
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

			db := inmStorage.New()
			oplati := &mockOplati{}
			service := auth.New(db, oplati)

			token, _ := service.SignUp(context.Background(), tt.login, tt.password)

			testID, err := service.GetAccountIdFromToken(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wanted error = %v", err, tt.wantErr)
				return
			}

			realID, _, err := db.GetUserByLogin(context.Background(), tt.login)
			if err != nil {
				t.Error("get user error")
			}

			if realID != testID {
				t.Error("real id doesn't equal to test id")
			}
		})
	}
}
