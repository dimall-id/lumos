package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Config struct {
	Host             string
	DatasourceConfig DatasourceConfig
	PoolDuration     time.Duration
}

type DatasourceConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Database string
	SslMode  string
}

type LumosOutbox struct {
	Id                string `json:"id" gorm:"id,primaryKey;type:varchar;size:50"`
	KafkaTopic        string `json:"kafka_topic" gorm:"kafka_topic;type:varchar;size:500;notNull"`
	KafkaKey          string `json:"kafka_key" gorm:"kafka_key;type:varchar;size:500;notNull"`
	KafkaValue        string `json:"kafka_value" gorm:"kafka_value;type:varchar;size:50000;notNull"`
	KafkaHeaderKeys   string `json:"kafka_header_keys" gorm:"kafka_header_keys;type:varchar;size:50000"`
	KafkaHeaderValues string `json:"kafka_header_values" gorm:"kafka_header_values;type:varchar;size:50000"`
	CreatedAt         int64  `json:"created_at" gorm:"created_at;notNull"`
	DeliveredAt       int64  `json:"delivered_at" gorm:"delivered_at"`
	Status            string `json:"status" gorm:"status,type:varchar;size:100;index:status_index;notNull"`
}

type LumosMessage struct {
	Topic   string
	Key     string
	Value   string
	Headers map[string]string
}

func initOutboxTable(DB *pgx.ConnPool) error {
	query := `
		CREATE OR REPLACE FUNCTION public.new_queue_message()
		 RETURNS trigger
		 LANGUAGE plpgsql
		AS $function$
			declare
				payload varchar;
			begin
				
				if new.status::varchar = 'QUEUE'::varchar then 
					payload = new.id::varchar;
					PERFORM pg_notify('lumos_ouboxes', payload::Text);
					return new;
				end if;
				return new;
			
			end; 
		$function$;
		
		CREATE table if not exists public.lumos_outboxes (
			id varchar(50) NOT NULL,
			kafka_topic varchar(255) NOT NULL,
			kafka_key varchar(500) NOT NULL,
			kafka_value text NOT NULL,
			kafka_header_keys varchar(50000) NULL,
			kafka_header_values text NULL,
			created_at int8 NOT NULL,
			delivered_at int8 NULL,
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
	_, err := DB.Exec(query)
	return err
}

func GenerateKafkaMessage(message LumosOutbox) (kafka.Message, error) {
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
	for idx, key := range keys {
		if idx < len(values) {
			headers = append(headers, kafka.Header{Key: key, Value: []byte(values[idx])})
		}
	}

	data := map[string]string{
		"id":   message.Id,
		"data": message.KafkaValue,
	}
	j, e := json.Marshal(data)
	if e != nil {
		return kafka.Message{}, e
	}
	return kafka.Message{
		Headers: headers,
		Key:     []byte(message.KafkaKey),
		Value:   j,
	}, nil
}

func GenerateOutbox(DB *gorm.DB, message LumosMessage) error {
	var keys = make([]string, len(message.Headers))
	var values = make([]string, len(message.Headers))
	var i = 0
	for k, v := range message.Headers {
		keys[i] = k
		values[i] = v
		i++
	}

	data := LumosOutbox{
		Id:                uuid.New().String(),
		KafkaTopic:        message.Topic,
		KafkaKey:          message.Key,
		KafkaValue:        message.Value,
		KafkaHeaderKeys:   strings.Join(keys[:], ","),
		KafkaHeaderValues: strings.Join(values[:], ","),
		CreatedAt:         time.Now().Unix(),
		Status:            "QUEUE",
	}

	tx := DB.Save(data)

	return tx.Error
}

func SendMessage(topic string, config Config, message kafka.Message) error {
	w := &kafka.Writer{
		Addr:     kafka.TCP(config.Host),
		Topic:    topic,
		Balancer: kafka.CRC32Balancer{},
	}

	err := w.WriteMessages(context.Background(), message)

	return err
}

func StartProducer(config Config) error {
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     config.DatasourceConfig.Host,
			Port:     config.DatasourceConfig.Port,
			User:     config.DatasourceConfig.User,
			Password: config.DatasourceConfig.Password,
			Database: config.DatasourceConfig.Database,
		},
		MaxConnections: 100,
	})
	if err != nil {
		log.Errorf("fail to open connection pool due to '%s", err)
		return err
	}
	defer conn.Close()

	log.Warn("migrating outbox table")
	err = initOutboxTable(conn)
	if err != nil {
		log.Errorf("fail to init outbox table due to '%s'", err)
		return err
	}
	log.Warn("success migrating outbox table")

	lconn, err := pgx.Connect(pgx.ConnConfig{
		Host:     config.DatasourceConfig.Host,
		Port:     config.DatasourceConfig.Port,
		User:     config.DatasourceConfig.User,
		Password: config.DatasourceConfig.Password,
		Database: config.DatasourceConfig.Database,
	})
	if err != nil {
		log.Errorf("fail to open listener connection due to '%s'", err)
		return err
	}

	err = lconn.Listen("lumos_ouboxes")
	if err != nil {
		log.Errorf("fail to listen to lumos_ouboxes channel due to '%s'", err)
		return err
	}

	for {
		msg, err := lconn.WaitForNotification(context.Background())
		if err != nil {
			log.Errorf("fail to get notification message due to '%s'", err)
			return err
		}
		go func() {
			id := msg.Payload
			var message LumosOutbox
			rows, err := conn.Query(fmt.Sprintf("SELECT * FROM public.lumos_outboxes WHERE id = '%s'", id))
			if err != nil {
				log.Error("fail to fetch lumos_outboxes data")
				log.Error(err)
			}
			err = rows.Scan(message)
			if err != nil {
				log.Error("fail to scan message data")
				log.Error(err)
			}
			kMessage, err := GenerateKafkaMessage(message)
			if err != nil {
				log.Errorf("fail to generate kafka message due to '%s'", err)
				return
			}
			conn.Exec(fmt.Sprintf("UPDATE public.lumos_outboxes SET status='DELIVERING' where id = '%s'", message.Id))
			err = SendMessage(message.KafkaTopic, config, kMessage)
			if err != nil {
				log.Errorf("put back message to QUEUE due to %s", err.Error())
				_, err = conn.Exec(fmt.Sprintf("UPDATE public.lumos_outboxes SET status='QUEUE' where id = '%s'", message.Id))
				if err != nil {
					log.Errorf("fail to update message to QUEUE due to '%s'", err)
				} else {
					log.Warn("message put backed to QUEUE")
				}
			} else {
				log.Warn("marking message as delivered")
				_, err = conn.Exec(fmt.Sprintf("UPDATE public.lumos_outboxes SET status='DELIVERED', delivered_at=%d where id = '%s'", time.Now().Unix(), message.Id))
				if err != nil {
					log.Errorf("fail to update message to DELIVERED due to '%s'", err)
				} else {
					log.Warn("message marked as delivered")
				}
			}
		}()
	}
}
