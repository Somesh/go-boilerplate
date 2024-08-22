package nsqrun

import (
	"testing"

	"github.com/nsqio/go-nsq"

	"github.com/Somesh/go-boilerplate/common/config"

	nsqConsumer "github.com/Somesh/go-boilerplate/event/nsq/consumer"
)

func dummyHandler(message *nsq.Message) error {
	return nil
}
func TestRegister(t *testing.T) {
	cfg := config.GetConfig()
	nsqcfg := Options{
		ListenAddress:  cfg.NSQ.ListenAddress,
		PublishAddress: cfg.NSQ.PublishAddress,
		Prefix:         cfg.NSQ.Prefix,
		LookUpAddress:  cfg.NSQ.LookUpAddress,
	}
	cons := &consumer.Consumer{}
	mq := New(&nsqcfg, cons)
	consumers := len(mq.consumers)
	mq.Register(constant.PaymentReconTopic, "recon", dummyHandler, nil)
	if consumers+1 != len(mq.consumers) {
		t.Errorf("Consumer Not Registered")
	}
	mq.RegisterWithoutPrefix(constant.PaymentReconTopic, "recon", dummyHandler, nil)
	if consumers+2 != len(mq.consumers) {
		t.Errorf("Consumer Not Registered")
	}
}
