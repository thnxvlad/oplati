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
	"github.com/thnxvlad/oplati/internal/service/oplati"
	"github.com/thnxvlad/oplati/internal/storages/inmemory"
)

const (
	publicAddr  = ":8080"
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
	oplatiService := oplati.New(inmemory.NewStorage())
	// TODO: доделать public server
	publicServer := hserver.NewPublicServer(oplatiService, publicAddr, hmiddlewares.LoggingMiddleware)
	privateServer := hserver.NewPrivateServer(oplatiService, privateAddr, hmiddlewares.LoggingMiddleware)

	/*	go func() {
			log.Info().Str("addr", publicAddr).Msg("public server started...")
			err := publicServer.ListenAndServe()
			if err != nil {
				log.Err(err).Msg("failed to start public server")
			}
		}()
	*/

	go func() {
		log.Info().Str("addr", privateAddr).Msg("private server started...")
		err := privateServer.ListenAndServe()
		if err != nil {
			log.Err(err).Msg("failed to start private server")
		}
	}()

	go func() {
		log.Info().Str("addr", publicAddr).Msg("public server started...")
		err := publicServer.ListenAndServe()
		if err != nil {
			log.Err(err).Msg("failed to start public server")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	stop()
}
