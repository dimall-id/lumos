package config

import (
	"github.com/spf13/viper"
)

var config = viper.New()

func InitConfig (filename string, path string) error {
	config.SetConfigName(filename)
	config.AddConfigPath(path)
	err := config.ReadInConfig()
	return err
}

func WatchConfig() {
	config.WatchConfig()
}

func Get(key string) interface{} {
	return config.Get(key)
}

func GetBool (key string) bool {
	return config.GetBool(key)
}

func GetInt (key string) int {
	return config.GetInt(key)
}

func GetInt32 (key string) int32{
	return config.GetInt32(key)
}

func GetInt64 (key string) int64 {
	return config.GetInt64(key)
}

func GetIntSlice (key string) []int {
	return config.GetIntSlice(key)
}

func GetString (key string) string {
	return config.GetString(key)
}

func GetStringSlice (key string) []string {
	return config.GetStringSlice(key)
}

func GetStringMap (key string) map[string]interface{} {
	return config.GetStringMap(key)
}

func GetStringMapString (key string) map[string]string {
	return config.GetStringMapString(key)
}

func GetStringMapStringSlice (key string) map[string][]string {
	return config.GetStringMapStringSlice(key)
}