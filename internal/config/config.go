package config

import (
	"github.com/pythoninja/go-redirect/internal/generator"
	"log/slog"
	"os"
	"strings"
)

type Application struct {
	Config Config
}

type Config struct {
	Addr     string
	Port     int
	Env      string
	Database struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	EnableRateLimiter bool
	APISecretKey      string
}

func InitConfiguration(cfg *Config) *Application {
	if strings.TrimSpace(cfg.APISecretKey) == "" {
		cfg.APISecretKey = generator.NewAPIKey()
	}

	slog.Info("starting application",
		slog.Group("settings",
			slog.String("env", cfg.Env),
			slog.String("addr", cfg.Addr),
			slog.Int("port", cfg.Port),
			slog.String("db-dsn", cfg.Database.Dsn),
			slog.Int("db-max-open-conns", cfg.Database.MaxOpenConns),
			slog.Int("db-max-idle-conns", cfg.Database.MaxIdleConns),
			slog.String("db-max-idle-time", cfg.Database.MaxIdleTime),
			slog.Bool("rate-limiter-enabled", cfg.EnableRateLimiter),
		))

	return &Application{
		Config: *cfg,
	}
}

func InitLogger() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	log.Info("initialize default logger")
	slog.SetDefault(log)
}
