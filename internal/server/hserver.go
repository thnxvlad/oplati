package hserver

import (
	"net/http"

	hmiddlewares "github.com/thnxvlad/oplati/internal/server/hmiddlewares"
)

func NewPublicServer(addr string, mws ...func(next http.Handler) http.Handler) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /newUser", NewUserFunc)
	// пополнить баланс
	// просмотреть информацию о всех пользователях

	server := http.Server{
		Addr:    addr,
		Handler: hmiddlewares.UseMiddlewares(mux, mws),
	}

	return &server
}
