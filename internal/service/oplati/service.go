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
	Deposit(ctx context.Context, userId uuid.UUID, amount int) (domain.UserInfo, error)
	GetUsersInfo(ctx context.Context)([]domain.UserInfo, error)
}

func (s *Service) GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error) {
	users, err := s.db.GetUsersInfo(ctx)
	if err != nil{
		return []domain.UserInfo{}, err
	}
	return users, err
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
	ui, err := s.db.Deposit(ctx, userId, amount)
	if err != nil {
		return domain.UserInfo{}, err
	}

	return ui, nil
}
