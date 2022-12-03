package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(level string) zerolog.Level {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Trace().Msgf("level: %s", level)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	log.Logger = log.With().Caller().Logger()
	switch {
	case level == "warn":
		return zerolog.WarnLevel
	case level == "panic":
		return zerolog.PanicLevel
	case level == "fatal":
		return zerolog.FatalLevel
	case level == "error":
		return zerolog.ErrorLevel
	case level == "info":
		return zerolog.InfoLevel
	case level == "trace":
		return zerolog.TraceLevel
	case level == "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel
	}
}
