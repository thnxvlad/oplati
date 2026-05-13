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
	GetUser(ctx context.Context, userId uuid.UUID) error
	GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error)
	Transfer(ctx context.Context, userIDFirst uuid.UUID, userIDSecond uuid.UUID, amount int) error
	Deposit(ctx context.Context, userId uuid.UUID, amount int) error
	Withdraw(ctx context.Context, userId uuid.UUID, amount int) error
}
