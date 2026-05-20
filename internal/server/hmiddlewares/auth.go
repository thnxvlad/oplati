package hmiddlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/thnxvlad/oplati/internal/service/auth"
)

type AccountIdContextKey struct{}


func NewAuthMiddleware(
	authService auth.Service, 
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")
			accountId, err := authService.GetAccountIdFromToken(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AccountIdContextKey{}, accountId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAccountIdFromContext(ctx context.Context) (uuid.UUID, error) {
	accountId, ok := ctx.Value(AccountIdContextKey{}).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("account id not found in context")
	}

	return accountId, nil
}
