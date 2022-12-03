package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/zhenisduissekov/testProject/pkg/connection"
	"github.com/zhenisduissekov/testProject/pkg/cryptocompare"
	"github.com/zhenisduissekov/testProject/pkg/scheduler"
)

const (
	filePath = "pkg/config/messages.yaml"
)

type Config struct {
	DB            connection.DBConfig
	HTTP          HTTPConfig
	Scheduler     scheduler.SchedulerConfig
	CryptoCompare cryptocompare.CryptoCompareConfig
	LogLevel      string
}

type HTTPConfig struct {
	Host string
	Port string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading env file")
	}
}

func New() *Config {
	return &Config{
		DB: connection.DBConfig{
			Service:         "testProject",
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Pass:            getEnv("DB_PASS", "postgres"),
			DB:              getEnv("DB_NAME", "postgres"),
			TimeOut:         getEnv("DB_TIMEOUT", "10"),
			MigrationPath:   getEnv("DB_MIGRATION_PATH", "pkg/migrations"),
			MigrationScheme: getEnv("DB_MIGRATION_SCHEME", "schema_migrations"),
		},
		HTTP: HTTPConfig{
			Host: getEnv("HTTP_HOST", "localhost"),
			Port: getEnv("HTTP_PORT", "8080"),
		},
		Scheduler: scheduler.SchedulerConfig{
			Interval: getEnvAsInt("SCHEDULER_INTERVAL_SECONDS", 10),
		},
		CryptoCompare: cryptocompare.CryptoCompareConfig{
			APIKey:       getEnv("CRYPTOCOMPARE_API_KEY", "https://min-api.cryptocompare.com/data/pricemultifull"),
			TimeOut:      getEnvAsInt("CRYPTOCOMPARE_TIMEOUT", 10),
			FSYMS:        getEnvSYMS("FSYMS", "BTC,ETH"),
			TSYMS:        getEnvSYMS("TSYMS", "USD"),
			WaitInterval: getEnvAsInt("CRYPTOCOMPARE_INTERVAL", 30),
		},
		LogLevel: strings.ToLower(getEnv("LOG_LEVEL", "info")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		num, err := strconv.Atoi(value)
		if err != nil {
			log.Err(err).Msg("error converting env variable to int")
			return fallback
		}
		return num
	}
	return fallback
}

func getEnvSYMS(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		value = strings.ReplaceAll(value, "[", "")
		value = strings.ReplaceAll(value, "]", "")
		value = strings.ReplaceAll(value, "\"", "")
		value = strings.ReplaceAll(value, " ", "")
		return value
	}
	return fallback
}
