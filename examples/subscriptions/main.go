package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func main() {
	// Create configuration from environment variables
	config := &mythic.Config{
		ServerURL:     os.Getenv("MYTHIC_URL"),       // e.g., "https://mythic.example.com:7443"
		Username:      os.Getenv("MYTHIC_USERNAME"),  // for login
		Password:      os.Getenv("MYTHIC_PASSWORD"),  // for login
		APIToken:      os.Getenv("MYTHIC_API_TOKEN"), // or use API token
		SkipTLSVerify: os.Getenv("MYTHIC_SKIP_TLS_VERIFY") == "true",
		SSL:           true,
		Timeout:       30 * time.Second,
	}

	// Create client
	client, err := mythic.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Authenticate (uses username/password or API token from config)
	ctx := context.Background()
	err = client.Login(ctx)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	log.Println("Authenticated successfully")

	// Set current operation
	operationID := 1 // Set to your operation ID
	client.SetCurrentOperation(operationID)

	// Example 1: Subscribe to task output
	log.Println("=== Example 1: Task Output Subscription ===")
	taskOutputSub := subscribeToTaskOutput(client)
	defer taskOutputSub.Close()

	// Example 2: Subscribe to callback updates
	log.Println("\n=== Example 2: Callback Subscription ===")
	callbackSub := subscribeToCallbacks(client)
	defer callbackSub.Close()

	// Example 3: Subscribe to file uploads
	log.Println("\n=== Example 3: File Subscription ===")
	fileSub := subscribeToFiles(client)
	defer fileSub.Close()

	// Wait for interrupt signal
	log.Println("\n=== Listening for events (Ctrl+C to exit) ===")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Process events from all subscriptions
	for {
		select {
		case event := <-taskOutputSub.Events:
			handleTaskOutputEvent(event)
		case err := <-taskOutputSub.Errors:
			log.Printf("[Task Output] Error: %v", err)

		case event := <-callbackSub.Events:
			handleCallbackEvent(event)
		case err := <-callbackSub.Errors:
			log.Printf("[Callback] Error: %v", err)

		case event := <-fileSub.Events:
			handleFileEvent(event)
		case err := <-fileSub.Errors:
			log.Printf("[File] Error: %v", err)

		case <-taskOutputSub.Done:
			log.Println("[Task Output] Subscription closed")
			return
		case <-callbackSub.Done:
			log.Println("[Callback] Subscription closed")
			return
		case <-fileSub.Done:
			log.Println("[File] Subscription closed")
			return

		case <-sigChan:
			log.Println("\nReceived interrupt signal, shutting down...")
			return
		}
	}
}

// subscribeToTaskOutput creates a subscription for task output events
func subscribeToTaskOutput(client *mythic.Client) *types.Subscription {
	config := &types.SubscriptionConfig{
		Type: types.SubscriptionTypeTaskOutput,
		Handler: func(event *types.SubscriptionEvent) error {
			// Handler is called for each event
			// You can process the event here or just let it flow to the Events channel
			return nil
		},
		BufferSize: 100,
		// Optional: Add filters
		Filter: map[string]interface{}{
			// "callback_id": 42,  // Only events from specific callback
			// "task_id": 123,     // Only events from specific task
		},
	}

	ctx := context.Background()
	sub, err := client.Subscribe(ctx, config)
	if err != nil {
		log.Fatalf("Failed to subscribe to task output: %v", err)
	}

	log.Printf("Subscribed to task output (ID: %s)", sub.ID)
	return sub
}

// subscribeToCallbacks creates a subscription for callback events
func subscribeToCallbacks(client *mythic.Client) *types.Subscription {
	config := &types.SubscriptionConfig{
		Type: types.SubscriptionTypeCallback,
		Handler: func(event *types.SubscriptionEvent) error {
			return nil
		},
		BufferSize: 50,
	}

	ctx := context.Background()
	sub, err := client.Subscribe(ctx, config)
	if err != nil {
		log.Fatalf("Failed to subscribe to callbacks: %v", err)
	}

	log.Printf("Subscribed to callbacks (ID: %s)", sub.ID)
	return sub
}

// subscribeToFiles creates a subscription for file events
func subscribeToFiles(client *mythic.Client) *types.Subscription {
	config := &types.SubscriptionConfig{
		Type: types.SubscriptionTypeFile,
		Handler: func(event *types.SubscriptionEvent) error {
			return nil
		},
		BufferSize: 50,
	}

	ctx := context.Background()
	sub, err := client.Subscribe(ctx, config)
	if err != nil {
		log.Fatalf("Failed to subscribe to files: %v", err)
	}

	log.Printf("Subscribed to files (ID: %s)", sub.ID)
	return sub
}

// Event Handlers

func handleTaskOutputEvent(event *types.SubscriptionEvent) {
	output, _ := event.GetDataField("output")
	taskID, _ := event.GetDataField("task_id")

	log.Printf("[Task Output] Task %v: %v", taskID, output)
}

func handleCallbackEvent(event *types.SubscriptionEvent) {
	host, _ := event.GetDataField("host")
	user, _ := event.GetDataField("user")
	active, _ := event.GetDataField("active")

	log.Printf("[Callback] %s@%s (active: %v)", user, host, active)
}

func handleFileEvent(event *types.SubscriptionEvent) {
	filename, _ := event.GetDataField("filename")
	complete, _ := event.GetDataField("complete")
	isDownload, _ := event.GetDataField("is_download_from_agent")

	direction := "Upload"
	if isDownload.(bool) {
		direction = "Download"
	}

	status := "In Progress"
	if complete.(bool) {
		status = "Complete"
	}

	log.Printf("[File] %s: %s (%s)", direction, filename, status)
}
