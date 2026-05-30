package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	hserver "github.com/thnxvlad/oplati/internal/server"
	"github.com/thnxvlad/oplati/internal/server/hmiddlewares"
	"github.com/thnxvlad/oplati/internal/service/auth"
	"github.com/thnxvlad/oplati/internal/service/oplati"
	authStorage "github.com/thnxvlad/oplati/internal/storages/inmemory/auth"
	postgresOplatiStorage "github.com/thnxvlad/oplati/internal/storages/postgres/oplati"
)

const (
	publicAddr         = ":8082"
	privateAddr        = ":8081"
	defaultDatabaseURL = "postgres://oplati:oplati@localhost:5432/oplati?sslmode=disable"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
		NoColor:    false,
	})
}

func main() {
	pool, err := pgxpool.New(context.Background(), defaultDatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	defer pool.Close()

	oplatiService := oplati.New(postgresOplatiStorage.New(pool))
	authService := auth.New(authStorage.New(), oplatiService)
	
	publicServer := hserver.NewPublicServer(
		oplatiService,
		authService,
		publicAddr,
		hmiddlewares.LoggingMiddleware,
	)
	privateServer := hserver.NewPrivateServer(
		oplatiService,
		privateAddr,
		hmiddlewares.LoggingMiddleware,
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
