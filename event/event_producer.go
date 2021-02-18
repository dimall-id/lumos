package event

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Config struct {
	KafkaConfig kafka.ConfigMap
	DatasourceConfig DatasourceConfig
	PoolDuration time.Duration
}

type DatasourceConfig struct {
	Host string
	Port string
	User string
	Password string
	Database string
	SslMode string
}

type LumosOutbox struct {
	Id string `gorm:"id,primaryKey"`
	KafkaTopic string `gorm:"kafka_topic,size:500"`
	KafkaKey string `gorm:"kafka_key,size:500"`
	KafkaValue string `gorm:"kafka_value,size:50000"`
	KafkaHeaderKeys string `gorm:"kafka_header_keys,size:50000"`
	KafkaHeaderValues string `gorm:"kafka_header_values,size:50000"`
	CreatedAt time.Time `gorm:"created_at,notNull"`
	DeliveredAt time.Time`gorm:"delivered_at"`
	Status string `gorm:"status,100"`
}

func initOutboxTable (DB *gorm.DB) error {
	err := DB.AutoMigrate(&LumosOutbox{})
	return err
}

func StartProducer (config Config) error {
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC sslmode=%s",
		config.DatasourceConfig.Host,
		config.DatasourceConfig.User,
		config.DatasourceConfig.Password,
		config.DatasourceConfig.Database,
		config.DatasourceConfig.Port,
		config.DatasourceConfig.SslMode)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return err
	}

	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)

	defer sqlDB.Close()

	err = initOutboxTable(db)
	if err != nil {
		return err
	}

	producer, err := kafka.NewProducer(&config.KafkaConfig)
	if err != nil {
		return err
	}
	defer producer.Close()

	/**
	Reading Kafka Event and update the lumos outbox table to ensure the delivered message as delivered and error message to queue for resend
	 */
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				var messageId string
				for _,header := range ev.Headers {
					if header.Key == "MESSAGE-ID" {
						messageId = string(header.Value)
						break
					}
				}
				if ev.TopicPartition.Error != nil {
					db.Model(&LumosOutbox{}).Where("id = ?", messageId).Update("status","QUEUE")
				} else {
					db.Model(&LumosOutbox{}).Where("id = ?", messageId).Updates(LumosOutbox{Status: "DELIVERED", DeliveredAt: time.Now()})
				}
			}
		}
	}()

	for {
		var messages []LumosOutbox
		db.Where("status = ?", "QUEUE").Find(&messages)
		if len(messages) > 0 {
			for _, message := range messages {
				var keys = strings.Split(message.KafkaHeaderKeys, ",")
				var values = strings.Split(message.KafkaHeaderValues, ",")
				var headers = make([]kafka.Header, len(keys))
				for idx , key := range keys {
					if idx < len(values) {
						headers = append(headers, kafka.Header{Key: key, Value: []byte(values[idx])})
					}
				}
				headers = append(headers, kafka.Header{
					Key: "MESSAGE-ID",
					Value: []byte(message.Id),
				})
				db.Model(&LumosOutbox{}).Where("id = ?", message.Id).Update("status","DELIVERING")
				producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &message.KafkaTopic, Partition: kafka.PartitionAny},
					Value: []byte(message.KafkaValue),
					Key: []byte(message.KafkaKey),
					Headers: headers,
				}, nil)
			}
		}

		var PoolDuration time.Duration = 60
		if &config.PoolDuration != nil {
			PoolDuration = config.PoolDuration
		}
		time.Sleep(PoolDuration)
		/**
		Clear the data for GC to collect
		 */
		messages = nil
	}
}

type LumosMessage struct {
	Topic string
	Key string
	Value string
	Headers map[string]string
}

func GenerateOutbox (DB *gorm.DB, message LumosMessage) error {

	var keys = make([]string, len(message.Headers))
	var values = make([]string, len(message.Headers))
	var i = 0
	for k, v := range message.Headers {
		keys[i] = k
		values[i] = v
		i ++
	}

	data := LumosOutbox{
		Id: uuid.New().String(),
		KafkaTopic: message.Topic,
		KafkaKey: message.Key,
		KafkaValue: message.Value,
		KafkaHeaderKeys: strings.Join(keys[:], ","),
		KafkaHeaderValues: strings.Join(values[:], ","),
		CreatedAt: time.Now(),
		Status: "QUEUE",
	}

	tx := DB.Save(data)

	return tx.Error
}