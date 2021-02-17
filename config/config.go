package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

var config = viper.New()

func InitConfig () {
	config.SetConfigName("config")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	if strings.ToUpper(config.GetString("service.env")) == "PROD" && config.InConfig("etcd.") {

	}
}

func GetString (key string) string {
	return config.GetString(key)
}