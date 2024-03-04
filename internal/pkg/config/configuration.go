package config

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	DbDriver string `mapstructure:"DBDRIVER"`
	DbName   string `mapstructure:"DBNAME"`
	DbHost   string `mapstructure:"DBHOST"`
	DbPort   string `mapstructure:"DBPORT"`
	DbUser   string `mapstructure:"DBUSER"`
	DbPass   string `mapstructure:"DBPASS"`
	DbSsl    string `mapstructure:"DBSSl"`
}

var Config *Configuration

func LoadConfig(path string) {
	var config *Configuration
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("err read config: %v", err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Printf("err unmarshal config: %v", err)
		return
	}
	Config = config
}

func GetConfig() *Configuration {
	return Config
}
