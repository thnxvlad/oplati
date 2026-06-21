package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	hserver "github.com/thnxvlad/oplati/internal/server"
	"github.com/thnxvlad/oplati/internal/server/hmiddlewares"
	"github.com/thnxvlad/oplati/internal/service/auth"
	"github.com/thnxvlad/oplati/internal/service/oplati"
	authStorage "github.com/thnxvlad/oplati/internal/storages/postgres/auth"
	postgresOplatiStorage "github.com/thnxvlad/oplati/internal/storages/postgres/oplati"
	redislimiter "github.com/thnxvlad/oplati/internal/storages/redis"
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
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://oplati:oplati@localhost:5432/oplati"
		//log.Fatal().Msg("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	defer pool.Close()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal().Msg("REDIS_ADDR is required")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}
	defer rdb.Close()

	redisStorage := redislimiter.New(rdb)

	authOplatiService := oplati.New(postgresOplatiStorage.New(pool))
	authService := auth.New(authStorage.New(pool), authOplatiService)
	publicServer := hserver.NewPublicServer(
		authOplatiService,
		authService,
		publicAddr,
		hmiddlewares.LoggingMiddleware,
		hmiddlewares.NewAuthMiddleware(authService),
		hmiddlewares.RateLimiterRedisMiddleware(redisStorage),
	)
	privateServer := hserver.NewPrivateServer(
		authOplatiService,
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
