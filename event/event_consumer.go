package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
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
			log.Println(err)
			continue
		}
		message, err := GenerateConsumerMessage(m)
		if err != nil {
			log.Println("error in converting data from kafka.message to consumer message")
			log.Println(err)
			continue
		}
		if callback, oke := callbacks[message.Topic]; oke {
			err := callback(message)
			if err != nil {
				log.Println("event callback fail to execute")
				log.Println(err)
				continue
			}
			err = r.CommitMessages(context.Background(), m)
			if err != nil {
				log.Println("fail to commit message")
				log.Println(err)
			}
		} else {
			log.Printf("no callback for topic '%s' found", message.Topic)
		}
	}
	//
	//group, err := kafka.NewConsumerGroup(newKafkaConsumerGroupConfig(config, topics))
	//if err != nil {return err}
	//defer group.Close()
	//
	//for {
	//	gen, err := group.Next(context.TODO())
	//	if err != nil {
	//		return err
	//	}
	//
	//	for topic, callback := range callbacks {
	//		assignments := gen.Assignments[topic]
	//		for _, assignment := range assignments {
	//			partition, offset := assignment.ID, assignment.Offset
	//			gen.Start(func (ctx context.Context) {
	//				reader := kafka.NewReader(newKafkaReadConfig(config, topic, partition))
	//				defer reader.Close()
	//
	//				reader.SetOffset(offset)
	//				for {
	//					msg, err := reader.FetchMessage(context.Background())
	//					switch err {
	//					case kafka.ErrGenerationEnded:
	//						// generation has ended.  commit offsets.  in a real app,
	//						// offsets would be committed periodically.
	//						gen.CommitOffsets(map[string]map[int]int64{"my-topic": {partition: offset + 1}})
	//						return
	//					case nil:
	//						message, err := GenerateConsumerMessage(msg)
	//						err = callback(message)
	//						if err == nil {
	//							offset = msg.Offset
	//						}
	//					default:
	//						fmt.Printf("error reading message: %+v\n", err)
	//					}
	//
	//				}
	//			})
	//		}
	//	}
	//}
}

