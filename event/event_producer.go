package event

import (
	"database/sql"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	harvest "github.com/obsidiandynamics/goharvest"
	"github.com/obsidiandynamics/goharvest/stasher"
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
	Sslmode string
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

	kafkaConfig := harvest.KafkaConfigMap{}
	for key, value := range config.KafkaConfig {
		kafkaConfig[key] = value
	}

	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC sslmode=%s",
		config.DatasourceConfig.Host,
		config.DatasourceConfig.User,
		config.DatasourceConfig.Password,
		config.DatasourceConfig.Database,
		config.DatasourceConfig.Port,
		config.DatasourceConfig.Sslmode)

	hConfig := harvest.Config{
		BaseKafkaConfig: kafkaConfig,
		DataSource: connString,
	}

	harvest, err := harvest.New(hConfig)
	if err != nil {
		return err
	}

	err = harvest.Start()
	if err != nil {
		return err
	}

	return nil
}

func SendOutbox (config DatasourceConfig, topic string, key string, value string, headers map[string]string) error {

	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.Database,
		config.Port,
		config.Sslmode)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	defer db.Close()

	st := stasher.New("outbox")

	// Begin a transaction.
	tx, _ := db.Begin()
	defer tx.Rollback()

	// Update other database entities in transaction scope.

	// Stash an outbox record for subsequent harvesting.
	kHeaders := harvest.KafkaHeaders{}
	for key, value := range headers {
		kHeaders = append(kHeaders, harvest.KafkaHeader{
			Key : key,
			Value: value,
		})
	}
	err = st.Stash(tx, harvest.OutboxRecord{
		KafkaTopic: topic,
		KafkaKey:   key,
		KafkaValue: harvest.String(value),
		KafkaHeaders: kHeaders,
	})
	if err != nil {
		return err
	}

	// Commit the transaction.
	tx.Commit()
	return nil
}