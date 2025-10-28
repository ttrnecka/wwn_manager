package utils

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func LoadEnv(logger *zerolog.Logger) {
	envFile := filepath.Join(BinaryOrBuildDir(), ".env")

	envMap, err := godotenv.Read(envFile)
	if err != nil {
		logger.Info().Msg("No .env file found")
	} else {
		logger.Info().Msg(".env file loaded")
	}

	// Only set env vars that aren't already set
	for k, v := range envMap {
		if _, exists := os.LookupEnv(k); !exists {
			err := os.Setenv(k, v)
			if err != nil {
				logger.Fatal().Err(err).Msg("")
			}
		}
	}
}
