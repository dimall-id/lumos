package config

import (
	"github.com/spf13/viper"
	"log"
)

var config = viper.New()

func InitConfig () {
	config.SetConfigName("config")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func GetString (key string) string {
	return config.GetString(key)
}