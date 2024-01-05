// Description: This file contains the main authentication logic for the server.
package authentication

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/env"
	"server/helpers"
	"server/models"
	"server/redis"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

var userModel models.User

type UserAuth struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Scope    string `json:"scope,omitempty"`
}

// Generates a JWT token for the user
func GenerateToken(w http.ResponseWriter, r *http.Request) {
	var user UserAuth

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Error().Err(err).Msg("Error decoding user")
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(user)

	var validationErrors []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Error().Err(err).Msg("Error validating user")

			validationErrors = append(validationErrors, err.Error())
		}

		if len(validationErrors) > 0 {
			var errorMessages string

			for _, errorMessage := range validationErrors {
				errorMessages += errorMessage + "\n"
			}

			helpers.ErrorJSON(w, errors.New(errorMessages), http.StatusBadRequest)
			return
		}
	}

	current_user, err := userModel.FindByEmail(user.UserName)
	if err != nil {
		log.Error().Err(err).Msg("Error finding user")
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//validate user credentials
	verified := helpers.ComparePasswords(current_user.Password, user.Password)
	if !verified {
		helpers.ErrorJSON(w, errors.New("Invalid Credentials Passed"), http.StatusBadRequest)
		return
	}

	// roles := []string{current_user.UserRole}

	if verified {
		//create token
		var token models.JWTClaims = models.JWTClaims{
			Email: user.UserName,
			AppMetadata: models.AppMetadata{
				Authorization: models.Authorization{
					// Roles: roles,
				},
			},
			Subject:    user.UserName,
			Audience:   "HOST", //TODO: Add audience from env
			Expiration: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:   time.Now().Unix(),
			//TODO: Generate JWTID from database
		}

		log.Info().Msgf("token: %v", token)

		//Create JWT token
		jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))

		jwtToken.Claims = token

		signedToken, err := jwtToken.SignedString([]byte(env.DefaultConfig.JWT_SECRET))

		if err != nil {
			log.Error().Err(err).Msg("Error signing token")
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		log.Info().Msgf("signedToken: %v", signedToken)

		jti := uuid.Must(uuid.NewV4()).String()

		refreshToken := jwt.New(jwt.GetSigningMethod("HS256"))
		rtClaims := refreshToken.Claims.(jwt.MapClaims)
		rtClaims["sub"] = user.UserName
		rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
		rtClaims["jti"] = jti

		rt, err := refreshToken.SignedString([]byte(env.DefaultConfig.JWT_SECRET))
		if err != nil {
			log.Error().Err(err).Msg("Error signing refresh token")
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		log.Info().Msgf("refreshToken: %v", rt)

		//save JTI to redis
		// err = redis.SetCache(jti, "true", rtClaims["exp"].(time.Duration))
		err = redis.SetCache(user.UserName, jti, time.Hour*24*7)

		if err != nil {
			log.Error().Err(err).Msg("Error saving JTI to redis")
			helpers.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}

		response := struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{
			AccessToken:  signedToken,
			RefreshToken: rt,
		}

		_ = helpers.WriteJSON(w, http.StatusOK, response)
		return
	}

	helpers.ErrorJSON(w, errors.New("Invalid Credentials Passed"), http.StatusBadRequest)
}

// Refreshes a JWT token for the user
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	//extract refresh token from request Header
	refreshToken := r.Header.Get("Authorization")

	if refreshToken == "" {
		helpers.ErrorJSON(w, errors.New("Refresh token not provided"), http.StatusBadRequest)
		return
	}

	refreshToken = refreshToken[7:]

	//validate refresh token
	refreshTokenClaims, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.DefaultConfig.JWT_SECRET), nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Error parsing refresh token")
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if !refreshTokenClaims.Valid {
		helpers.ErrorJSON(w, errors.New("Invalid refresh token"), http.StatusBadRequest)
		return
	}

	//check if JTI is in redis
	jti, err := redis.GetCache(refreshTokenClaims.Claims.(jwt.MapClaims)["sub"].(string))

	if err != nil {
		log.Error().Err(err).Msg("Error getting JTI from redis")
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if jti == "" || jti != refreshTokenClaims.Claims.(jwt.MapClaims)["jti"].(string) {
		helpers.ErrorJSON(w, errors.New("Invalid refresh token"), http.StatusBadRequest)
		return
	}

	//get user from database
	user, err := userModel.FindByEmail(refreshTokenClaims.Claims.(jwt.MapClaims)["sub"].(string))
	if err != nil {
		log.Error().Err(err).Msg("Error finding user")
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// roles := []string{user.UserRole}

	//create token
	var token models.JWTClaims = models.JWTClaims{
		Email: user.Email,
		AppMetadata: models.AppMetadata{
			Authorization: models.Authorization{
				// Roles: roles,
			},
		},
		Subject:    user.Email,
		Audience:   "HOST", //TODO: Add audience from env
		Expiration: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:   time.Now().Unix(),
	}

	log.Info().Msgf("token: %v", token)

	//Create JWT token
	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))

	jwtToken.Claims = token

	signedToken, err := jwtToken.SignedString([]byte(env.DefaultConfig.JWT_SECRET))

	if err != nil {
		log.Error().Err(err).Msg("Error signing token")
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	log.Info().Msgf("signedToken: %v", signedToken)

	response := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: signedToken,
	}

	_ = helpers.WriteJSON(w, http.StatusOK, response)
}

// Revokes a JWT token for the user
func RevokeToken(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email string `json:"email" validate:"required"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Error().Err(err).Msg("Error decoding body")
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	err = validate.Struct(body)

	var validationErrors []string

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Error().Err(err).Msg("Error validating user")

			validationErrors = append(validationErrors, err.Error())
		}

		if len(validationErrors) > 0 {
			var errorMessages string

			for _, errorMessage := range validationErrors {
				errorMessages += errorMessage + "\n"
			}

			helpers.ErrorJSON(w, errors.New(errorMessages), http.StatusBadRequest)
			return
		}
	}

	//revoke token
	err = redis.DeleteCache(body.Email)

	if err != nil {
		log.Error().Err(err).Msg("Error revoking token")
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = helpers.WriteJSON(w, http.StatusOK, "Token revoked successfully")
}
