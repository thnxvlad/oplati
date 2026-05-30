package oplati

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thnxvlad/oplati/internal/domain"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{db: pool}
}

func (s *Storage) CreateUser(ctx context.Context, userId uuid.UUID) error {
	cmdTag, err := s.db.Exec(ctx, `
		INSERT INTO users (id, balance) 
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING
	`, userId, 0)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user %s already exists", userId.String())
	}

	return nil
}

func (s *Storage) GetUser(ctx context.Context, userId uuid.UUID) (domain.UserInfo, error) {
	row := s.db.QueryRow(ctx, `
		SELECT id, balance
		FROM users
		WHERE id = $1
	`, userId)
	ui := domain.UserInfo{}
	if err := row.Scan(&ui.Id, &ui.Balance); err != nil{
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserInfo{}, errors.New("user does not exist")
		}
		return domain.UserInfo{}, fmt.Errorf("Scan: %w", err)
	}

	return ui, nil
}

func (s *Storage) GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, balance
		FROM users
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	users := []domain.UserInfo{}
	for rows.Next() {
		user := domain.UserInfo{}
		if err := rows.Scan(&user.Id, &user.Balance); err != nil {
			return nil, fmt.Errorf("Scan: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) Transfer(ctx context.Context, userFrom, userTo uuid.UUID, amount int) error {
	// tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
	// 	IsoLevel: pgx.ReadCommitted,
	// })
	// if err != nil {
	// 	return fmt.Errorf("BeginTx: %w", err)
	// }
	// defer func(){
	// 	_ = tx.Rollback(ctx) 
	// }()

	// row := tx.QueryRow(ctx, `
	// 		SELECT  users
	// 		SET balance = balance - $1
	// 		WHERE id = $2
	// 	`, amount, userFrom)

	// if userFrom.String() < userTo.String() {
	// 	row := tx.QueryRow(ctx, `
	// 		UPDATE users
	// 		SET balance = balance - $1
	// 		WHERE id = $2
	// 		RETURNING balance
	// 	`, amount, userFrom)
	// 	var balance int
	// 	if err := row.Scan(&balance); err != nil {
	// 		if errors.Is(err, pgx.ErrNoRows){
	// 			return fmt.Errorf("userFrom does not exist")
	// 		}
	// 		return err
	// 	}
	// 	if balance < 0 {
	// 		return fmt.Errorf("not enough founds")
	// 	}
	// 	tmdTag, err := tx.Exec(ctx, `
	// 		UPDATE users
	// 		SET balance = balance + $1
	// 		WHERE id = $2
	// 	`, amount, userTo)
	// 	if tmdTag.RowsAffected()
	// }


	
	return errors.New("not implemented")
}

func (s *Storage) Deposit(ctx context.Context, userId uuid.UUID, amount int) error {
	cmdTag, err := s.db.Exec(ctx, `
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
	`, amount, userId)
	if err != nil {
		return fmt.Errorf("Exec: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user does not exist")
	}
	return nil
}

func (s *Storage) Withdraw(ctx context.Context, userId uuid.UUID, amount int) error {
	cmdTag, err := s.db.Exec(ctx, `
		UPDATE users
		SET balance = balance - $1
		WHERE 
			id = $2 AND
			balance >= $1
	`, amount, userId)
	if err != nil {
		return fmt.Errorf("Exec: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("user does not exist or not enough founds")
	}
	return nil
}
