package nsqPublisher

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Somesh/go-boilerplate/common/config"
	"github.com/nsqio/go-nsq"

	"gopkg.in/tokopedia/logging.v1"
)

var producer *NSQPublisher

func init() {
	cfg := config.GetConfig()
	publisher, err := NewNSQPublisher(cfg.NSQ, nil)
	if err == nil {
		producer = publisher
	}
}

// Return the global nsq publisher
func GetNSQPublisher() Publisher {
	return producer
}

type NSQPublisher struct {
	producer        *nsq.Producer
	publish         func(producer *nsq.Producer, topic string, data []byte) error
	deferredPublish func(producer *nsq.Producer, topic string, data []byte, delay time.Duration) error
	Prefix          string
}

func (n *NSQPublisher) FormatTopicName(topic string) string {
	if len(n.Prefix) == 0 {
		log.Printf("[MQ] FormatTopicName : Invalid Prefix %+v", n)
		return topic
	}
	return n.Prefix + topic
}

func nsqPublish(producer *nsq.Producer, topic string, data []byte) error {
	return producer.Publish(topic, data)
}
func nsqDeferredPublish(producer *nsq.Producer, topic string, data []byte, delay time.Duration) error {
	return producer.DeferredPublish(topic, delay, data)
}

func NewNSQPublisher(appCfg config.NSQCfg, config *nsq.Config) (*NSQPublisher, error) {
	if config == nil {
		config = nsq.NewConfig()
	}
	p, err := nsq.NewProducer(appCfg.PublishAddress, config)
	if err != nil {
		log.Printf("[MQ] NewNSQPublisher : Unable To Init Publisher %+v", err)
		return nil, err
	}
	logging.Debug.Printf("[MQ] NewNSQPublisher : Successfully connected to client producer %+v and prefix %+v", p, appCfg.Prefix)
	return &NSQPublisher{
		producer:        p,
		publish:         nsqPublish,
		deferredPublish: nsqDeferredPublish,
		Prefix:          appCfg.Prefix,
	}, nil
}

// Publish data to message queue.
func (n *NSQPublisher) Publish(topic string, data interface{}) (err error) {
	if data == nil {
		return fmt.Errorf("Data is nil")
	}
	var payload []byte
	payload, err = json.Marshal(data)
	if err != nil {
		log.Printf("[MQ] Publish : Unable to parse payload err %v", err)
		return fmt.Errorf("Error in Payload : %+v", err)
	}

	if n.producer == nil {
		log.Println("[MQ] Publish : Unexpexted Nil producer")
		return fmt.Errorf("No Producer Found")
	}

	//Info :: Prepand Not Empty Prefix
	topic = n.FormatTopicName(topic)

	fmt.Printf("[MQ] Publishing topic=%s payload=%s\n", topic, string(payload))
	return n.publish(n.producer, topic, payload)
}

func (n *NSQPublisher) PublishWithoutPrefix(topic string, data interface{}) (err error) {
	if data == nil {
		return fmt.Errorf("Data is nil")
	}
	var payload []byte
	payload, err = json.Marshal(data)
	if err != nil {
		log.Printf("[MQ] PublishWithoutPrefix : Unable to parse payload err %v", err)
		return fmt.Errorf("Error in Payload : %+v", err)
	}

	if n.producer == nil {
		log.Println("[MQ] PublishWithoutPrefix : Unexpexted Nil producer")
		return fmt.Errorf("No Producer Found")
	}

	fmt.Printf("[MQ] PublishWithoutPrefix topic=%s payload=%s\n", topic, string(payload))
	return n.publish(n.producer, topic, payload)
}

func (n *NSQPublisher) DeferredPublish(topic string, data interface{}, delay time.Duration) (err error) {
	if data == nil {
		return fmt.Errorf("Data is nil")
	}
	var payload []byte
	payload, err = json.Marshal(data)
	if err != nil {
		log.Printf("[MQ] DeferredPublish : Unable to parse payload err %v", err)
		return fmt.Errorf("Error in Payload : %+v", err)
	}
	if n.producer == nil {
		log.Println("[MQ] DeferredPublish : Unexpexted Nil producer")
		return fmt.Errorf("No Producer Found")
	}

	//Info :: Prepand Not Empty Prefix
	topic = n.FormatTopicName(topic)
	fmt.Printf("[MQ] DeferredPublish topic=%s payload=%s delay=%s\n", topic, string(payload), delay)
	return n.deferredPublish(n.producer, topic, payload, delay)
}

func (n *NSQPublisher) DeferredPublishWithoutPrefix(topic string, data interface{}, delay time.Duration) (err error) {
	if data == nil {
		return fmt.Errorf("Data is nil")
	}
	var payload []byte
	payload, err = json.Marshal(data)
	if err != nil {
		log.Printf("[MQ] DeferredPublishWithoutPrefix : Unable to parse payload err %v", err)
		return fmt.Errorf("Error in Payload : %+v", err)
	}
	if n.producer == nil {
		log.Println("[MQ] DeferredPublishWithoutPrefix : Unexpexted Nil producer")
		return fmt.Errorf("No Producer Found")
	}

	fmt.Printf("[MQ] DeferredPublishWithoutPrefix topic=%s payload=%s delay=%s\n", topic, string(payload), delay)
	return n.deferredPublish(n.producer, topic, payload, delay)
}
