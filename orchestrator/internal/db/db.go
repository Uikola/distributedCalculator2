package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	"os"
)

func InitPostgres(pgURL string) *sqlx.DB {
	db, err := sqlx.Connect("pgx", pgURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to the db")
		os.Exit(1)
	}

	return db
}

func OtelInitPostgres(pgURL string) (*sqlx.DB, error) {
	db, err := otelsqlx.Connect("pgx", pgURL)
	if err != nil {
		return nil, err
	}
	return db, err
}
