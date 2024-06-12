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

	return &Application{
		Config: *cfg,
	}
}

func InitLogger() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	log.Info("initialized default logger")
	slog.SetDefault(log)
}
