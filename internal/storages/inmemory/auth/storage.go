package auth

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ToDo: реализовать всё
var (
jwtSecret = os.Getenv("JWT_SECRET")
loginData = make(map[string]string)
accountData = make(map[string]string)
authMu sync.RWMutex
) 

type userClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func generateToken(userID string) (string, error) {
	now := time.Now()
	claims := &userClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24*time.Hour)),
			IssuedAt: jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer: "oplati-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func validateLogin(login, password string) bool {
	hashedPassword, ok := loginData[login]
	if !ok {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func Login(login, password string) (string, error) {
	authMu.Lock()
	defer authMu.Unlock()

	if !validateLogin(login, password) {
		return "",errors.New("invalid login or password")
	}

	accountId, ok := accountData[login]
	if !ok {
		return "",errors.New("account not found")
	}

	return generateToken(accountId)
}

func GetAccountIdFromToken(token string) (string, error) {
	claims, err := jwt.ParseWithClaims(token, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	parsedClaims, ok := claims.Claims.(*userClaims)
	if !ok {
		return "", errors.New("invalid claims type")
	}

	return parsedClaims.UserID, nil
}