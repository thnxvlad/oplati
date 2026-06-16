package auth

import (
	"context"
	"errors"
	"sync"
)

type Storage struct {
	authMu      sync.RWMutex
	loginData   map[string]string
	accountData map[string]string
}

func New() *Storage {
	return &Storage{
		loginData:   make(map[string]string),
		accountData: make(map[string]string),
	}
}

func (s *Storage) SignUp(ctx context.Context, login, password, userID string) error {
	s.authMu.Lock()
	defer s.authMu.Unlock()

	_, exists := s.loginData[login]

	if exists {
		return errors.New("login already exists")
	}

	s.loginData[login] = string(password)
	s.accountData[login] = userID

	return nil
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (userID, passwordHash string, err error) {
	s.authMu.Lock()
	defer s.authMu.Unlock()

	hash, ok := s.loginData[login]
	if !ok {
		return "", "", errors.New("user not found")
	}

	id := s.accountData[login]

	return id, hash, nil
}
