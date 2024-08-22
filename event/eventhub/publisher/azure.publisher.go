package publishser

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Somesh/go-boilerplate/common/config"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

// EventHubProducer manages multiple Event Hub clients and sends messages to specified partitions.
type EventHubProducer struct {
	hubs   map[string]*eventhub.Hub // Map of hub names to Event Hub clients
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	mu     sync.Mutex
}

// NewEventHubProducer initializes a new EventHubProducer with multiple hubs.
func NewEventHubProducer(connStrs map[string]string) (*EventHubProducer, error) {
	hubs := make(map[string]*eventhub.Hub)
	for hubName, connStr := range connStrs {
		hub, err := eventhub.NewHubFromConnectionString(connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create hub client for %s: %w", hubName, err)
		}
		hubs[hubName] = hub
	}

	_, cancel := context.WithCancel(context.Background())
	return &EventHubProducer{
		hubs:   hubs,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}, nil
}

// SendMessage sends a message to the specified Event Hub. Partitioning is managed by Event Hubs or by partition key.
func (p *EventHubProducer) SendMessage(hubName string, message string) error {
	p.mu.Lock()
	hub, ok := p.hubs[hubName]
	p.mu.Unlock()

	if !ok {
		return fmt.Errorf("hub not found: %s", hubName)
	}

	p.wg.Add(1)
	defer p.wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	event := eventhub.NewEventFromString(message)
	err := hub.Send(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to send message to hub %s: %w", hubName, err)
	}

	log.Printf("Message sent to hub %s: %s", hubName, message)
	return nil
}

// Close gracefully closes all Event Hub clients, ensuring all messages are sent.
func (p *EventHubProducer) Close() {
	p.cancel()
	p.wg.Wait()

	p.mu.Lock()
	defer p.mu.Unlock()
	for hubName, hub := range p.hubs {
		if err := hub.Close(context.Background()); err != nil {
			log.Printf("Error closing hub %s: %v", hubName, err)
		}
	}
}

// setupSignalHandler handles OS signals for graceful shutdown.
func setupSignalHandler(producer *EventHubProducer) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\nReceived interrupt signal, shutting down gracefully...")

		// Initiate graceful shutdown
		producer.Close()
		fmt.Println("Shutdown complete.")
	}()
}

// Example usage of the EventHubProducer
func init() {

	cfg := config.GetConfig()

	// connStrs := map[string]string{
	// 	"hub1": "your-connection-string-for-hub1",
	// 	"hub2": "your-connection-string-for-hub2",
	// 	// Add more hubs as needed
	// }

	var connStrs map[string]string = make(map[string]string)
	// var partitionMap map[string][]string = make(map[string][]string)
	connectionHub := cfg.ConnectionHub

	for connStr := range connectionHub {
		conn := connectionHub[connStr]
		connStrs[conn.HubName] = conn.HubConnString

		// partitionMap[conn.HubName] = conn.Partitions
	}

	// Initialize the EventHubProducer
	producer, err := NewEventHubProducer(connStrs)
	if err != nil {
		log.Fatalf("Failed to create EventHubProducer: %s", err)
	}

	// Setup signal handler for graceful shutdown
	setupSignalHandler(producer)

	// Example: Sending a message to multiple hubs
	for hubName := range connStrs {
		message := fmt.Sprintf("Hello from %s!", hubName)
		err := producer.SendMessage(hubName, message)
		if err != nil {
			log.Printf("Failed to send message: %s", err)
		}
	}

	// Continue with more business logic or wait for signal to shut down
	time.Sleep(10 * time.Second) // Simulate work

	// Manually close the producer if not using signal handling
	// producer.Close()
}
