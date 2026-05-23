package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	hserver "github.com/thnxvlad/oplati/internal/server"
	"github.com/thnxvlad/oplati/internal/server/hmiddlewares"
	"github.com/thnxvlad/oplati/internal/service/auth"
	"github.com/thnxvlad/oplati/internal/service/oplati"
	authStorage "github.com/thnxvlad/oplati/internal/storages/inmemory/auth"
	oplatiStorage "github.com/thnxvlad/oplati/internal/storages/inmemory/oplati"
)

const (
	publicAddr  = ":8082"
	privateAddr = ":8081"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
		NoColor:    false,
	})
}

func main() {
	oplatiService := oplati.New(oplatiStorage.NewStorage())
	authService := auth.New(authStorage.New(), oplatiService)
	publicServer := hserver.NewPublicServer(oplatiService, authService, publicAddr, hmiddlewares.LoggingMiddleware)
	privateServer := hserver.NewPrivateServer(
		oplatiService,
		privateAddr,
		hmiddlewares.LoggingMiddleware,
		hmiddlewares.NewAuthMiddleware(authService),
	)

	go func() {
		log.Info().Str("addr", publicAddr).Msg("public server started...")
		err := publicServer.ListenAndServe()
		if err != nil {
			log.Err(err).Msg("failed to start public server")
		}
	}()

	go func() {
		log.Info().Str("addr", privateAddr).Msg("private server started...")
		err := privateServer.ListenAndServe()
		if err != nil {
			log.Err(err).Msg("failed to start private server")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	stop()
}
