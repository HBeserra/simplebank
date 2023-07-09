package util

import (
	"github.com/spf13/viper"
	"log"
)

// Config stores all configurations of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBUri         string `mapstructure:"DB_URI"`
	DBMaxNConn    int    `mapstructure:"DB_MAX_NUMBER_CONNECTIONS"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (conf Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") //ENV, JSON, XML

	viper.SetDefault("DB_MAX_NUMBER_CONNECTIONS", 5) // Define the default max number of connections to db

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Panicln("erro can't load configuration:", err.Error())
		return conf, err
	}

	err = viper.Unmarshal(&conf)
	return conf, err
}
