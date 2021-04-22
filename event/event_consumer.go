package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ExistingCallbackError struct {
	topic string
}
func (ex *ExistingCallbackError) Error() string {
	return fmt.Sprintf("Existing callback for topic %s is founded", ex.topic)
}

type NoExistingCallbackError struct {
	topic string
}
func (ex *NoExistingCallbackError) Error() string {
	return fmt.Sprintf("No Existing callback for topic %s is founded", ex.topic)
}

type NoCallbackRegisteredError struct {}
func (n *NoCallbackRegisteredError) Error() string {
	return "No callback registered"
}

type ConsumerMessage struct {
	Topic string
	MessageId string
	Headers []kafka.Header
	Value string
}
func GenerateConsumerMessage (message kafka.Message) (ConsumerMessage, error) {
	var value map[string]string
	err := json.Unmarshal(message.Value, &value)
	if err != nil {
		return ConsumerMessage{}, nil
	}

	return ConsumerMessage{
		Topic: message.Topic,
		MessageId: value["id"],
		Headers: message.Headers,
		Value: value["data"],
	}, nil
}
type Callback func(message ConsumerMessage, logger *log.Entry) error

var callbacks = make(map[string]Callback)
func AddCallback (topic string, callback Callback) error {
	if _, oke := callbacks[topic]; oke {
		return &ExistingCallbackError{topic}
	} else {
		callbacks[topic] = callback
	}
	return nil
}
func RemoveCallback (topic string) error {
	if _, oke := callbacks[topic]; oke {
		delete(callbacks, topic)
		return nil
	}
	return &NoExistingCallbackError{topic}
}

type ConsumerConfig struct {
	ConsumerGroupId  string
	Brokers []string
}
func NewConsumerConfig (Broker string, ConsumerGroupId string) ConsumerConfig {
	Brokers := strings.Split(Broker, ",")
	return ConsumerConfig{
		Brokers: Brokers,
		ConsumerGroupId: ConsumerGroupId,
	}
}

func StartConsumers (config ConsumerConfig) error {
	if len(callbacks) == 0 {return &NoCallbackRegisteredError{}}

	topics := make([]string, len(callbacks))
	i := 0
	for topic, _ := range callbacks {
		topics[i] = topic
		i ++
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Brokers,
		GroupID: config.ConsumerGroupId,
		GroupTopics: topics,
	})

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Errorf("fail to fetch message due to %s", err.Error())
			continue
		}
		message, err := GenerateConsumerMessage(m)
		if err != nil {
			log.Errorf("error in converting data from kafka.message to ConsumerMessage %s", err.Error())
			continue
		}
		logger := log.WithFields(log.Fields{
			"topic" : message.Topic,
			"partition" : m.Partition,
			"offset" : m.Offset,
			"message_id" : message.MessageId,
		})
		if callback, oke := callbacks[message.Topic]; oke {
			err := callback(message, logger)
			if err != nil {
				logger.Errorf("event callback fail to finish it job due to %s", err.Error())
				continue
			}
			err = r.CommitMessages(context.Background(), m)
			if err != nil {
				logger.Errorf("fail to commit message due to %s", err.Error())
			}
		} else {
			log.Errorf("no callback for topic '%s' found", message.Topic)
		}
	}
}

