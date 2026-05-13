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
	mux := http.NewServeMux()

	httpServer := http.Server{
		Addr:    addr,
		Handler: hmiddlewares.UseMiddlewares(mux, mws),
	}

	server := &Server{
		oplatiService: oplatiService,
		Server:        &httpServer,
	}

	mux.HandleFunc("PUT /deposit", server.depositHandler)
	mux.HandleFunc("PUT /withdraw", server.withdrawHandler)
	mux.HandleFunc("GET /getUser", server.getUserHandler)
	mux.HandleFunc("PUT /transfer", server.transferHandler)
	mux.HandleFunc("POST /newUser", server.newUserHandler)

	return server
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

	mux.HandleFunc("GET /getUsersInfo", server.getUsersInfoHandler)
	
	return server
}
