package hmiddlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AccountIdContextKey struct{}

type userClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func getAccountIdFromToken(token string) (string, error) {
	claims, err := jwt.ParseWithClaims(token, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.JwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	return claims.Claims.(*userClaims).UserID, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		accountId, err := getAccountIdFromToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), AccountIdContextKey{}, accountId)))
	})
}
