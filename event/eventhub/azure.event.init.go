package event

// Azure Eventhub

import (
	"context"
	"log"
	"sync"
	"time"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/tokopedia/logging.v1"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"

	"github.com/Somesh/go-boilerplate/common/config"
	"github.com/Somesh/go-boilerplate/event/eventhub/consumer"
)

var (
	err error
)

type Handler func(c context.Context, event *eventhub.Event) error

type MQ struct {
	cfg                 *config.Config
	options             *Options
	partitionHandlerMap map[string]Handler
	birth               time.Time
}

type Options struct {
	HubConnString string
	HubName       string
	HubNameSpace  string
}

// New MQ
func New(o *Options, consmer *consumer.Consumer, cfg *config.Config) *MQ {

	m := &MQ{
		cfg:     cfg,
		options: o,
		birth:   time.Now(),
	}
	// Register consumers here
	m.partitionHandlerMap = make(map[string]Handler)
	m.partitionHandlerMap["0"] = consmer.HandlerA
	m.partitionHandlerMap["1"] = consmer.HandlerB
	return m
}

func (m *MQ) Run() {

	// hubName := m.cfg.Event.HubName
	// namespaceName := m.cfg.Event.HubNameSpace
	connStr := m.cfg.Event.HubConnString

	logging.Debug.Printf("Connection config %+v", m.cfg.Event)
	// Create a new Event Hub client
	hub, err := eventhub.NewHubFromConnectionString(connStr)
	if err != nil {
		log.Fatalf("Failed to create hub client: %s", err)
	}
	defer hub.Close(context.Background())

	// new
	// Context and wait group for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Get partition IDs from Event Hub runtime information
	runtimeInfo, err := hub.GetRuntimeInformation(ctx)
	if err != nil {
		log.Fatalf("Failed to get runtime information: %s", err)
	}

	// Start receiving messages from partitions
	startReceivingMessages(ctx, wg, hub, runtimeInfo.PartitionIDs, m.partitionHandlerMap)

	// Wait for interrupt signal for graceful shutdown
	<-signalChan
	fmt.Println("\nReceived interrupt signal, shutting down gracefully...")

	// Initiate graceful shutdown
	cancel()
	wg.Wait()
	fmt.Println("Shutdown complete.")
}

// startReceivingMessages starts receiving messages from each partition in a separate goroutine.
func startReceivingMessages(ctx context.Context, wg *sync.WaitGroup, hub *eventhub.Hub, partitionIDs []string, partitionHandlerMap map[string]Handler) {
	for _, partitionID := range partitionIDs {
		if handler, ok := partitionHandlerMap[partitionID]; ok {
			wg.Add(1)
			go func(partitionID string, handler Handler) {
				defer wg.Done()
				receiveMessages(ctx, hub, partitionID, handler)
			}(partitionID, handler)
		}
	}
}

// receiveMessages starts receiving messages for a specific partition and routes them to the appropriate handler.
func receiveMessages(ctx context.Context, hub *eventhub.Hub, partitionID string, handler Handler) {
	_, err := hub.Receive(ctx, partitionID, func(c context.Context, event *eventhub.Event) error {
		return handler(c, event)
	}, eventhub.ReceiveWithLatestOffset())

	if err != nil {
		log.Printf("Failed to receive messages from partition %s: %s", partitionID, err)
	} else {
		log.Printf("Stopped receiving messages from partition %s", partitionID)
	}
}

// setupSignalHandler sets up a signal channel for handling OS interrupts, waits for the signal, and then gracefully shuts down.
func setupSignalHandler(cancel context.CancelFunc, wg *sync.WaitGroup) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal for graceful shutdown
	go func() {
		<-signalChan
		fmt.Println("\nReceived interrupt signal, shutting down gracefully...")

		// Initiate graceful shutdown
		cancel()
		wg.Wait()
		fmt.Println("Shutdown complete.")
	}()
}
