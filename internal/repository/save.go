package repository

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
)

// saving data to DB repository
func (r Repository) Save(raw, display []byte, fsyms, tsyms string, error error) error {
	log.Trace().Msgf("Saving quote for fsyms: [%s] and tsyms: [%s]", fsyms, tsyms)
	errorMsg := "nil"
	if error != nil {
		errorMsg = error.Error()
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Insert(tableName).Columns("fsyms", "tsyms", "raw", "display", "error_message").
		Values(fsyms, tsyms, string(raw), string(display), errorMsg).
		Suffix("ON CONFLICT (fsyms, tsyms) DO UPDATE SET raw = EXCLUDED.raw, display = EXCLUDED.display, error_message = EXCLUDED.error_message, updated_at=Now()").
		ToSql()
	if err != nil {
		return fmt.Errorf("error while building sql: %w", err)
	}

	log.Debug().Msgf("Executing sql: [%s] with args: [%v]", sql, args[0])
	return r.connection.Save(sql, args)
}
