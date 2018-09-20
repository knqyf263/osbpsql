package config

import (
	"log"
	"strings"

	"github.com/pivotal-cf/brokerapi"

	"github.com/spf13/viper"
)

const (
	cfgFile              = "config.toml"
	environmentVarPrefix = "osbpsql"
)

var (
	propertyToEnvReplacer = strings.NewReplacer(".", "_", "-", "_")
)

func init() {
	initConfig()
	viper.SetEnvPrefix(environmentVarPrefix)
	viper.SetEnvKeyReplacer(propertyToEnvReplacer)
	viper.AutomaticEnv()
}

func initConfig() {
	if cfgFile == "" {
		return
	}

	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can't read config: %v\n", err)
	}
}

// BrokerConfig is the configurations for ServiceBroker
type BrokerConfig struct {
	Credentials brokerapi.BrokerCredentials
	Port        string
	DB          DBConfig
}

// DBConfig is the configurations for PostgreSQL
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// NewBrokerConfigFromEnv returns BrokerConfig from environment variables
func NewBrokerConfigFromEnv() *BrokerConfig {
	return &BrokerConfig{
		Credentials: brokerapi.BrokerCredentials{
			Username: viper.GetString("basic_auth.username"),
			Password: viper.GetString("basic_auth.password"),
		},
		Port: viper.GetString("broker.port"),
		DB: DBConfig{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
			Database: viper.GetString("db.database"),
		},
	}
}
