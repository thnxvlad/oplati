package oplati

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thnxvlad/oplati/internal/domain"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{db: pool}
}

const createUserQuery = `
INSERT INTO users (id, balance) 
VALUES ($1, $2)
ON CONFLICT (id) DO NOTHING
`

func (s *Storage) CreateUser(ctx context.Context, userId uuid.UUID) error {
	cmdTag, err := s.db.Exec(ctx, createUserQuery, userId, 0)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user %s already exists", userId.String())
	}

	return nil
}

func (s *Storage) GetUser(ctx context.Context, userId uuid.UUID) (domain.UserInfo, error) {
	return domain.UserInfo{}, errors.New("not implemented")
}

func (s *Storage) GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error) {
	return nil, errors.New("not implemented")
}

func (s *Storage) Transfer(ctx context.Context, userIDFirst uuid.UUID, userIDSecond uuid.UUID, amount int) error {
	return errors.New("not implemented")
}

func (s *Storage) Deposit(ctx context.Context, userId uuid.UUID, amount int) error {
	return errors.New("not implemented")
}

func (s *Storage) Withdraw(ctx context.Context, userId uuid.UUID, amount int) error {
	return errors.New("not implemented")
}
