package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"dynamic-user-segmentation/pkg/config"
	"dynamic-user-segmentation/pkg/controller"
	"dynamic-user-segmentation/pkg/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	localEnv = "local"
	demoEnv  = "demo"
)

func main() {
	cfg := config.MustLoad()
	log := setUpLogger(cfg.Env)
	log.Info(
		"starting dynamic-user-segmentation service",
		slog.String("env", cfg.Env),
	)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	db, err := repository.NewDB()
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
	}

	repo := repository.NewRepository(db)
	controllers := controller.NewController(repo, log)
	controllers.InitRoutes(router)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case localEnv:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case demoEnv:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
