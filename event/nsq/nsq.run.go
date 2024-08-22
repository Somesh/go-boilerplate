package nsqrun

import (
	"os"
	"strconv"
	"time"

	"fmt"

	"github.com/nsqio/go-nsq"

	nsqConsumer "github.com/Somesh/go-boilerplate/event/nsq/consumer"
)

var (
	err error
)

type subscriber struct {
	topic      string
	channel    string
	handler    func(message *nsq.Message) error
	concurrent int
	config     *nsq.Config
}

type MQ struct {
	options     *Options
	consumers   []subscriber
	listenErrCh chan error
	birth       time.Time
}

type Options struct {
	ListenAddress      []string
	LookUpAddress      []string
	PublishAddress     string
	PublishUberAddress string
	Prefix             string
}

// New MQ
func New(o *Options, consumer *nsqConsumer.Consumer) *MQ {

	m := &MQ{
		options:     o,
		listenErrCh: make(chan error),
		birth:       time.Now(),
	}

	// Register consumers here
	m.Register("topic_for_handerA", "channel_for_handlerA", consumer.HandlerA, nil)

	return m
}

// Register consumer
func (m *MQ) Register(topic, channel string, handler func(message *nsq.Message) error, config *nsq.Config, concurrents ...int) {
	var concurrent int
	if len(concurrents) > 0 && concurrents[0] > 0 {
		concurrent = concurrents[0]
	}

	topic = m.options.Prefix + topic
	channel = m.options.Prefix + channel
	fmt.Printf("Registering MQ Consumer : %s/%s concurrent=%d\n", topic, channel, concurrent)

	m.consumers = append(m.consumers, subscriber{
		topic:      topic,
		channel:    channel,
		handler:    handler,
		concurrent: concurrent,
		config:     config,
	})
}

func (m *MQ) RegisterWithoutPrefix(topic, channel string, handler func(message *nsq.Message) error, config *nsq.Config, concurrents ...int) {
	var concurrent int
	if len(concurrents) > 0 && concurrents[0] > 0 {
		concurrent = concurrents[0]
	}

	fmt.Printf("Registering MQ Consumer : %s/%s concurrent=%d\n", topic, channel, concurrent)

	m.consumers = append(m.consumers, subscriber{
		topic:      topic,
		channel:    channel,
		handler:    handler,
		concurrent: concurrent,
		config:     config,
	})
}

// Run listener
func (m *MQ) Run() {
	config := nsq.NewConfig()
	if m.options.PublishAddress == "" {
		m.options.PublishAddress = m.options.ListenAddress[0]
	}
	c, _ := strconv.Atoi(os.Getenv("CONC"))
	if c == 0 {
		c = 1
	}
	config.MaxInFlight = c
	config.MaxAttempts = 5
	for _, consumer := range m.consumers {
		consumerConfig := consumer.config
		if consumerConfig == nil {
			consumerConfig = config
		}
		q, _ := nsq.NewConsumer(consumer.topic, consumer.channel, consumerConfig)

		if consumer.concurrent != 0 {
			q.AddConcurrentHandlers(nsq.HandlerFunc(consumer.handler), consumer.concurrent)
		} else {
			q.AddHandler(nsq.HandlerFunc(consumer.handler))
		}

		var err error
		if len(m.options.LookUpAddress) > 1 {
			err = q.ConnectToNSQLookupds(m.options.LookUpAddress)
		} else {
			err = q.ConnectToNSQLookupd(m.options.LookUpAddress[0])
		}

		if err != nil {
			m.listenErrCh <- err
		}
	}
}
