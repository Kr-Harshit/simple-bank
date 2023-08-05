package util

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The value are read by viper from a config file or enviroment variables
type Config struct {
	Database DBConfig
	App      APPConfig
	Token    TokenConfig
}

type DBConfig struct {
	Name          string
	Source        string
	MigrateSource string `mapstructure:"migrate_source"`
}

type APPConfig struct {
	Address string
	Port    string
}

type TokenConfig struct {
	Key      string
	Duration time.Duration
}

func LoadConfig(cfgFile string) (config Config, err error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// setting default config path
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("ENV")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindEnv("token.key")
	viper.BindEnv("database.source")
	viper.BindEnv("database.MigrateSource")

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	log.Printf("viper settings: %v", viper.AllSettings())

	err = viper.Unmarshal(&config)
	log.Printf("config: %+v", config)
	return
}
