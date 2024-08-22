package consumer

import (
	"context"
	"fmt"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"

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

// SS: HINT: In real implemenation pass on the core module to interact with its fuctions while processing the messages in consumer

// func Setup(mod *manager.Module, cfg *config.Config) *Consumer {

// 	return consumerModule(mod, cfg)
// }

// func consumerModule(mod *manager.Module, cfg *config.Config) *Consumer {

// 	return &Consumer{
// 		cfg:              cfg,
// 		managerModule:    mod,
// 		slackDebugLog:    slack.GetDebugLogger(),
// 		slackLog:         slack.GetLogger(),
// 		slackBusinessLog: slack.GetBusinessLogger(),
// 	}
// }

func (c *Consumer) HandlerA(ctx context.Context, event *eventhub.Event) error {
	fmt.Printf("Handler A processing message from Partition %s: %s\n", *event, string(event.Data))

	// fmt.Printf("Handler A processing message from Partition %s: %s\n", *event.SystemProperties.PartitionID, string(event.Data))
	// Add specific logic for Handler A
	return nil
}

func (c *Consumer) HandlerB(ctx context.Context, event *eventhub.Event) error {
	fmt.Printf("Handler B processing message from Partition %s: %s\n", *event.SystemProperties.PartitionID, string(event.Data))
	// Add specific logic for Handler A
	return nil
}
