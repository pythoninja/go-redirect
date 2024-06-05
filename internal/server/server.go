package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/pythoninja/go-redirect/internal/api"
	"github.com/pythoninja/go-redirect/internal/config"
	"github.com/pythoninja/go-redirect/internal/storage"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(app *config.Application, store *storage.Storage) error {
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	serverLogger := slog.NewLogLogger(logHandler, slog.LevelDebug)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.Config.Addr, app.Config.Port),
		Handler:      api.Router(app, store),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorLog:     serverLogger,
	}

	shutdownError := make(chan error)

	go func() {
		slog.Info("starting background job to listen for signals")
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		slog.Info("shutting down the server", slog.Any("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	slog.Info("starting server", slog.Any("addr", srv.Addr))
	slog.Info("api secret key", slog.Any("key", app.Config.APISecretKey))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	slog.Info("server gracefully stopped")

	return nil
}
