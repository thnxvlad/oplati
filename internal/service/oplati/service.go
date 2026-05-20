package oplati

import (
	"context"

	"github.com/google/uuid"
	"github.com/thnxvlad/oplati/internal/domain"
)

type Service struct {
	db OplatiDatabase
}

func New(db OplatiDatabase) *Service {
	return &Service{
		db: db,
	}
}

type OplatiDatabase interface {
	CreateUser(ctx context.Context, userId uuid.UUID) error
	GetUser(ctx context.Context, userId uuid.UUID) (domain.UserInfo, error)
	GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error)
	Transfer(ctx context.Context, userIDFirst uuid.UUID, userIDSecond uuid.UUID, amount int) error
	UpdateBalance(ctx context.Context, userId uuid.UUID, amount int) error
}

func (s *Service) GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error) {
	users, err := s.db.GetUsersInfo(ctx)
	return users, err
}

func (s *Service) Deposit(ctx context.Context, userId uuid.UUID, amount int) error {
	err := s.db.UpdateBalance(ctx, userId, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Withdraw(ctx context.Context, userId uuid.UUID, amount int) error {
	err := s.db.UpdateBalance(ctx, userId, -amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Transfer(ctx context.Context, senderID uuid.UUID, recipientID uuid.UUID, amount int) error {
	err := s.db.Transfer(ctx, senderID, recipientID, amount)
	return err
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (domain.UserInfo, error) {
	// ui, err := s.db.GetUser(ctx, id)

	// if err != nil {
	// 	return domain.UserInfo{}, err
	// }

	// return ui, nil
	return domain.UserInfo{},nil
}
