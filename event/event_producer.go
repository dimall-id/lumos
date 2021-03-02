package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	Host string
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
	Id string `gorm:"id,primaryKey;type:varchar;size:50"`
	KafkaTopic string `gorm:"kafka_topic;type:varchar;size:500;notNull"`
	KafkaKey string `gorm:"kafka_key;type:varchar;size:500;notNull"`
	KafkaValue string `gorm:"kafka_value;type:varchar;size:50000;notNull"`
	KafkaHeaderKeys string `gorm:"kafka_header_keys;type:varchar;size:50000"`
	KafkaHeaderValues string `gorm:"kafka_header_values;type:varchar;size:50000"`
	CreatedAt time.Time `gorm:"created_at;notNull"`
	DeliveredAt time.Time`gorm:"delivered_at"`
	Status string `gorm:"status,type:varchar;size:100;index:status_index;notNull"`
}

type LumosMessage struct {
	Topic string
	Key string
	Value string
	Headers map[string]string
}

func initOutboxTable (DB *gorm.DB) error {
	err := DB.AutoMigrate(&LumosOutbox{})
	return err
}

func GenerateKafkaMessage (message LumosOutbox) (kafka.Message,error) {
	var headers []kafka.Header
	var keys []string
	var values []string
	if message.KafkaHeaderKeys != "" {
		keys = strings.Split(message.KafkaHeaderKeys, ",")
		values = strings.Split(message.KafkaHeaderValues, ",")
		headers = make([]kafka.Header, len(keys))
	} else {
		keys = make([]string, 0)
		values = make([]string, 0)
		headers = make([]kafka.Header, 0)
	}
	for idx , key := range keys {
		if idx < len(values) {
			headers = append(headers, kafka.Header{Key: key, Value: []byte(values[idx])})
		}
	}

	data := map[string]string {
		"id" : message.Id,
		"data" : message.KafkaValue,
	}
	j,e := json.Marshal(data)
	if e != nil {
		return kafka.Message{}, e
	}
	return kafka.Message{
		Headers: headers,
		Key: []byte(message.KafkaKey),
		Value: j,
	}, nil
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

func SendMessage (topic string, config Config, message kafka.Message) error {
	w := &kafka.Writer{
		Addr : kafka.TCP(config.Host),
		Topic: topic,
		Balancer: kafka.CRC32Balancer{},
	}

	err := w.WriteMessages(context.Background(), message)

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

	fmt.Printf("[%s] Starting DB Connection\n", time.Now().Format(time.RFC850))
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		PrepareStmt: true,
	})
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

	fmt.Printf("[%s] Migrating Outbox Table\n", time.Now().Format(time.RFC850))
	err = initOutboxTable(db)
	if err != nil {
		return err
	}
	fmt.Printf("[%s] Done Migrating Outbox Table\n", time.Now().Format(time.RFC850))

	var messages []LumosOutbox
	for {
		fmt.Printf("[%s] Fetching Messaging ... \n", time.Now().Format(time.RFC850))
		db.Where("status = ?", "QUEUE").Find(&messages)
		fmt.Printf("[%s] Processing %d amount of message \n", time.Now().Format(time.RFC850), len(messages))
		if len(messages) > 0 {
			for _, message := range messages {
				kMessage, err := GenerateKafkaMessage(message)
				if err != nil {
					return err
				}
				db.Model(&LumosOutbox{}).Where("id = ?", message.Id).Update("status","DELIVERING")
				err = SendMessage(message.KafkaTopic, config, kMessage)
				if err != nil {
					fmt.Printf("[%s] Put Backn Message to QUEUE \n", time.Now().Format(time.RFC850))
					db.Model(&LumosOutbox{}).Where("id = ?", message.Id).Update("status", "QUEUE")
					return err
				} else {
					fmt.Printf("[%s] Mark Message as Delivered \n", time.Now().Format(time.RFC850))
					db.Model(&LumosOutbox{}).Where("id = ?", message.Id).Updates(LumosOutbox{Status: "DELIVERED", DeliveredAt: time.Now()})
				}
			}
		}

		var PoolDuration time.Duration = 10
		if &config.PoolDuration != nil {
			PoolDuration = config.PoolDuration
		}
		fmt.Printf("[%s] Sleep for %d Seconds\n", time.Now().Format(time.RFC850), PoolDuration)
		time.Sleep(PoolDuration * time.Second)
		/**
		Clear the data for GC to collect
		 */
		messages = nil
	}
}
