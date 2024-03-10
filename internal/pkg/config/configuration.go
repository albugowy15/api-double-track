package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	AppEnv      string `mapstructure:"APP_ENV"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
	Port        string `mapstructure:"PORT"`
	Secret      string `mapstructure:"SECRET"`
}

var config *Configuration

func LoadConfig(path string) {
	var conf *Configuration
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("err read config: %v", err)
		return
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Printf("err unmarshal config: %v", err)
		return
	}
	config = conf
}

func GetConfig() *Configuration {
	return config
}
