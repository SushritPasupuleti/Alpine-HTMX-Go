// Handles loading and validation of environment variables
package env

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"errors"
	"os"
)

// Config struct with environment variables
type Config struct {
	PORT        string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	JWT_SECRET  string
	ENVIRONMENT string
	REDIS_HOST  string
	REDIS_PORT  string
}

var DefaultConfig Config

// Load and validate environment variables
// `DefaultConfig` is now available to use
func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error loading .env file")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal().
			Err(errors.New("$PORT must be set")).
			Msg("$PORT must be set")
		os.Exit(1)
	}

	db_host := os.Getenv("DB_HOST")
	if db_host == "" {
		log.Fatal().
			Err(errors.New("$DB_HOST must be set")).
			Msg("$DB_HOST must be set")
		os.Exit(1)
	}

	db_port := os.Getenv("DB_PORT")
	if db_port == "" {
		log.Fatal().
			Err(errors.New("$DB_PORT must be set")).
			Msg("$DB_PORT must be set")
		os.Exit(1)
	}

	db_user := os.Getenv("DB_USER")
	if db_user == "" {
		log.Fatal().
			Err(errors.New("$DB_USER must be set")).
			Msg("$DB_USER must be set")
		os.Exit(1)
	}

	db_password := os.Getenv("DB_PASSWORD")
	if db_password == "" {
		log.Fatal().
			Err(errors.New("$DB_PASSWORD must be set")).
			Msg("$DB_PASSWORD must be set")
		os.Exit(1)
	}

	db_name := os.Getenv("DB_NAME")
	if db_name == "" {
		log.Fatal().
			Err(errors.New("$DB_NAME must be set")).
			Msg("$DB_NAME must be set")
		os.Exit(1)
	}

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal().
			Err(errors.New("$JWT_SECRET must be set")).
			Msg("$JWT_SECRET must be set")
		os.Exit(1)
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		log.Fatal().
			Err(errors.New("$ENVIRONMENT must be set")).
			Msg("$ENVIRONMENT must be set")
		os.Exit(1)
	}

	redis_host := os.Getenv("REDIS_HOST")
	if redis_host == "" {
		log.Fatal().
			Err(errors.New("$REDIS_HOST must be set")).
			Msg("$REDIS_HosT must be set")
		os.Exit(1)
	}

	redis_port := os.Getenv("REDIS_PORT")
	if redis_port == "" {
		log.Fatal().
			Err(errors.New("$REDIS_PORT must be set")).
			Msg("$REDIS_PORT must be set")
		os.Exit(1)
	}

	DefaultConfig = Config{
		PORT:        port,
		DB_HOST:     db_host,
		DB_PORT:     db_port,
		DB_USER:     db_user,
		DB_PASSWORD: db_password,
		DB_NAME:     db_name,
		JWT_SECRET:  jwt_secret,
		ENVIRONMENT: environment,
		REDIS_HOST:  redis_host,
		REDIS_PORT:  redis_port,
	}

	// log.Info().Msgf("Successfully loaded environment variables: %v", DefaultConfig)
	log.Info().Msg("Successfully loaded environment variables")

	return DefaultConfig, nil
}
