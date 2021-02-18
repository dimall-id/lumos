package config

import (
	"github.com/spf13/viper"
	"time"
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

func GetUint (key string) uint {
	return config.GetUint(key)
}

func GetUint32 (key string) uint32 {
	return config.GetUint32(key)
}

func GetUint64 (key string) uint64 {
	return config.GetUint64(key)
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

func GetDuration (key string) time.Duration {
	return config.GetDuration(key)
}

func GetFloat64 (key string) float64 {
	return config.GetFloat64(key)
}