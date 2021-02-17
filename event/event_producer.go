package event

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Config struct {
	KafkaConfig kafka.ConfigMap
	DatasourceConfig DatasourceConfig
}

type DatasourceConfig struct {
	Host string
	Port string
	User string
	Password string
	Database string
}

type Outbox struct {
	Id int `gorm:"id,type:BIGSERIAL"`
	createTime time.Time `gorm:"create_time"`
	KafkaTopic string `gorm:"kafka_topic,size:255,notNull"`
	KafkaKey string `gorm:"kafka_topic,size:255,notNull"`
	KafkaValue string `gorm:"kafka_topic,size:100000"`
	KafkaHeaderKeys []string `gorm:"kafka_header_keys,type:TEXT[],notNull"`
	KafkaHeaderValues []string`gorm:"kafka_header_values,type:TEXT[],notNull"`
	LeaderId string `gorm:"leader_id,type:UUID"`
}

func initOutboxTable (config Config) error {
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC",
		viper.GetString(config.DatasourceConfig.Host),
		viper.GetString(config.DatasourceConfig.User),
		viper.GetString(config.DatasourceConfig.Password),
		viper.GetString(config.DatasourceConfig.Database),
		viper.GetString(config.DatasourceConfig.Port))

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return err
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	err = db.AutoMigrate(Outbox{})
	if err != nil {
		return err
	}
	return nil
}

func StartProducer (config Config) error {
	err := initOutboxTable(config)
	if err != nil {
		return nil
	}
	return nil
}