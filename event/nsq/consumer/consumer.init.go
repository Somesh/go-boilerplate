package nsqConsumer

import (
	"fmt"
	"github.com/nsqio/go-nsq"

	"github.com/Somesh/go-boilerplate/common/config"
	"github.com/Somesh/go-boilerplate/src/manager"
)

type Consumer struct {
	cfg           *config.Config
	managerModule *manager.Module
}

func Setup(mod *manager.Module, cfg *config.Config) *Consumer {
	return consumerModule(mod, cfg)
}

func consumerModule(mod *manager.Module, cfg *config.Config) *Consumer {

	return &Consumer{
		cfg:           cfg,
		managerModule: mod,
	}
}

func (c *Consumer) HandlerA(msg *nsq.Message) error {
	fmt.Printf("[NSQ][HandlerA] Processing message : %s", string(msg.Body))

	// Your logic
	return nil
}
