// Custom Middlware for RBAC and Authorization
package middleware

import (
	"context"
	"net/http"

	"time"

	"github.com/rs/zerolog/log"

	"server/authorization"

	"server/helpers"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = authorization.InitJWTAuth()
}

// Extracts the claims from the JWT token and adds them to the context
// Returns 401 if token is expired
func RBACMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		ctx := context.WithValue(r.Context(), "claims", claims)

		exp := claims["exp"].(time.Time).Unix()
		now := time.Now().Unix()
		// log.Info().Msgf("RBACMiddleware: exp=%v\n", exp)
		// log.Info().Msgf("RBACMiddleware: now=%v\n", now)

		if exp < now {
			log.Info().Msgf("RBACMiddleware: token expired\n")

			un := struct {
				Error   bool   `json:"error"`
				Message string `json:"message"`
			}{
				Error:   true,
				Message: "Token expired.",
			}

			helpers.WriteJSON(w, http.StatusUnauthorized, un)
			return
		}

		scope := claims["app_metadata"].(map[string]interface{})["authorization"] //.(map[string]interface{})["roles"]

		log.Info().Msgf("RBACMiddleware: scope=%v\n", scope)

		var scopeArray []string

		//extract roles from scope
		for key, value := range scope.(map[string]interface{}) {

			valueArray := value.([]interface{})

			if key == "roles" {
				for _, v := range valueArray {
					scopeArray = append(scopeArray, v.(string))
				}
			}
		}

		log.Info().Msgf("RBACMiddleware: scopeArray=%v\n", scopeArray)

		ctx = context.WithValue(ctx, "scope", scopeArray)

		//TODO: Add other claims to context

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Checks if the user has the required scope to access the route
func RBACMiddlewareProtectedRoute(scopeRequired string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scope := r.Context().Value("scope").([]string)
			// scope := r.Context().Value("scope").([]interface{})

			log.Info().Msgf("RBACMiddlewareProtectedRoute: scope=%v\n", scope)
			log.Info().Msgf("RBACMiddlewareProtectedRoute: scopeRequired=%v\n", scopeRequired)

			if !helpers.Contains(scope, scopeRequired) {
				log.Info().Msgf("RBACMiddlewareProtectedRoute: scopeRequired=%v not found in scope=%v\n", scopeRequired, scope)

				un := struct {
					Error   bool   `json:"error"`
					Message string `json:"message"`
				}{
					Error:   true,
					Message: "You do not have the required scope to access this resource.",
				}

				helpers.WriteJSON(w, http.StatusForbidden, un)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
