package event

import (
	"database/sql"
	"encoding/json"
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
		TopicPartition: kafka.TopicPartition{
			Topic: &message.KafkaTopic,
			Partition: kafka.PartitionAny,
		},
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

func StartProducer (config Config) error {
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC sslmode=%s",
		config.DatasourceConfig.Host,
		config.DatasourceConfig.User,
		config.DatasourceConfig.Password,
		config.DatasourceConfig.Database,
		config.DatasourceConfig.Port,
		config.DatasourceConfig.SslMode)

	fmt.Println("Starting DB Connection");
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

	fmt.Println("Migrating Outbox Table")
	err = initOutboxTable(db)
	if err != nil {
		return err
	}
	fmt.Println("Done Migrating Outbox Table")

	fmt.Println("Starting Kafka Producer")
	producer, err := kafka.NewProducer(&config.KafkaConfig)
	if err != nil {
		return err
	}
	defer producer.Close()
	fmt.Println("Done Kafka Producer")

	errs := make(chan error)
	if errs != nil {
		return <- errs
	}

	/**
	Reading Kafka Event and update the lumos outbox table to ensure the delivered message as delivered and error message to queue for resend
	 */
	go func() {
		for e := range producer.Events() {
			fmt.Println("Getting New Producer Event")
			switch ev := e.(type) {
			case *kafka.Message:
				var messageId string
				var data map[string]string
				err := json.Unmarshal(ev.Value, &data)
				errs <- err
				messageId = data["id"]
				if ev.TopicPartition.Error != nil {
					db.Model(&LumosOutbox{}).Where("id = ?", messageId).Update("status","QUEUE")
				} else {
					db.Model(&LumosOutbox{}).Where("id = ?", messageId).Updates(LumosOutbox{Status: "DELIVERED", DeliveredAt: time.Now()})
				}
			}
			fmt.Print("Done Processing Producer Event")
		}
	}()

	var messages []LumosOutbox
	for {
		fmt.Printf("Fetching Messaging ...")
		db.Where("status = ?", "QUEUE").Find(&messages)
		fmt.Printf("Processing %d amount of message \n", len(messages))
		if len(messages) > 0 {
			for _, message := range messages {
				kMessage, err := GenerateKafkaMessage(message)
				if err != nil {
					return err
				}
				db.Model(&LumosOutbox{}).Where("id = ?", message.Id).Update("status","DELIVERING")
				err = producer.Produce(&kMessage, nil)
				if err != nil {
					return err
				}
			}
		}

		var PoolDuration time.Duration = 10
		if &config.PoolDuration != nil {
			PoolDuration = config.PoolDuration
		}
		fmt.Printf("Sleep for %d", PoolDuration * time.Second)
		time.Sleep(PoolDuration * time.Second)
		/**
		Clear the data for GC to collect
		 */
		messages = nil
	}
}
