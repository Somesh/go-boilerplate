package nsqPublisher

import "time"

type QMessage struct {
	Source  string      `json:"source"`
	Message interface{} `json:"message"`
}

type Publisher interface {
	Publish(topic string, data interface{}) error
	DeferredPublish(topic string, data interface{}, delay time.Duration) error
	PublishWithoutPrefix(topic string, data interface{}) error
	DeferredPublishWithoutPrefix(topic string, data interface{}, delay time.Duration) error
}
