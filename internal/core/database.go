package core

import (
	"github.com/rs/zerolog/log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func MustSetupMySQLDatabase(c AppConfig) *sqlx.DB {
	dbURL := strings.TrimPrefix(c.DbURL, "mysql://")
	db, err := sqlx.Connect("mysql", dbURL)
	if err != nil {
		log.Panic().Err(err).Msg("failed to connect to MySQL database")
	}
	return db
}
