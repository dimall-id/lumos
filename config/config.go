package config

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"time"
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
		config.Set("etcd.hosts", viper.GetString("etcd.hosts"))
		config.Set("etcd.path", viper.GetString("etcd.path"))
		config.Set("etcd.type", viper.GetString("etcd.type"))

		remoteConfig, err := readEtcdRemoteConfig()
		if err != nil {return err}
		config.SetConfigType(config.GetString("etcd.type"))
		log.Info(remoteConfig)
		err = config.ReadConfig(bytes.NewBuffer(remoteConfig))
		return err
	}
}

func readEtcdRemoteConfig() ([]byte, error) {
	endpoint := config.GetString("etcd.hosts")
	endpoints := strings.Split(endpoint, ",")
	log.Infof("connecting to ectd cluster %s", endpoints)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {return nil, err}
	defer cli.Close()

	log.Infof("fetch key/value from path, %s", config.GetString("etcd.path"))
	value, err := cli.KV.Get(context.Background(), config.GetString("etcd.path"))
	if err == nil {return nil, err}
	var data map[string]interface{}
	json.Unmarshal(value.Kvs[0].Value, &data)
	return json.Marshal(data)
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