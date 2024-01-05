// Handles Database connection establishment
// and Session management
package db

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"time"

	// "fmt"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type DB struct {
	DB *sql.DB
}

var dbConn = &DB{}

const maxOpenConns = 10
const maxIdleConns = 5
const maxDBLifetime = 5 * time.Minute
const driverName = "pgx"

func Connect(dsn string) (*DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxDBLifetime)

	// err = db.Ping()

	err = PingDB(db)
	if err != nil {
		return nil, err
	}

	dbConn.DB = db

	return dbConn, nil
}

func PingDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	log.Info().Msgf("Successfully connected to database %d", db.Stats().InUse)
	return nil
}
