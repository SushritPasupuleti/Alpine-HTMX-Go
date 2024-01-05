package authorization

import (
	"log"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
)

// Initializes a JWTAuth instance with the secret from the .env file
// Returns a pointer to the JWTAuth instance
// Ensure that a `JWT_SECRET` environment variable is set
func InitJWTAuth() *jwtauth.JWTAuth {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("$JWT_SECRET must be set")
		os.Exit(1)
	}

	return jwtauth.New("HS256", []byte(jwtSecret), nil)
}
