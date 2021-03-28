package config

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"strings"
	"time"
	log "github.com/sirupsen/logrus"
)

var config = viper.New()

func InitConfig (env string) error {
	if strings.ToUpper(env) == "DEBUG" {
		config.SetConfigName("config")
		config.AddConfigPath(".")
		err := config.ReadInConfig()
		return err
	} else {
		viper.SetEnvKeyReplacer(strings.NewReplacer(".","_"))
		viper.SetEnvPrefix("DIMALL")
		viper.SetDefault("etcd.host", "http://localhost:2379")
		viper.AutomaticEnv()

		config.Set("service.name", viper.GetString("service.name"))
		config.Set("db.host", viper.GetString("db.host"))
		config.Set("db.username", viper.GetString("db.username"))
		config.Set("db.password", viper.GetString("db.password"))
		config.Set("db.database", viper.GetString("db.database"))
		config.Set("etcd.host", viper.GetString("etcd.host"))
		config.Set("etcd.path", viper.GetString("etcd.path"))
		config.Set("etcd.type", viper.GetString("etcd.type"))

		log.Infoln("ETCD HOST", config.GetString("etcd.host"))
		log.Infoln("ETCD PATH", config.GetString("etcd.path"))
		err := config.AddRemoteProvider("etcd", config.GetString("etcd.host"), config.GetString("etcd.path"))
		if err != nil {
			log.Errorln(err)
			return err
		}
		config.SetConfigFile(config.GetString("etcd.type"))
		err = config.ReadRemoteConfig()
		return err
	}
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