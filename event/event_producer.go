package event

import (
	"database/sql"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func initOutboxTable (config Config) error {
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC",
		config.DatasourceConfig.Host,
		config.DatasourceConfig.User,
		config.DatasourceConfig.Password,
		config.DatasourceConfig.Database,
		config.DatasourceConfig.Port)

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
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS outbox (
		  id                  BIGSERIAL PRIMARY KEY,
		  create_time         TIMESTAMP WITH TIME ZONE NOT NULL,
		  kafka_topic         VARCHAR(500) NOT NULL,
		  kafka_key           VARCHAR(500) NOT NULL,  -- pick your own maximum key size
		  kafka_value         VARCHAR(40000),         -- pick your own maximum value size
		  kafka_header_keys   TEXT[] NOT NULL,
		  kafka_header_values TEXT[] NOT NULL,
		  leader_id           UUID
		)`)
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