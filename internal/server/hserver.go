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
	SignIn(ctx context.Context, login, password string) (string, error)
	SignUp(ctx context.Context, login, password string) (string, error)
}

type PrivateServer struct {
	oplatiService PrivateOplatiService
	*http.Server
}

type PublicServer struct {
	oplatiService PublicOplatiService
	authService   AuthService
	*http.Server
}

func NewPublicServer(
	oplatiService PublicOplatiService,
	authService AuthService,
	addr string,
	mws ...func(next http.Handler) http.Handler,
) *PublicServer {
	mux := http.NewServeMux()

	httpServer := http.Server{
		Addr:    addr,
		Handler: hmiddlewares.UseMiddlewares(mux, mws),
	}

	server := &PublicServer{
		oplatiService: oplatiService,
		authService:   authService,
		Server:        &httpServer,
	}

	mux.HandleFunc("POST /deposit", server.depositHandler)
	mux.HandleFunc("POST /withdraw", server.withdrawHandler)
	mux.HandleFunc("GET /getUser", server.getUserHandler)
	mux.HandleFunc("POST /transfer", server.transferHandler)
	mux.HandleFunc("POST /signin", server.signInHandler)
	mux.HandleFunc("POST /signup", server.signUpHandler)

	return server
}

func NewPrivateServer(
	oplatiService PrivateOplatiService,
	addr string,
	mws ...func(next http.Handler) http.Handler,
) *PrivateServer {
	mux := http.NewServeMux()

	httpServer := http.Server{
		Addr:    addr,
		Handler: hmiddlewares.UseMiddlewares(mux, mws),
	}

	server := &PrivateServer{
		oplatiService: oplatiService,
		Server:        &httpServer,
	}

	mux.HandleFunc("GET /getUsersInfo", server.getUsersInfoHandler)

	return server
}
