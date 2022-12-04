package repository

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
)

func (r Repository) Read(fsyms, tsyms string) (map[string]string, error) {
	log.Trace().Msgf("Reading quote for fsyms: [%s] and tsyms: [%s]", fsyms, tsyms)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("raw", "display").From(tableName).Where(sq.Eq{"fsyms": fsyms, "tsyms": tsyms}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error while building sql: %w", err)
	}
	log.Trace().Msgf("Executing sql: [%s] with args: [%v]", sql, args)
	return r.connection.Read(sql, args)
}
