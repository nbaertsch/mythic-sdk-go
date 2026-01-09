package mythic

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// Subscribe creates a GraphQL subscription for real-time event updates.
// Subscriptions use WebSocket connections to stream events like task output,
// callback status changes, and file uploads in real-time.
//
// Subscription Types:
//   - task_output: Real-time task output as it's generated
//   - callback: Callback status changes (new, active, dead, etc.)
//   - file: New file uploads and downloads
//   - all: All events across the operation
//
// The subscription runs in a goroutine and calls the provided handler for each event.
// Events are delivered through channels for thread-safe consumption.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - config: Subscription configuration including type, handler, and filters
//
// Returns:
//   - *types.Subscription: Active subscription with event/error channels
//   - error: Error if subscription creation fails
//
// Example:
//
//	// Subscribe to task output updates
//	config := &types.SubscriptionConfig{
//	    Type: types.SubscriptionTypeTaskOutput,
//	    Handler: func(event *types.SubscriptionEvent) error {
//	        fmt.Printf("Task output: %v\n", event.Data)
//	        return nil
//	    },
//	    BufferSize: 100,
//	}
//
//	sub, err := client.Subscribe(ctx, config)
//	if err != nil {
//	    return err
//	}
//	defer sub.Close()
//
//	// Process events until context is cancelled
//	for {
//	    select {
//	    case event := <-sub.Events:
//	        // Event is automatically handled by the handler
//	        fmt.Printf("Received event: %s\n", event.String())
//	    case err := <-sub.Errors:
//	        fmt.Printf("Subscription error: %v\n", err)
//	    case <-sub.Done:
//	        fmt.Println("Subscription closed")
//	        return nil
//	    case <-ctx.Done():
//	        fmt.Println("Context cancelled")
//	        return ctx.Err()
//	    }
//	}
//
// WebSocket Connection:
// Subscriptions use the graphql-transport-ws protocol over WebSocket connections.
// The connection is automatically established on first subscription and reused
// for subsequent subscriptions. Authentication is handled via connection parameters.
func (c *Client) Subscribe(ctx context.Context, config *types.SubscriptionConfig) (*types.Subscription, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, WrapError("Subscribe", ErrInvalidInput, err.Error())
	}

	// Set default buffer size
	bufferSize := config.BufferSize
	if bufferSize == 0 {
		bufferSize = 100
	}

	// Set operation ID if not specified
	var operationID int
	if config.OperationID != 0 {
		operationID = config.OperationID
	} else {
		opID := c.GetCurrentOperation()
		if opID == nil {
			return nil, WrapError("Subscribe", ErrNotAuthenticated, "no current operation set")
		}
		operationID = *opID
	}

	// Generate unique subscription ID
	subID := generateSubscriptionID()

	// Create subscription object
	sub := &types.Subscription{
		ID:     subID,
		Type:   config.Type,
		Active: true,
		Events: make(chan *types.SubscriptionEvent, bufferSize),
		Errors: make(chan error, 10),
		Done:   make(chan struct{}),
	}

	// Create context for subscription lifecycle
	subCtx, cancel := context.WithCancel(context.Background())

	// Get subscription client (establishes WebSocket connection if needed)
	subscriptionClient := c.getSubscriptionClient()

	// Build GraphQL subscription query based on type
	query, variables := buildSubscriptionQuery(config.Type, operationID, config.Filter)

	// Start subscription in background goroutine
	go func() {
		defer func() {
			// Clean up on exit
			sub.Close()
			cancel()

			// Remove from active subscriptions
			c.subscriptionsMutex.Lock()
			delete(c.activeSubscriptions, subID)
			c.subscriptionsMutex.Unlock()
		}()

		// Subscribe using WebSocket client
		graphqlSubID, err := subscriptionClient.Subscribe(query, variables, func(dataValue []byte, errValue error) error {
			// Handle errors from subscription
			if errValue != nil {
				select {
				case sub.Errors <- WrapError("Subscribe", ErrOperationFailed, errValue.Error()):
				case <-subCtx.Done():
					return subCtx.Err()
				}
				return errValue
			}

			// Parse event data
			event := &types.SubscriptionEvent{
				Type:      config.Type,
				Data:      make(map[string]interface{}),
				Timestamp: time.Now().Format(time.RFC3339),
			}

			// Parse JSON data into event
			if err := parseJSON(dataValue, &event.Data); err != nil {
				select {
				case sub.Errors <- WrapError("Subscribe", ErrOperationFailed, fmt.Sprintf("failed to parse event data: %v", err)):
				case <-subCtx.Done():
					return subCtx.Err()
				}
				return err
			}

			// Call user handler
			if config.Handler != nil {
				if err := config.Handler(event); err != nil {
					select {
					case sub.Errors <- WrapError("Subscribe", ErrOperationFailed, fmt.Sprintf("handler error: %v", err)):
					case <-subCtx.Done():
						return subCtx.Err()
					}
					// Continue processing even if handler returns error
				}
			}

			// Send event to channel
			select {
			case sub.Events <- event:
			case <-subCtx.Done():
				return subCtx.Err()
			}

			return nil
		})

		if err != nil {
			select {
			case sub.Errors <- WrapError("Subscribe", ErrOperationFailed, fmt.Sprintf("subscription failed: %v", err)):
			case <-subCtx.Done():
			}
			return
		}

		// Wait for cancellation
		<-subCtx.Done()

		// Unsubscribe using the graphqlSubID
		if graphqlSubID != "" {
			_ = subscriptionClient.Unsubscribe(graphqlSubID)
		}
	}()

	// Track active subscription
	c.subscriptionsMutex.Lock()
	c.activeSubscriptions[subID] = &subscriptionContext{
		cancel:  cancel,
		closeFn: nil, // unsubscribe function is handled in goroutine
	}
	c.subscriptionsMutex.Unlock()

	return sub, nil
}

// Unsubscribe closes an active subscription and cleans up resources.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - subscription: The subscription to close
//
// Returns:
//   - error: Error if unsubscribe fails
//
// Example:
//
//	err := client.Unsubscribe(ctx, sub)
//	if err != nil {
//	    log.Printf("Failed to unsubscribe: %v", err)
//	}
func (c *Client) Unsubscribe(ctx context.Context, subscription *types.Subscription) error {
	if subscription == nil {
		return WrapError("Unsubscribe", ErrInvalidInput, "subscription cannot be nil")
	}

	if !subscription.Active {
		return WrapError("Unsubscribe", ErrInvalidInput, "subscription is not active")
	}

	// Get subscription context
	c.subscriptionsMutex.RLock()
	subCtx, exists := c.activeSubscriptions[subscription.ID]
	c.subscriptionsMutex.RUnlock()

	if !exists {
		// Already unsubscribed or never existed
		subscription.Close()
		return nil
	}

	// Cancel subscription context (triggers cleanup in goroutine)
	if subCtx.cancel != nil {
		subCtx.cancel()
	}

	// Close subscription object
	subscription.Close()

	return nil
}

// generateSubscriptionID generates a unique subscription identifier.
func generateSubscriptionID() string {
	return uuid.New().String()
}

// buildSubscriptionQuery constructs a GraphQL subscription query based on type.
func buildSubscriptionQuery(subType types.SubscriptionType, operationID int, filter map[string]interface{}) (interface{}, map[string]interface{}) {
	variables := map[string]interface{}{
		"operation_id": operationID,
	}

	// Add filter parameters if provided
	if filter != nil {
		for k, v := range filter {
			variables[k] = v
		}
	}

	// Build query based on subscription type
	switch subType {
	case types.SubscriptionTypeTaskOutput:
		// Subscribe to task output updates
		var query struct {
			TaskOutput []struct {
				ID        int    `graphql:"id"`
				Output    string `graphql:"output"`
				Timestamp string `graphql:"timestamp"`
				TaskID    int    `graphql:"task_id"`
				Task      struct {
					ID              int    `graphql:"id"`
					Command         string `graphql:"command"`
					Params          string `graphql:"params"`
					OriginalParams  string `graphql:"original_params"`
					DisplayParams   string `graphql:"display_params"`
					Status          string `graphql:"status"`
					Timestamp       string `graphql:"timestamp"`
					CompletedTime   string `graphql:"completed_time"`
					CallbackID      int    `graphql:"callback_id"`
					OperatorID      int    `graphql:"operator_id"`
					CommentOperator string `graphql:"comment_operator"`
				} `graphql:"task"`
			} `graphql:"task_output(where: {task: {callback: {operation_id: {_eq: $operation_id}}}}, order_by: {id: desc})"`
		}
		return &query, variables

	case types.SubscriptionTypeCallback:
		// Subscribe to callback updates
		var query struct {
			Callback []struct {
				ID                  int    `graphql:"id"`
				DisplayID           int    `graphql:"display_id"`
				AgentCallbackID     string `graphql:"agent_callback_id"`
				InitCallback        string `graphql:"init_callback"`
				LastCheckin         string `graphql:"last_checkin"`
				User                string `graphql:"user"`
				Host                string `graphql:"host"`
				PID                 int    `graphql:"pid"`
				IP                  string `graphql:"ip"`
				ExternalIP          string `graphql:"external_ip"`
				ProcessName         string `graphql:"process_name"`
				Description         string `graphql:"description"`
				OperatorID          int    `graphql:"operator_id"`
				Active              bool   `graphql:"active"`
				RegisteredPayloadID int    `graphql:"registered_payload_id"`
				IntegrityLevel      int    `graphql:"integrity_level"`
				Locked              bool   `graphql:"locked"`
				OperationID         int    `graphql:"operation_id"`
				SleepInfo           string `graphql:"sleep_info"`
				Architecture        string `graphql:"architecture"`
				Domain              string `graphql:"domain"`
				Os                  string `graphql:"os"`
			} `graphql:"callback(where: {operation_id: {_eq: $operation_id}}, order_by: {id: desc})"`
		}
		return &query, variables

	case types.SubscriptionTypeFile:
		// Subscribe to file updates
		var query struct {
			FileMeta []struct {
				ID                  int    `graphql:"id"`
				AgentFileID         string `graphql:"agent_file_id"`
				TotalChunks         int    `graphql:"total_chunks"`
				ChunksReceived      int    `graphql:"chunks_received"`
				ChunkSize           int    `graphql:"chunk_size"`
				Path                string `graphql:"full_remote_path"`
				Host                string `graphql:"host"`
				IsDownloadFromAgent bool   `graphql:"is_download_from_agent"`
				IsScreenshot        bool   `graphql:"is_screenshot"`
				IsPayload           bool   `graphql:"is_payload"`
				Timestamp           string `graphql:"timestamp"`
				Complete            bool   `graphql:"complete"`
				Deleted             bool   `graphql:"deleted"`
				OperatorID          int    `graphql:"operator_id"`
				OperationID         int    `graphql:"operation_id"`
				TaskID              *int   `graphql:"task_id"`
				Filename            string `graphql:"filename_text"`
				Md5                 string `graphql:"md5"`
				Sha1                string `graphql:"sha1"`
			} `graphql:"filemeta(where: {operation_id: {_eq: $operation_id}}, order_by: {id: desc})"`
		}
		return &query, variables

	case types.SubscriptionTypeAlert:
		// Subscribe to operational alerts
		var query struct {
			OperationalAlert []struct {
				ID          int    `graphql:"id"`
				Message     string `graphql:"message"`
				Alert       string `graphql:"alert"`
				Source      string `graphql:"source"`
				Severity    int    `graphql:"severity"`
				Resolved    bool   `graphql:"resolved"`
				OperationID int    `graphql:"operation_id"`
				CallbackID  *int   `graphql:"callback_id"`
				Timestamp   string `graphql:"timestamp"`
			} `graphql:"operationalert(where: {operation_id: {_eq: $operation_id}}, order_by: {id: desc})"`
		}
		return &query, variables

	case types.SubscriptionTypeAll:
		// Subscribe to all events (task output, callbacks, files)
		// Note: This would require multiple subscriptions or a complex query
		// For now, default to task output as primary event type
		fallthrough

	default:
		// Default to task output subscription
		var query struct {
			TaskOutput []struct {
				ID        int    `graphql:"id"`
				Output    string `graphql:"output"`
				Timestamp string `graphql:"timestamp"`
				TaskID    int    `graphql:"task_id"`
			} `graphql:"task_output(where: {task: {callback: {operation_id: {_eq: $operation_id}}}}, order_by: {id: desc})"`
		}
		return &query, variables
	}
}
