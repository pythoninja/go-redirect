package server

import (
	"fmt"
	"github.com/pythoninja/go-redirect/internal/api"
	"github.com/pythoninja/go-redirect/internal/config"
	"github.com/pythoninja/go-redirect/internal/storage"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func Serve(app *config.Application, store *storage.Storage) {
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	serverLogger := slog.NewLogLogger(logHandler, slog.LevelDebug)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      api.Routes(app, store),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorLog:     serverLogger,
	}

	slog.Info("starting server", slog.Any("addr", srv.Addr))
	slog.Info("database info", slog.Any("dsn", app.Config.Database.Dsn))

	err := srv.ListenAndServe()
	if err != nil {
		slog.Error(err.Error())
	}
}
