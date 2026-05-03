package hserver

import (
	"net/http"

	hmiddlewares "github.com/thnxvlad/oplati/internal/server/hmiddlewares"
	"github.com/thnxvlad/oplati/internal/service/oplati"
)

type Server struct {
	oplatiService *oplati.Service
	*http.Server
}

func NewPublicServer(
	oplatiService *oplati.Service,
	addr string,
	mws ...func(next http.Handler) http.Handler,
) *Server {
	//  TODO: пополнить баланс
	//  TODO: просмотреть информацию о конкретном пользователе по id
	//  TODO: снять деньги с баланса
	//  TODO: перевести сумму денег с одного пользователя на другой

	// TODO: доделать public server

	return &Server{}
}

func NewPrivateServer(
	oplatiService *oplati.Service,
	addr string,
	mws ...func(next http.Handler) http.Handler,
) *Server {
	mux := http.NewServeMux()

	httpServer := http.Server{
		Addr:    addr,
		Handler: hmiddlewares.UseMiddlewares(mux, mws),
	}

	server := &Server{
		oplatiService: oplatiService,
		Server:        &httpServer,
	}

	mux.HandleFunc("POST /newUser", server.newUserHandler)
	//  TODO: просмотреть информацию о всех пользователях

	return server
}
