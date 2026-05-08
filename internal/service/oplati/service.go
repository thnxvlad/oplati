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
	CreateUser(ctx context.Context, ui domain.UserInfo) error
	UpdateBalance(ctx context.Context, userId uuid.UUID, amount int) (domain.UserInfo, error)
	GetUsersInfo(ctx context.Context)([]domain.UserInfo)
	GetUser(ctx context.Context, id uuid.UUID)(domain.UserInfo, error)
	Transfer(ctx context.Context, senderID uuid.UUID, recipientID uuid.UUID, amount int) (error)
}

func (s *Service) Transfer(ctx context.Context, senderID uuid.UUID, recipientID uuid.UUID, amount int) (error) {
	err := s.db.Transfer(ctx, senderID, recipientID, amount)
	return err
}

func (s *Service) GetUsersInfo(ctx context.Context) ([]domain.UserInfo) {
	users := s.db.GetUsersInfo(ctx)
	return users
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (domain.UserInfo, error) {
	ui, err := s.db.GetUser(ctx, id)

	if err != nil {
		return domain.UserInfo{}, err
	}

	return ui, nil
} 

func (s *Service) CreateUser(ctx context.Context, name string) (domain.UserInfo, error) {
	ui := domain.UserInfo{
		Id:      uuid.New(),
		Name:    name,
		Balance: 0,
	}

	err := s.db.CreateUser(ctx, ui)
	if err != nil {
		return domain.UserInfo{}, err
	}

	return ui, nil
}

func (s *Service) Deposit(ctx context.Context, userId uuid.UUID, amount int) (domain.UserInfo, error) {
	ui, err := s.db.UpdateBalance(ctx, userId, amount)
	if err != nil {
		return domain.UserInfo{}, err
	}

	return ui, nil
}

func (s *Service) Withdraw(ctx context.Context, userId uuid.UUID, amount int) (domain.UserInfo, error) {
	ui, err := s.db.UpdateBalance(ctx, userId, -amount)
	if err != nil {
		return domain.UserInfo{}, err
	}

	return ui, nil
}