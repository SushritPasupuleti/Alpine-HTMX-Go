// Root router of the application
package routes

import (
	"fmt"
	"html/template"
	"time"

	// "log"
	"net/http"
	// "os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"github.com/unrolled/secure"

	// "github.com/go-chi/oauth"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	authentication "server/authentication"
	"server/authorization"
	"server/env"
	"server/handlers"
	"server/helpers"
	middlewareCustom "server/middleware"

	// "server/models"
	_ "server/docs"

	"github.com/rs/zerolog/log"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = authorization.InitJWTAuth()
}

// Returns a router with all routes configured
func Routes() http.Handler {

	//INFO: Refer [to](https://github.com/unrolled/secure?tab=readme-ov-file#default-options)
	secureMiddleware := secure.New(secure.Options{
		IsDevelopment:      env.DefaultConfig.ENVIRONMENT == "development",
		AllowedHosts:       []string{"localhost:5000", "localhost:3000"},
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		// ContentSecurityPolicy: "default-src 'self'",
		// This allows htmx's script to be loaded from unpkg.com
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' https://unpkg.com 'nonce-a23gbfz9e'; style-src 'self';",
		// ContentSecurityPolicy: "default-src 'self'; script-src 'self' https://unpkg.com https://cdn.jsdelivr.net;",
	})

	router := chi.NewRouter()
	router.Use(secureMiddleware.Handler)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// OAuth2 INFO: Should be run as a separate service in production
	router.Group(func(r chi.Router) {
		router.Route("/oauth", func(r chi.Router) {
			r.Post("/token", authentication.GenerateToken)
			r.Get("/token/refresh", authentication.RefreshToken)
			r.Post("/token/revoke", authentication.RevokeToken)
		})
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// if request header accepts text/html, then return a html template
		if helpers.SupportsHTML(r) {
			log.Info().Msg("Accepts text/html")
			w.Header().Set("Content-Type", "text/html")

			t, _ := template.ParseFiles("templates/index.html")
			err := t.Execute(w, nil)
			if err != nil {
				log.Error().Err(err).Msg("Error executing template")
				w.Write([]byte("Error executing template"))
				return
			}
			return
		}

		w.Write([]byte("API is up and running"))
	})

	//serve css
	router.Get("/dist/main.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "templates/main.css")
	})

	//serve js
	router.Get("/dist/main.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		http.ServeFile(w, r, "templates/main.js")
	})

	//serve login page
	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/login.html")
		err := t.Execute(w, nil)

		if err != nil {
			log.Error().Err(err).Msg("Error executing template")
			w.Write([]byte("Error executing template"))
			return
		}
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:5000/swagger/doc.json"), //The url pointing to API definition
	))

	router.Group(func(r chi.Router) {
		router.Route("/api/v1/users", func(r chi.Router) {
			// r.Get("/", handlers.GetAllUsers)
			r.With(middlewareCustom.CacheMiddleware(0)).Get("/", handlers.GetAllUsers) //Response is cached
			r.Post("/", handlers.CreateUser)
			r.Get("/{email}", handlers.FindUserByEmail)
			r.Put("/", handlers.UpdateUserByEmail)

			//Rate limit by IP for 3 requests per 30 minutes
			r.With(httprate.LimitByIP(3, 30*time.Minute)).Post("/check-password", handlers.CheckUserPassword)

		})
	})

	// Protected routes
	router.Group(func(r chi.Router) {
		router.Route("/api/v1/admin", func(r chi.Router) {
			//1. Verify token
			//2. Authenticate token
			//3. Populate roles from token into context
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator) //TODO: Add custom Authenticator to verify token against DB
			r.Use(middlewareCustom.RBACMiddleware)

			r.With(middlewareCustom.RBACMiddlewareProtectedRoute("admin")).Get("/", func(w http.ResponseWriter, r *http.Request) {
				_, claims, _ := jwtauth.FromContext(r.Context())
				w.Write([]byte(fmt.Sprintf("Hello, %v you are authorized to view this.", claims["user_id"])))
			})
		})
	})

	return router
}
