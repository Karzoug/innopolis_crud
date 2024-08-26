package app

import (
	"context"
	"crud/internal/handler"
	"crud/internal/pkg/authclient"
	"crud/internal/pkg/server"
	"crud/internal/repository/cache"
	"crud/internal/service"
	"errors"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	// initialize dbs
	db, err := cache.RecipeCacheInit(ctx, &wg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize user database")
	}

	authclient.Init("localhost:8000")

	// initialize service
	service.Init(db)

	go func() {
		log.Info().Str("port", "8080").Msg("starting CRUD server")
		err := server.Run("localhost:8080", handler.ServerHandler)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server run")
		}
	}()

	<-ctx.Done()

	if err = server.Stop(); err != nil {
		log.Error().Err(err).Msg("server was not gracefully shutdown")
	}
	wg.Wait()

	log.Info().Msg("CRUD service stopped")
}
