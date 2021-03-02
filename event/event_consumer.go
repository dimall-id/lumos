package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
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
type Callback func(message ConsumerMessage) error

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
func newConsumerConfig (Broker string, ConsumerGroupId string) ConsumerConfig {
	Brokers := strings.Split(Broker, ",")
	return ConsumerConfig{
		Brokers: Brokers,
		ConsumerGroupId: ConsumerGroupId,
	}
}
func newKafkaReadConfig(config ConsumerConfig, topic string) kafka.ReaderConfig {
	c := kafka.ReaderConfig{
		Brokers: config.Brokers,
		GroupID: config.ConsumerGroupId,
		Topic: topic,
	}
	return c
}

func StartConsumer(topic string, callback Callback, config ConsumerConfig) error {
	errs := make(chan error, 1)
	go func() {
		r := kafka.NewReader(newKafkaReadConfig(config, topic))
		defer r.Close()
		for {
			m, err := r.FetchMessage(context.Background())
			if err != nil {
				errs <- err
			}
			message, err := GenerateConsumerMessage(m)
			err = callback(message)
			if err == nil {
				if err := r.CommitMessages(context.Background(), m); err != nil {
					errs <- err
				}
			}
		}
	}()
	return <-errs
}

func StartConsumers (config ConsumerConfig) error {
	if len(callbacks) == 0 {return &NoCallbackRegisteredError{}}

	for topic, callback := range callbacks {
		err := StartConsumer(topic, callback, config)
		if err != nil {
			return err
		}
	}
	return nil
}