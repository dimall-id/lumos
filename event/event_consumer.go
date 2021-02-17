package event

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
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

type Callback func(message kafka.Message) error

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

func StartConsumer(config *kafka.ConfigMap) error {
	if len(callbacks) == 0 {return &NoCallbackRegisteredError{}}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return err
	}

	topics := make([]string, 0)
	for key, _ := range callbacks {
		topics = append(topics, key)
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		return err
	}
	for {
		ev := c.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			topic := *e.TopicPartition.Topic
			if callback, oke := callbacks[topic]; oke {
				err = callback(*e)
				if err == nil {
					_, err = c.CommitMessage(e)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			}
			break
		case *kafka.Error:
			fmt.Printf("%% Error %v\n", e)
		case *kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		}
	}
}