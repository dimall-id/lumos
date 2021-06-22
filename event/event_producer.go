package event

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Config struct {
	Host string
	DatasourceConfig DatasourceConfig
	PoolDuration time.Duration
}

type DatasourceConfig struct {
	Host string
	Port uint16
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
	query := `
		CREATE OR REPLACE FUNCTION public.new_queue_message()
		 RETURNS trigger
		 LANGUAGE plpgsql
		AS $function$
		declare
			payload jsonb;
		begin
			
			if new.status::varchar = 'QUEUE'::varchar then 
				payload = row_to_json(NEW);
				PERFORM pg_notify('lumos_ouboxes', payload::Text);
				return new;
			end if;
		
		end;
		$function$
		;
		
		CREATE IF NOT EXISTS TABLE public.lumos_outboxes (
			id varchar(50) NOT NULL,
			kafka_topic varchar(255) NOT NULL,
			kafka_key varchar(500) NOT NULL,
			kafka_value varchar(50000) NOT NULL,
			kafka_header_keys varchar(50000) NULL,
			kafka_header_values varchar(50000) NULL,
			created_at timestamp NOT NULL,
			delivered_at timestamp NULL,
			status varchar(100) NOT NULL,
			CONSTRAINT lumos_outboxes_pkey PRIMARY KEY (id)
		);
		
		drop trigger if exists lumos_outbox_inserted on lumos_outboxes;
		create trigger lumos_outbox_inserted
			after insert
			on public.lumos_outboxes
			for each row
			execute procedure public.new_queue_message();
			
		drop trigger if exists lumos_outbox_updated on lumos_outboxes;
		create trigger lumos_outbox_updated
			after update
			on public.lumos_outboxes
			for each row
			execute procedure public.new_queue_message();
	`
	tx := DB.Exec(query)
	return tx.Error
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

func StartProducer (config Config, db *gorm.DB) error {
	log.Info("migrating outbox table")
	err := initOutboxTable(db)
	if err != nil {return err}
	log.Info("done migrating outbox table")

	conn, err := pgx.Connect(pgx.ConnConfig{Host: config.DatasourceConfig.Host, Port: config.DatasourceConfig.Port, User: config.DatasourceConfig.User, Password: config.DatasourceConfig.Password, Database: config.DatasourceConfig.Database})
	if err != nil {
		log.Errorf("Fail to open database connection due to '%s'", err)
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Errorf("fail to close database donnection due to '%s'", err)
		}
	}()

	err = conn.Listen("lumos_outbox")
	if err != nil {
		log.Errorf("fail to listen to lumos_outbox notify due to '%s'", err)
		return err
	}

	for {
		msg, err := conn.WaitForNotification(context.Background())
		if err != nil {
			log.Errorf("fail to get notification message due to '%s'", err)
			return err
		}
		go func() {
			message := LumosOutbox{}
			err = json.Unmarshal([]byte(msg.Payload), &message)
			if err != nil {
				log.Errorf("fail to decode message due to '%s'", err)
				return
			}

			kMessage, err := GenerateKafkaMessage(message)
			if err != nil {
				log.Errorf("fail to generate kafka message due to '%s'", err)
				return
			}
			db.Model(&LumosOutbox{}).Where("id = ?", message.Id).Update("status","DELIVERING")
			err = SendMessage(message.KafkaTopic, config, kMessage)
			if err != nil {
				log.Errorf("put back message to QUEUE due to %s", err.Error())
				db.Model(&LumosOutbox{}).Where("id = ?", message.Id).Update("status", "QUEUE")
				log.Info("message put backed to QUEUE")
			} else {
				log.Info("marking message as delivered")
				db.Model(&LumosOutbox{}).Where("id = ?", message.Id).Updates(LumosOutbox{Status: "DELIVERED", DeliveredAt: time.Now()})
				log.Info("message marked as delivered")
			}
		}()
	}
}
