package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUri          string `mapstructure:"MONGODB_LOCAL_URI"`
	DBName         string `mapstructure:"MONGODB_DB_NAME"`
	URLCollection  string `mapstructure:"MONGODB_URL_COLLECTION"`
	UserCollection string `mapstructure:"MONGODB_USER_COLLECTION"`
	APIDefaultPath string `mapstructure:"API_DEFAULT_PATH"`
	APIVersion     string `mapstructure:"API_VERSION"`
	APIShort       string `mapstructure:"API_SHORT"`
	Port           string `mapstructure:"PORT"`
	Host           string `mapstructure:"HOST"`
	Origin         string `mapstructure:"CLIENT_ORIGIN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("dev")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
