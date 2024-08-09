package main

import (
	"fmt"
	"os"

	"github.com/coneno/logger"
	"github.com/tekenradar/content-service/pkg/types"
	"gopkg.in/yaml.v2"
)

const (
	ENV_CONFIG_PATH         = "CONFIG_PATH"
	ENV_CONTENT_DB_USERNAME = "CONTENT_DB_USERNAME"
	ENV_CONTENT_DB_PASSWORD = "CONTENT_DB_PASSWORD"
)

type Config struct {
	InstanceIDs     []string     `yaml:"instanceIDs"`
	ContentDBConfig DBConfigYaml `yaml:"contentDBConfig"`
	EmailClientURL  string       `yaml:"emailClientURL"`

	RunTasks struct {
		ParticipantCreation bool `yaml:"participantCreation"`
		InvitationSending   bool `yaml:"invitationSending"`
		ReminderSending     bool `yaml:"reminderSending"`
	} `yaml:"runTasks"`

	ParticipantCreation struct {
		CSVPath      string `yaml:"csvPath"`
		Separator    string `yaml:"separator"`
		ForceReplace bool   `yaml:"forceReplace"`
	} `yaml:"participantCreation"`

	WebsiteURL     string `yaml:"websiteURL"`
	MessageConfigs map[string]struct {
		InvitationTemplatePath string `yaml:"invitationTemplatePath"`
		InvitationSubject      string `yaml:"invitationSubject"`
		ReminderTemplatePath   string `yaml:"reminderTemplatePath"`
		ReminderSubject        string `yaml:"reminderSubject"`
	} `yaml:"messageConfigs"`
}

type DBConfigYaml struct {
	ConnectionStr      string `yaml:"connection_str"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	ConnectionPrefix   string `yaml:"connection_prefix"`
	Timeout            int    `yaml:"timeout"`
	IdleConnTimeout    int    `yaml:"idle_conn_timeout"`
	MaxPoolSize        int    `yaml:"max_pool_size"`
	UseNoCursorTimeout bool   `yaml:"use_no_cursor_timeout"`
	DBNamePrefix       string `yaml:"db_name_prefix"`
	RunIndexCreation   bool   `yaml:"run_index_creation"`
}

func init() {
	// Read config from file
	yamlFile, err := os.ReadFile(os.Getenv(ENV_CONFIG_PATH))
	if err != nil {
		panic(err)
	}

	err = yaml.UnmarshalStrict(yamlFile, &conf)
	if err != nil {
		panic(err)
	}

	secretsOverride()
}

func getContentDBConfig() types.DBConfig {
	connStr := conf.ContentDBConfig.ConnectionStr
	username := conf.ContentDBConfig.Username
	password := conf.ContentDBConfig.Password
	prefix := conf.ContentDBConfig.ConnectionPrefix // Used in test mode
	if connStr == "" || username == "" || password == "" {
		logger.Error.Fatal("Couldn't read DB credentials.")
	}
	URI := fmt.Sprintf(`mongodb%s://%s:%s@%s`, prefix, username, password, connStr)

	return types.DBConfig{
		URI:             URI,
		Timeout:         conf.ContentDBConfig.Timeout,
		IdleConnTimeout: conf.ContentDBConfig.IdleConnTimeout,
		MaxPoolSize:     uint64(conf.ContentDBConfig.MaxPoolSize),
		DBNamePrefix:    conf.ContentDBConfig.DBNamePrefix,
	}
}

func secretsOverride() {
	// Override secrets from environment variables

	if dbUsername := os.Getenv(ENV_CONTENT_DB_USERNAME); dbUsername != "" {
		conf.ContentDBConfig.Username = dbUsername
	}

	if dbPassword := os.Getenv(ENV_CONTENT_DB_PASSWORD); dbPassword != "" {
		conf.ContentDBConfig.Password = dbPassword
	}
}
