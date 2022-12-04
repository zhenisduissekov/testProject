package connection

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
	"github.com/rs/zerolog/log"
)

func migrateDB(pool *pgxpool.Pool, schemeTable, path string) (err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to acquire connection")
		return err
	}

	migrator, err := migrate.NewMigrator(context.Background(), conn.Conn(), schemeTable)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create migrator")
		return err
	}
	err = migrator.LoadMigrations(path)
	if err != nil {
		log.Error().Err(err).Msg("Unable to load migrations")
		return err
	}

	err = migrator.Migrate(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to migrate")
		return err
	}
	ver, err := migrator.GetCurrentVersion(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Unable to get current version")
		return err
	}
	log.Info().Msgf("Current version: %d", ver)
	return nil
}
