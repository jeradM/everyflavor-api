package cmd

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type migrationConfig struct {
	DbUrl            string
	MigrationVersion int
}

type migrationLogger struct{}

func (m migrationLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (m migrationLogger) Verbose() bool {
	return true
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		var config migrationConfig
		if err := viper.Unmarshal(&config); err != nil {
			panic(errors.Wrap(err, "unable to parse config"))
		}
		fmt.Println(config.DbUrl)
		m, err := migrate.New("file://migrations/mysql", "mysql://"+config.DbUrl)
		if err != nil {
			log.Panic().Err(err).Msg("migration failed")
		}
		m.Log = migrationLogger{}
		if config.MigrationVersion == -1 {
			err = m.Up()
		} else if config.MigrationVersion == 0 {
			err = m.Down()
		} else if config.MigrationVersion > 0 {
			err = m.Migrate(uint(config.MigrationVersion))
		} else {
			log.Panic().Int("version", config.MigrationVersion).Msg("Invalid migration version")
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Panic().Err(err).Msg("migration failed")
		}

		var msg string
		if err != nil {
			msg = "nothing to migrate"
		} else {
			msg = "migrations complete"
		}
		log.Debug().Int("version", config.MigrationVersion).Msg(msg)
	},
}
