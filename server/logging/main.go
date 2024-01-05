// Contains the configuration for the logger.
// Utilizes the `zerolog` package.
package logging

import (
	"errors"
	"os"
	"server/env"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Initializes the logger with enironment-specific settings
func InitLogging() {
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger() // Enable Caller Info

	ENVIRONMENT := env.DefaultConfig.ENVIRONMENT
	if ENVIRONMENT == "" {
		log.Fatal().
			Err(errors.New("$ENVIRONMENT must be set")).
			Msg("$ENVIRONMENT must be set")
	}

	// Use Pretty Console Logging in Development
	if ENVIRONMENT == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Msg("Using pretty console logging")
	}
}
