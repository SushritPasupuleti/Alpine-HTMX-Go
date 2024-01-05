// Functions for hashing and comparing passwords
package helpers

import (
	// "fmt"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// hashes and salts password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return "", err
	}

	// log.Info().Msg("Password hashed successfully")
	// log.Info().Msg(fmt.Sprintf("Hashed password: %s", string(hash)))

	return string(hash), nil
}

// compares password with hashed password in database
func ComparePasswords(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Error().Err(err).Msg("Error comparing passwords")
		return false
	}

	// log.Info().Msg("Password matched successfully")

	return true
}
