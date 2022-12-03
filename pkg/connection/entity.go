package connection

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
	"github.com/rs/zerolog/log"
)

const (
	connectionStringTemplate = `host=%s port=%s user=%s password=%s dbname=%s application_name=%s sslmode=disable`
)

type DBConfig struct {
	Service         string
	Host            string
	Port            string
	DB              string
	User            string
	Pass            string
	TimeOut         string
	MigrationPath   string
	MigrationScheme string
}

type Client struct {
	ctx  *context.Context
	pool *pgxpool.Pool
}

func New(cnf DBConfig) Client {
	ctx := context.TODO()

	connectionString := fmt.Sprintf(connectionStringTemplate, cnf.Host, cnf.Port, cnf.User, cnf.Pass, cnf.DB, cnf.Service)
	log.Trace().Msgf("connection string %s", connectionString)

	pool, err := pgxpool.Connect(ctx, connectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to DB")
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("could not ping")
	}

	err = migrateDB(pool, "schema_migrations", cnf.MigrationPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not migrate DB")
	}

	return Client{
		ctx:  &ctx,
		pool: pool,
	}
}

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

func (c *Client) Save(query string, args []interface{}) (err error) {
	resp, err := c.pool.Exec(*c.ctx, query, args...)
	if err != nil {
		return fmt.Errorf("could not exec: %w", err)
	}
	log.Info().Msgf("rows affected: %s", resp)
	return nil
}

func (c *Client) Read(query string, args []interface{}) (data string, err error) {
	rows, err := c.pool.Query(*c.ctx, query, args...)
	if err != nil {
		return "", fmt.Errorf("could not Query: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&data)
		if err != nil {
			return "", fmt.Errorf("could not scan: %w", err)
		}
	}

	return data, nil
}
