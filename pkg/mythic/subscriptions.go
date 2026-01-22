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
			sub.Close() //nolint:errcheck // Cleanup function, error not critical
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
			if err := subscriptionClient.Unsubscribe(graphqlSubID); err != nil {
				// Log error but continue cleanup
				select {
				case sub.Errors <- WrapError("Subscribe", ErrOperationFailed, fmt.Sprintf("unsubscribe error: %v", err)):
				default:
				}
			}
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
		subscription.Close() //nolint:errcheck // Best effort cleanup
		return nil
	}

	// Cancel subscription context (triggers cleanup in goroutine)
	if subCtx.cancel != nil {
		subCtx.cancel()
	}

	// Close subscription object
	subscription.Close() //nolint:errcheck // Best effort cleanup

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
	for k, v := range filter {
		variables[k] = v
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

	case types.SubscriptionTypeScreenshot:
		// Subscribe to screenshot uploads (filemeta with is_screenshot=true)
		var query struct {
			FileMeta []struct {
				ID          int    `graphql:"id"`
				AgentFileID string `graphql:"agent_file_id"`
				Filename    string `graphql:"filename_text"`
				Path        string `graphql:"full_remote_path"`
				Host        string `graphql:"host"`
				Timestamp   string `graphql:"timestamp"`
				Complete    bool   `graphql:"complete"`
				TaskID      *int   `graphql:"task_id"`
				CallbackID  *int   `graphql:"callback_id"`
				OperationID int    `graphql:"operation_id"`
			} `graphql:"filemeta(where: {operation_id: {_eq: $operation_id}, is_screenshot: {_eq: true}, deleted: {_eq: false}}, order_by: {id: desc})"`
		}
		return &query, variables

	case types.SubscriptionTypeKeylog:
		// Subscribe to keylog entries
		var query struct {
			Keylog []struct {
				ID          int    `graphql:"id"`
				TaskID      int    `graphql:"task_id"`
				Keystrokes  string `graphql:"keystrokes"`
				Window      string `graphql:"window"`
				Timestamp   string `graphql:"timestamp"`
				OperationID int    `graphql:"operation_id"`
				User        string `graphql:"user"`
				CallbackID  int    `graphql:"callback_id"`
			} `graphql:"keylog(where: {operation_id: {_eq: $operation_id}}, order_by: {id: desc})"`
		}
		return &query, variables

	case types.SubscriptionTypeProcess:
		// Subscribe to process tracking updates
		var query struct {
			Process []struct {
				ID              int    `graphql:"id"`
				Name            string `graphql:"name"`
				ProcessID       int    `graphql:"process_id"`
				ParentProcessID int    `graphql:"parent_process_id"`
				Architecture    string `graphql:"architecture"`
				BinPath         string `graphql:"bin_path"`
				User            string `graphql:"user"`
				CommandLine     string `graphql:"command_line"`
				IntegrityLevel  int    `graphql:"integrity_level"`
				OperationID     int    `graphql:"operation_id"`
				HostID          int    `graphql:"host_id"`
				CallbackID      *int   `graphql:"callback_id"`
				Timestamp       string `graphql:"timestamp"`
				Deleted         bool   `graphql:"deleted"`
			} `graphql:"process(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {id: desc})"`
		}
		return &query, variables

	case types.SubscriptionTypeCredential:
		// Subscribe to credential discoveries
		var query struct {
			Credential []struct {
				ID          int    `graphql:"id"`
				Type        string `graphql:"type"`
				Account     string `graphql:"account"`
				Realm       string `graphql:"realm"`
				Credential  string `graphql:"credential"`
				Comment     string `graphql:"comment"`
				OperationID int    `graphql:"operation_id"`
				OperatorID  int    `graphql:"operator_id"`
				TaskID      *int   `graphql:"task_id"`
				Timestamp   string `graphql:"timestamp"`
				Deleted     bool   `graphql:"deleted"`
				Metadata    string `graphql:"metadata"`
			} `graphql:"credential(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {id: desc})"`
		}
		return &query, variables

	case types.SubscriptionTypeArtifact:
		// Subscribe to artifact/IOC tracking
		var query struct {
			Artifact []struct {
				ID           int    `graphql:"id"`
				Artifact     string `graphql:"artifact"`
				BaseArtifact string `graphql:"base_artifact"`
				Host         string `graphql:"host"`
				ArtifactType string `graphql:"artifact_type"`
				OperationID  int    `graphql:"operation_id"`
				OperatorID   int    `graphql:"operator_id"`
				TaskID       *int   `graphql:"task_id"`
				Timestamp    string `graphql:"timestamp"`
				Deleted      bool   `graphql:"deleted"`
				Metadata     string `graphql:"metadata"`
			} `graphql:"artifact(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {id: desc})"`
		}
		return &query, variables

	case types.SubscriptionTypeToken:
		// Subscribe to token discoveries
		var query struct {
			Token []struct {
				ID         int    `graphql:"id"`
				TokenID    string `graphql:"token_id"`
				User       string `graphql:"user"`
				Groups     string `graphql:"groups"`
				Privileges string `graphql:"privileges"`
				ThreadID   int    `graphql:"thread_id"`
				ProcessID  int    `graphql:"process_id"`
				SessionID  int    `graphql:"session_id"`
				LogonSID   string `graphql:"logon_sid"`
				// IntegrityLevel field removed - integrity_level_int not available in Mythic v3.4.20 schema
				Restricted  bool   `graphql:"restricted"`
				TaskID      *int   `graphql:"task_id"`
				OperationID int    `graphql:"operation_id"`
				Timestamp   string `graphql:"timestamp"`
				Host        string `graphql:"host"`
				Deleted     bool   `graphql:"deleted"`
			} `graphql:"token(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {id: desc})"`
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
