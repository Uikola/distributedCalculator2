package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func InitDB(pgURL string) *sqlx.DB {
	db, err := sqlx.Connect("pgx", pgURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to the database")
	}
	return db
}
