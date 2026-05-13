package oplati

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/thnxvlad/oplati/internal/domain"
)



type userClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
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
	GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.UserInfo, error)
	Transfer(ctx context.Context, senderID uuid.UUID, recipientID uuid.UUID, amount int) error
}

func (s *Service) Transfer(ctx context.Context, senderID uuid.UUID, recipientID uuid.UUID, amount int) error {
	err := s.db.Transfer(ctx, senderID, recipientID, amount)
	return err
}

func (s *Service) GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error) {
	users, err := s.db.GetUsersInfo(ctx)
	return users, err
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (domain.UserInfo, error) {
	ui, err := s.db.GetUser(ctx, id)

	if err != nil {
		return domain.UserInfo{}, err
	}

	return ui, nil
}

func generateToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := &userClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "oplati-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(domain.JwtSecret))
}

func (s *Service) CreateUser(ctx context.Context, login, password string) (string, error) {
	ui := domain.UserInfo{
		Id:       uuid.New(),
		Name:     login,
		Balance:  0,
		Login:    login,
		Password: password,
	}
	
	err := s.db.CreateUser(ctx, ui)
	if err != nil {
		return "", err
	}
	token, err := generateToken(ui.Id);

	return token, err
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
