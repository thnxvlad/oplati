package hserver

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/thnxvlad/oplati/internal/domain"
	hmiddlewares "github.com/thnxvlad/oplati/internal/server/hmiddlewares"
)

type PublicOplatiService interface {
	Deposit(ctx context.Context, userId uuid.UUID, amount int) error
	Withdraw(ctx context.Context, userId uuid.UUID, amount int) error
	GetUser(ctx context.Context, id uuid.UUID) (domain.UserInfo, error)
	Transfer(ctx context.Context, senderID uuid.UUID, recipientID uuid.UUID, amount int) error
}

type PrivateOplatiService interface {
	GetUsersInfo(ctx context.Context) ([]domain.UserInfo, error)
}

type AuthService interface {
	SignIn(login, password string) (string, error)
	SignUp(login, password string) (string, error)
}

type PrivateServer struct {
	oplatiService PrivateOplatiService
	authService   AuthService
	*http.Server
}

type PublicServer struct {
	oplatiService PublicOplatiService
	*http.Server
}

func NewPublicServer(
	oplatiService PublicOplatiService,
	addr string,
	mws ...func(next http.Handler) http.Handler,
) *PublicServer {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /deposit", nil)
	mux.HandleFunc("PUT /withdraw", nil)
	mux.HandleFunc("GET /getUser", nil)
	mux.HandleFunc("PUT /transfer", nil)
	mux.HandleFunc("POST /newUser", nil)

	httpServer := http.Server{
		Addr:    addr,
		Handler: hmiddlewares.UseMiddlewares(mux, mws),
	}

	server := &PublicServer{
		oplatiService: oplatiService,
		Server:        &httpServer,
	}

	return server
}

func NewPrivateServer(
	oplatiService PrivateOplatiService,
	authService AuthService,
	addr string,
	mws ...func(next http.Handler) http.Handler,
) *PrivateServer {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /getUsersInfo", nil)
	mux.HandleFunc("POST /signin", nil)
	mux.HandleFunc("POST /signup", nil)

	httpServer := http.Server{
		Addr:    addr,
		Handler: hmiddlewares.UseMiddlewares(mux, mws),
	}

	server := &PrivateServer{
		oplatiService: oplatiService,
		Server:        &httpServer,
	}

	return server
}
