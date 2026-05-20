package auth

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	OplatiService OplatiService
	db            AuthDB
}

func New(db AuthDB, oplatiService OplatiService) *Service {
	return &Service{
		db:            db,
		OplatiService: oplatiService,
	}
}

type AuthDB interface {
	SignIn(login, password string) (string, error)
	GetAccountIdFromToken(token string) (string, error)
	SignUp(login, password, userID string) error
}

type OplatiService interface {
	CreateUser(ctx context.Context, userID uuid.UUID) error
}

func (s *Service) SignIn(login, password string) (string, error) {
	token, err := s.db.SignIn(login, password)
	if err != nil {
		return "",err
	}

	return token, nil
}

func (s *Service) SignUp(login, password, userID string) error {
	err := s.db.SignUp(login, password, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateUser(ctx context.Context, userID uuid.UUID) error {
	err := s.OplatiService.CreateUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAccountIdFromToken(token string) (string, error) {
	userID, err := s.db.GetAccountIdFromToken(token)
	if err != nil {
		return "",err
	}

	return userID, nil
}

