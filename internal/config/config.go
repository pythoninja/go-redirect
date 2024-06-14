package config

import (
	"github.com/pythoninja/go-redirect/internal/generator"
	"log/slog"
	"os"
	"strconv"
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

func SetupLogger(env string) {
	const (
		envDev  = "development"
		envProd = "production"
	)

	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	log.Info("initialized default logger")
	slog.SetDefault(log)
}

// GetEnv retrieves the value of an environment variable corresponding to the given key.
// If the environment variable is not found, it returns the fallback value provided.
//
// [os.LookupEnv] always returns a string, but real return type depends on the fallback argument.
func GetEnv[T bool | string | int](key string, fallback T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	var parsedValue any

	switch any(fallback).(type) {
	case string:
		parsedValue = value
	case int:
		var intValue int

		intValue, err := strconv.Atoi(value)
		if err != nil {
			slog.Error("failed to convert env variable to integer", slog.Any("error", err.Error()))

			return fallback
		}

		parsedValue = intValue
	case bool:
		var boolValue bool

		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			slog.Error("failed to convert env variable to boolean", slog.Any("error", err.Error()))

			return fallback
		}

		parsedValue = boolValue
	default:
		return fallback
	}

	if res, ok := parsedValue.(T); ok {
		return res
	}

	return fallback
}
