package main

import (
	// "errors"
	"fmt"
	// "log"
	"net/http"
	// "os"

	"github.com/redis/go-redis/v9"
	"server/db"
	"server/env"
	"server/logging"
	"server/models"
	rc "server/redis"
	"server/routes"

	"github.com/rs/zerolog/log"
)

type Application struct {
	Config env.Config
	Models models.Models
	Redis  *redis.Client
}

func (app *Application) Serve() error {
	port := app.Config.PORT

	log.Printf("🚀 Server listening on port %s", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: routes.Routes(),
	}

	return srv.ListenAndServe()
}

func init() {
	env.Load()
	logging.InitLogging()
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://localhost:6000/terms/

// @contact.name API Support
// @contact.url http://www.localhost:6000/support
// @contact.email support@localhost:6000

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /
func main() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.DefaultConfig.DB_HOST, env.DefaultConfig.DB_PORT, env.DefaultConfig.DB_USER, env.DefaultConfig.DB_PASSWORD, env.DefaultConfig.DB_NAME)

	dbConn, err := db.Connect(dsn)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error connecting to database")
	}

	defer dbConn.DB.Close()

	redisClient, err := rc.InitRedisClient()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error connecting to Redis")
	}

	defer redisClient.Close()

	app := Application{
		Config: env.DefaultConfig,
		Models: models.New(dbConn.DB),
		Redis:  redisClient,
	}

	err = app.Serve()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error starting server")
	}
}

