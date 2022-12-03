package repository

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/internal/connection"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "crypto.quotes"
)

type IndexedMessages struct {
	Index     int    `json:"index"`
	Number    int32  `json:"number"`
	Text      string `json:"text"`
	Status    string `json:"status"`
	CreatedOn string `json:"createdOn"`
	CreatedBy string `json:"createdBy"`
}

type Repository struct {
	connection connection.Client
}

func New(conn connection.Client) Repository {
	return Repository{
		connection: conn,
	}
}

func (r Repository) Save(raw, display []byte, fsyms, tsyms string, error error) error {
	error_message := "nil"
	if error != nil {
		error_message = error.Error()
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Insert(tableName).Columns("fsyms", "tsyms", "raw", "display", "error_message").
		Values(fsyms, tsyms, string(raw), string(display), error_message).
		Suffix("ON CONFLICT (fsyms, tsyms) DO UPDATE SET raw = EXCLUDED.raw, display = EXCLUDED.display, error_message = EXCLUDED.error_message, updated_at=Now()").
		ToSql()
	if err != nil {
		return fmt.Errorf("error while building sql: %w", err)
	}

	log.Debug().Msgf("Executing sql: [%s] with args: [%v]", sql, args[0])
	return r.connection.Save(sql, args)
}

func (r Repository) Get(fsyms, tsyms string) (string, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("raw", "display").From(tableName).Where(sq.Eq{"fsyms": fsyms, "tsyms": tsyms}).ToSql()
	if err != nil {
		return "", err
	}
	log.Debug().Msgf("Executing sql: [%s] with args: [%v]", sql, args[0])
	return r.connection.Read(sql, args)
}
