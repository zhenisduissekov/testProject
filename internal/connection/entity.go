package connection

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

type DB interface {
	Ping(ctx context.Context) error
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Client struct {
	ctx  *context.Context
	pool DB
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
