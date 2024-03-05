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
	Port     string `mapstructure:"PORT"`
	Secret   string `mapstructure:"SECRET"`
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
