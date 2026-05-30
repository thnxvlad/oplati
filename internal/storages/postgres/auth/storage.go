package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	cmdTag, err := s.db.Exec(ctx, signUpQuery, login, hashedPassword, userID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user %s already exists", userID)
	}
	return nil
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (userID, passwordHash string, err error) {
	row := s.db.QueryRow(ctx, getUserByLoginQuery, login)
	var returnedUserID, returnedHash string

	if err1 := row.Scan(&returnedUserID, &returnedHash); err1 != nil {
		if errors.Is(err1, pgx.ErrNoRows) {
			return returnedUserID, returnedHash, errors.New("user does not exist")
		}
		return returnedUserID, returnedHash, fmt.Errorf("Scan: %w", err1)
	}
	return returnedUserID, returnedHash, nil
}
