package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/coneno/logger"
	"github.com/tekenradar/content-service/pkg/types"
)

const (
	ENV_INSTANCE_IDS     = "INSTANCE_IDS"
	ENV_EMAIL_CLIENT_URL = "EMAIL_CLIENT_URL"
	ENV_FORCE_REPLACE    = "FORCE_REPLACE"
	ENV_CSV_PATH         = "CSV_PATH"
	ENV_SEPARATOR        = "SEPARATOR"

	ENV_CONTENT_DB_CONNECTION_STR    = "CONTENT_DB_CONNECTION_STR"
	ENV_CONTENT_DB_USERNAME          = "CONTENT_DB_USERNAME"
	ENV_CONTENT_DB_PASSWORD          = "CONTENT_DB_PASSWORD"
	ENV_CONTENT_DB_CONNECTION_PREFIX = "CONTENT_DB_CONNECTION_PREFIX"
	ENV_DB_TIMEOUT                   = "DB_TIMEOUT"
	ENV_DB_IDLE_CONN_TIMEOUT         = "DB_IDLE_CONN_TIMEOUT"
	ENV_DB_MAX_POOL_SIZE             = "DB_MAX_POOL_SIZE"
	ENV_DB_NAME_PREFIX               = "DB_DB_NAME_PREFIX"

	ENV_INVITATION_EMAIL_TEMPLATE_PATH = "INVITATION_EMAIL_TEMPLATE_PATH"
	ENV_INVITATION_EMAIL_SUBJECT       = "INVITATION_EMAIL_SUBJECT"

	ENV_RUN_PARTICIPANT_CREATION = "RUN_PARTICIPANT_CREATION"
	ENV_RUN_INVITATION_SENDING   = "RUN_INVITATION_SENDING"

	defaultCSVFile   = "participants.csv"
	defaultSeparator = ";"
)

type Config struct {
	InstanceIDs                 []string
	ContentDBConfig             types.DBConfig
	EmailClientURL              string
	ForceReplace                bool
	CSVPath                     string
	Separator                   string
	InvitationEmailTemplatePath string
	InvitationEmailSubject      string
	RunParticipantCreation      bool
	RunInvitationSending        bool
}

func readConfig() Config {
	conf := Config{}
	conf.EmailClientURL = os.Getenv(ENV_EMAIL_CLIENT_URL)
	conf.ForceReplace = os.Getenv(ENV_FORCE_REPLACE) == "true"
	conf.CSVPath = os.Getenv(ENV_CSV_PATH)
	if conf.CSVPath == "" {
		conf.CSVPath = defaultCSVFile
	}
	conf.Separator = os.Getenv(ENV_SEPARATOR)
	if conf.Separator == "" {
		conf.Separator = defaultSeparator
	}
	conf.InstanceIDs = strings.Split(os.Getenv(ENV_INSTANCE_IDS), ",")
	conf.ContentDBConfig = getContentDBConfig()

	conf.InvitationEmailTemplatePath = os.Getenv(ENV_INVITATION_EMAIL_TEMPLATE_PATH)
	conf.InvitationEmailSubject = os.Getenv(ENV_INVITATION_EMAIL_SUBJECT)

	conf.RunParticipantCreation = os.Getenv(ENV_RUN_PARTICIPANT_CREATION) == "true"
	conf.RunInvitationSending = os.Getenv(ENV_RUN_INVITATION_SENDING) == "true"

	return conf
}

func getContentDBConfig() types.DBConfig {
	connStr := os.Getenv(ENV_CONTENT_DB_CONNECTION_STR)
	username := os.Getenv(ENV_CONTENT_DB_USERNAME)
	password := os.Getenv(ENV_CONTENT_DB_PASSWORD)
	prefix := os.Getenv(ENV_CONTENT_DB_CONNECTION_PREFIX) // Used in test mode
	if connStr == "" || username == "" || password == "" {
		logger.Error.Fatal("Couldn't read DB credentials.")
	}
	URI := fmt.Sprintf(`mongodb%s://%s:%s@%s`, prefix, username, password, connStr)

	var err error
	Timeout, err := strconv.Atoi(os.Getenv(ENV_DB_TIMEOUT))
	if err != nil {
		logger.Error.Fatal("DB_TIMEOUT: " + err.Error())
	}
	IdleConnTimeout, err := strconv.Atoi(os.Getenv(ENV_DB_IDLE_CONN_TIMEOUT))
	if err != nil {
		logger.Error.Fatal("DB_IDLE_CONN_TIMEOUT" + err.Error())
	}
	mps, err := strconv.Atoi(os.Getenv(ENV_DB_MAX_POOL_SIZE))
	MaxPoolSize := uint64(mps)
	if err != nil {
		logger.Error.Fatal("DB_MAX_POOL_SIZE: " + err.Error())
	}

	DBNamePrefix := os.Getenv(ENV_DB_NAME_PREFIX)

	return types.DBConfig{
		URI:             URI,
		Timeout:         Timeout,
		IdleConnTimeout: IdleConnTimeout,
		MaxPoolSize:     MaxPoolSize,
		DBNamePrefix:    DBNamePrefix,
	}
}
