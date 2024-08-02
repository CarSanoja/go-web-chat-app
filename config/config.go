package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string
	ServerPort    string
}

var config Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	config = Config{
		ServerAddress: viper.GetString("server_address"),
		ServerPort:    viper.GetString("server_port"),
	}
}

func GetConfig() Config {
	return config
}
