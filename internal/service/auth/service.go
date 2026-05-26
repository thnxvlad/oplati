package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	OplatiService OplatiService
	db            AuthDB
	jwtSecret     string
}

type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthDB interface {
	GetUserByLogin(ctx context.Context, login string) (userID, passwordHash string, err error)
	SignUp(ctx context.Context, login, password, userID string) error
}

type OplatiService interface {
	CreateUser(ctx context.Context, userID uuid.UUID) error
}

func New(db AuthDB, oplatiService OplatiService) *Service {
	return &Service{
		db:            db,
		OplatiService: oplatiService,
		jwtSecret:     os.Getenv("JWT_SECRET"),
	}
}

func (s *Service) SignIn(ctx context.Context, login, password string) (string, error) {
	id, passwordHash, err := s.db.GetUserByLogin(ctx, login)
	if err != nil {
		return "", errors.New("invalid password or login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password or login")
	}

	return s.generateToken(id)
}

func (s *Service) SignUp(ctx context.Context, login, password string) (string, error) {
	id := uuid.New()

	err := s.db.SignUp(ctx, login, password, id.String())
	if err != nil {
		return "", err
	}

	err = s.OplatiService.CreateUser(ctx, id)
	if err != nil {
		return "", err
	}

	token, err := s.SignIn(ctx, login, password)
	if err != nil {
		return "", err
	}

	return token, err
}

func (s *Service) GetAccountIdFromToken(token string) (string, error) {
	claims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	parsedClaims, ok := claims.Claims.(*UserClaims)
	if !ok {
		return "", errors.New("invalid claims type")
	}

	return parsedClaims.UserID, nil
}

func (s *Service) generateToken(userID string) (string, error) {
	now := time.Now()
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "oplati-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
