package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"os"
	"regexp"
	"testing"
)

func testLogger() {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	slog.SetDefault(log)
}

//nolint:funlen
//revive:disable function-length
func TestGetEnv(t *testing.T) {
	testLogger()

	cases := []struct {
		name     string
		key      string
		fallback any
		setValue string
		want     any
	}{
		{
			name:     "IntegerFallback",
			key:      "TEST_ENV_INT",
			fallback: 0,
			setValue: "123",
			want:     123,
		},
		{
			name:     "BoolFallback",
			key:      "TEST_ENV_BOOL",
			fallback: false,
			setValue: "true",
			want:     true,
		},
		{
			name:     "StringFallback",
			key:      "TEST_ENV_STR",
			fallback: "fallback",
			setValue: "env_value",
			want:     "env_value",
		},
		{
			name:     "InvalidIntegerParse",
			key:      "TEST_ENV_INVALID_INT",
			fallback: 0,
			setValue: "invalid_int",
			want:     0,
		},
		{
			name:     "InvalidBoolParse",
			key:      "TEST_ENV_INVALID_BOOL",
			fallback: false,
			setValue: "invalid_bool",
			want:     false,
		},
		{
			name:     "StringFallbackWithoutEnvSet",
			key:      "TEST_ENV_RANDOM_NON_EXISTENT",
			fallback: "fallback",
			want:     "fallback",
		},
		{
			name:     "IntFallbackWithoutEnvSet",
			key:      "TEST_ENV_RANDOM_NON_EXISTENT",
			fallback: 1,
			want:     1,
		},
		{
			name:     "BoolFallbackWithoutEnvSet",
			key:      "TEST_ENV_RANDOM_NON_EXISTENT",
			fallback: true,
			want:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setValue != "" {
				err := os.Setenv(tc.key, tc.setValue)
				require.NoError(t, err)
			}

			switch fallbackType := tc.fallback.(type) {
			case int:
				val := GetEnv(tc.key, fallbackType)
				assert.Equal(t, tc.want, val)
			case bool:
				val := GetEnv(tc.key, fallbackType)
				assert.Equal(t, tc.want, val)
			case string:
				val := GetEnv(tc.key, fallbackType)
				assert.Equal(t, tc.want, val)
			}

			err := os.Unsetenv(tc.key)
			require.NoError(t, err)
		})
	}
}

func TestInitConfiguration(t *testing.T) {
	testLogger()

	cfg := &Config{}

	t.Run("EmptySecretKeyReturnsNewRandomKey", func(t *testing.T) {
		cfg.APISecretKey = ""
		expectedPattern := "^[A-Z0-9]{52}$"

		app := InitConfiguration(cfg)

		matched, err := regexp.MatchString(expectedPattern, app.Config.APISecretKey)
		require.NoError(t, err)

		assert.NotNil(t, app, "app should be non nil after initialization")
		assert.True(t, matched, "api key should match the pattern: %s", expectedPattern)
	})

	t.Run("PredefinedSecretKeyReturnsSelf", func(t *testing.T) {
		cfg.APISecretKey = "K123"
		expectedValue := "K123"

		app := InitConfiguration(cfg)

		assert.Equal(t, expectedValue, app.Config.APISecretKey)
	})
}
