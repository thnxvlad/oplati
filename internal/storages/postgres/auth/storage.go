package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{db: pool}
}

const signUpQuery = `
INSERT INTO accounts (login, password_hash, user_id)
VALUES ($1, $2, $3)
ON CONFLICT (login) DO NOTHING
`

const getUserByLoginQuery = `
SELECT user_id, password_hash from accounts
WHERE login = $1
`

func (s *Storage) SignUp(ctx context.Context, login, password, userID string) error {
	cmdTag, err := s.db.Exec(ctx, signUpQuery, login, password, userID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user %s already exists", userID)
	}
	return nil
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (returnedUserID, returnedHash string, err error) {
	row := s.db.QueryRow(ctx, getUserByLoginQuery, login)

	if err1 := row.Scan(&returnedUserID, &returnedHash); err1 != nil {
		if errors.Is(err1, pgx.ErrNoRows) {
			return "", "", errors.New("user does not exist")
		}
		return "", "", fmt.Errorf("Scan: %w", err1)
	}
	return returnedUserID, returnedHash, nil
}
