package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/coneno/logger"
	"github.com/tekenradar/content-service/pkg/types"
)

// Config is the structure that holds all global configuration data
type Config struct {
	Port            string
	LogLevel        logger.LogLevel
	ContentDBConfig types.DBConfig
}

func InitConfig() Config {
	conf := Config{}
	conf.Port = os.Getenv("CONTENT_SERVICE_LISTEN_PORT")

	conf.LogLevel = getLogLevel()
	conf.ContentDBConfig = getContentDBConfig()

	return conf
}

func getLogLevel() logger.LogLevel {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return logger.LEVEL_DEBUG
	case "info":
		return logger.LEVEL_INFO
	case "error":
		return logger.LEVEL_ERROR
	case "warning":
		return logger.LEVEL_WARNING
	default:
		return logger.LEVEL_INFO
	}
}

func getContentDBConfig() types.DBConfig {
	connStr := os.Getenv("CONTENT_DB_CONNECTION_STR")
	username := os.Getenv("CONTENT_DB_USERNAME")
	password := os.Getenv("CONTENT_DB_PASSWORD")
	prefix := os.Getenv("CONTENT_DB_CONNECTION_PREFIX") // Used in test mode
	if connStr == "" || username == "" || password == "" {
		logger.Error.Fatal("Couldn't read DB credentials.")
	}
	URI := fmt.Sprintf(`mongodb%s://%s:%s@%s`, prefix, username, password, connStr)

	var err error
	Timeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	if err != nil {
		logger.Error.Fatal("DB_TIMEOUT: " + err.Error())
	}
	IdleConnTimeout, err := strconv.Atoi(os.Getenv("DB_IDLE_CONN_TIMEOUT"))
	if err != nil {
		logger.Error.Fatal("DB_IDLE_CONN_TIMEOUT" + err.Error())
	}
	mps, err := strconv.Atoi(os.Getenv("DB_MAX_POOL_SIZE"))
	MaxPoolSize := uint64(mps)
	if err != nil {
		logger.Error.Fatal("DB_MAX_POOL_SIZE: " + err.Error())
	}

	DBNamePrefix := os.Getenv("DB_DB_NAME_PREFIX")

	return types.DBConfig{
		URI:             URI,
		Timeout:         Timeout,
		IdleConnTimeout: IdleConnTimeout,
		MaxPoolSize:     MaxPoolSize,
		DBNamePrefix:    DBNamePrefix,
	}
}
