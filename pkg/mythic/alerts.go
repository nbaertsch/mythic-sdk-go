package mythic

import (
	"context"
	"fmt"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetAlerts retrieves alerts for an operation with optional filtering.
//
// Alerts provide automated security monitoring and OPSEC notifications
// for suspicious activities, policy violations, or security events.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - operationID: ID of the operation (0 for current operation)
//   - filter: Optional filtering options (nil for defaults)
//
// Returns:
//   - []*types.Alert: List of alerts (most recent first)
//   - error: Error if operation ID is invalid or query fails
//
// Example:
//
//	// Get all unresolved high-severity alerts
//	resolved := false
//	minSeverity := types.AlertSeverityHigh
//	filter := &types.AlertFilter{
//	    Resolved:    &resolved,
//	    MinSeverity: &minSeverity,
//	    Limit:       50,
//	}
//	alerts, err := client.GetAlerts(ctx, 0, filter)
//	if err != nil {
//	    return err
//	}
//	for _, alert := range alerts {
//	    fmt.Printf("[%s] %s\n", alert.AlertType, alert.Message)
//	}
func (c *Client) GetAlerts(ctx context.Context, operationID int, filter *types.AlertFilter) ([]*types.Alert, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Use current operation if not specified
	if operationID == 0 {
		currentOp := c.GetCurrentOperation()
		if currentOp == nil {
			return nil, WrapError("GetAlerts", ErrNotAuthenticated, "no current operation set")
		}
		operationID = *currentOp
	}

	// Set defaults if no filter provided
	if filter == nil {
		filter = &types.AlertFilter{}
	}
	filter.SetDefaults()

	// Build variables
	variables := map[string]interface{}{
		"operation_id": operationID,
		"limit":        filter.Limit,
	}

	// We'll fetch all alerts and filter client-side for simplicity
	// since GraphQL dynamic where clauses are complex
	var query struct {
		OperationalAlert []struct {
			ID          int       `graphql:"id"`
			Message     string    `graphql:"message"`
			Alert       string    `graphql:"alert"`
			Source      string    `graphql:"source"`
			Severity    int       `graphql:"severity"`
			Resolved    bool      `graphql:"resolved"`
			OperationID int       `graphql:"operation_id"`
			CallbackID  *int      `graphql:"callback_id"`
			Timestamp   time.Time `graphql:"timestamp"`
		} `graphql:"operationalert(where: {operation_id: {_eq: $operation_id}}, order_by: {timestamp: desc}, limit: $limit)"`
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetAlerts", err, "failed to query alerts")
	}

	// Convert and apply client-side filters
	alerts := make([]*types.Alert, 0, len(query.OperationalAlert))
	for _, alertData := range query.OperationalAlert {
		// Apply filters
		if filter.AlertType != nil && alertData.Alert != *filter.AlertType {
			continue
		}
		if filter.Severity != nil && alertData.Severity != *filter.Severity {
			continue
		}
		if filter.MinSeverity != nil && alertData.Severity < *filter.MinSeverity {
			continue
		}
		if filter.Resolved != nil && alertData.Resolved != *filter.Resolved {
			continue
		}
		if filter.CallbackID != nil && (alertData.CallbackID == nil || *alertData.CallbackID != *filter.CallbackID) {
			continue
		}

		alert := &types.Alert{
			ID:          alertData.ID,
			Message:     alertData.Message,
			AlertType:   alertData.Alert,
			Source:      alertData.Source,
			Severity:    alertData.Severity,
			Resolved:    alertData.Resolved,
			OperationID: alertData.OperationID,
			CallbackID:  alertData.CallbackID,
			Timestamp:   alertData.Timestamp,
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// GetAlertByID retrieves a specific alert by its ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - alertID: ID of the alert to retrieve
//
// Returns:
//   - *types.Alert: The alert object
//   - error: Error if alert ID is invalid or not found
//
// Example:
//
//	alert, err := client.GetAlertByID(ctx, 42)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Alert: %s\n", alert.String())
func (c *Client) GetAlertByID(ctx context.Context, alertID int) (*types.Alert, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if alertID == 0 {
		return nil, WrapError("GetAlertByID", ErrInvalidInput, "alert ID is required")
	}

	var query struct {
		OperationalAlert []struct {
			ID          int       `graphql:"id"`
			Message     string    `graphql:"message"`
			Alert       string    `graphql:"alert"`
			Source      string    `graphql:"source"`
			Severity    int       `graphql:"severity"`
			Resolved    bool      `graphql:"resolved"`
			OperationID int       `graphql:"operation_id"`
			CallbackID  *int      `graphql:"callback_id"`
			Timestamp   time.Time `graphql:"timestamp"`
		} `graphql:"operationalert(where: {id: {_eq: $alert_id}})"`
	}

	variables := map[string]interface{}{
		"alert_id": alertID,
	}

	err := c.executeQuery(ctx, &query, variables)
	if err != nil {
		return nil, WrapError("GetAlertByID", err, "failed to query alert")
	}

	if len(query.OperationalAlert) == 0 {
		return nil, WrapError("GetAlertByID", ErrNotFound, fmt.Sprintf("alert %d not found", alertID))
	}

	alertData := query.OperationalAlert[0]
	return &types.Alert{
		ID:          alertData.ID,
		Message:     alertData.Message,
		AlertType:   alertData.Alert,
		Source:      alertData.Source,
		Severity:    alertData.Severity,
		Resolved:    alertData.Resolved,
		OperationID: alertData.OperationID,
		CallbackID:  alertData.CallbackID,
		Timestamp:   alertData.Timestamp,
	}, nil
}

// GetUnresolvedAlerts retrieves all unresolved (active) alerts for an operation.
//
// This is a convenience method for monitoring active security events.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - operationID: ID of the operation (0 for current operation)
//
// Returns:
//   - []*types.Alert: List of unresolved alerts (most recent first)
//   - error: Error if operation ID is invalid or query fails
//
// Example:
//
//	alerts, err := client.GetUnresolvedAlerts(ctx, 0)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Active alerts: %d\n", len(alerts))
//	for _, alert := range alerts {
//	    if alert.IsCritical() {
//	        fmt.Printf("CRITICAL: %s\n", alert.Message)
//	    }
//	}
func (c *Client) GetUnresolvedAlerts(ctx context.Context, operationID int) ([]*types.Alert, error) {
	resolved := false
	filter := &types.AlertFilter{
		Resolved: &resolved,
		Limit:    1000, // High limit for unresolved alerts
	}
	return c.GetAlerts(ctx, operationID, filter)
}

// ResolveAlert marks an alert as resolved/acknowledged.
//
// This indicates the alert has been reviewed and addressed by an operator.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - req: Resolution request with alert ID and optional notes
//
// Returns:
//   - error: Error if request is invalid or update fails
//
// Example:
//
//	req := &types.ResolveAlertRequest{
//	    AlertID: 42,
//	    Notes:   "False positive - expected behavior",
//	}
//	err := client.ResolveAlert(ctx, req)
//	if err != nil {
//	    return err
//	}
//	fmt.Println("Alert resolved successfully")
func (c *Client) ResolveAlert(ctx context.Context, req *types.ResolveAlertRequest) error {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return err
	}

	if req == nil {
		return WrapError("ResolveAlert", ErrInvalidInput, "request is required")
	}

	if err := req.Validate(); err != nil {
		return WrapError("ResolveAlert", ErrInvalidInput, err.Error())
	}

	// Update alert to mark as resolved
	var mutation struct {
		UpdateOperationalert struct {
			Affected_rows int `graphql:"affected_rows"`
		} `graphql:"update_operationalert(where: {id: {_eq: $alert_id}}, _set: {resolved: true})"`
	}

	variables := map[string]interface{}{
		"alert_id": req.AlertID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return WrapError("ResolveAlert", err, "failed to resolve alert")
	}

	if mutation.UpdateOperationalert.Affected_rows == 0 {
		return WrapError("ResolveAlert", ErrNotFound, fmt.Sprintf("alert %d not found or already resolved", req.AlertID))
	}

	return nil
}

// CreateCustomAlert creates a new manual alert in the operation.
//
// This allows operators to manually flag issues, suspicious activities,
// or important events that should be tracked as alerts.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - req: Alert creation request with message, type, source, and severity
//
// Returns:
//   - *types.Alert: The created alert
//   - error: Error if request is invalid or creation fails
//
// Example:
//
//	req := &types.CreateAlertRequest{
//	    Message:     "Unusual login pattern detected",
//	    AlertType:   types.AlertTypeOPSEC,
//	    Source:      "manual",
//	    Severity:    types.AlertSeverityHigh,
//	    OperationID: operationID,
//	}
//	alert, err := client.CreateCustomAlert(ctx, req)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Alert created: %s\n", alert.String())
func (c *Client) CreateCustomAlert(ctx context.Context, req *types.CreateAlertRequest) (*types.Alert, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	if req == nil {
		return nil, WrapError("CreateCustomAlert", ErrInvalidInput, "request is required")
	}

	if err := req.Validate(); err != nil {
		return nil, WrapError("CreateCustomAlert", ErrInvalidInput, err.Error())
	}

	// Create alert
	var mutation struct {
		InsertOperationalertOne struct {
			ID          int       `graphql:"id"`
			Message     string    `graphql:"message"`
			Alert       string    `graphql:"alert"`
			Source      string    `graphql:"source"`
			Severity    int       `graphql:"severity"`
			Resolved    bool      `graphql:"resolved"`
			OperationID int       `graphql:"operation_id"`
			CallbackID  *int      `graphql:"callback_id"`
			Timestamp   time.Time `graphql:"timestamp"`
		} `graphql:"insert_operationalert_one(object: {message: $message, alert: $alert, source: $source, severity: $severity, operation_id: $operation_id, callback_id: $callback_id, resolved: false})"`
	}

	variables := map[string]interface{}{
		"message":      req.Message,
		"alert":        req.AlertType,
		"source":       req.Source,
		"severity":     req.Severity,
		"operation_id": req.OperationID,
		"callback_id":  req.CallbackID,
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return nil, WrapError("CreateCustomAlert", err, "failed to create alert")
	}

	return &types.Alert{
		ID:          mutation.InsertOperationalertOne.ID,
		Message:     mutation.InsertOperationalertOne.Message,
		AlertType:   mutation.InsertOperationalertOne.Alert,
		Source:      mutation.InsertOperationalertOne.Source,
		Severity:    mutation.InsertOperationalertOne.Severity,
		Resolved:    mutation.InsertOperationalertOne.Resolved,
		OperationID: mutation.InsertOperationalertOne.OperationID,
		CallbackID:  mutation.InsertOperationalertOne.CallbackID,
		Timestamp:   mutation.InsertOperationalertOne.Timestamp,
	}, nil
}

// GetAlertStatistics retrieves aggregated alert statistics for an operation.
//
// This provides overview metrics for monitoring operation health and security.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - operationID: ID of the operation (0 for current operation)
//
// Returns:
//   - map[string]int: Statistics including total, unresolved, by_severity, by_type
//   - error: Error if operation ID is invalid or query fails
//
// Example:
//
//	stats, err := client.GetAlertStatistics(ctx, 0)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Total alerts: %d\n", stats["total"])
//	fmt.Printf("Unresolved: %d\n", stats["unresolved"])
//	fmt.Printf("Critical: %d\n", stats["critical"])
func (c *Client) GetAlertStatistics(ctx context.Context, operationID int) (map[string]int, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Use current operation if not specified
	if operationID == 0 {
		currentOp := c.GetCurrentOperation()
		if currentOp == nil {
			return nil, WrapError("GetAlertStatistics", ErrNotAuthenticated, "no current operation set")
		}
		operationID = *currentOp
	}

	// Fetch all alerts for the operation
	allAlerts, err := c.GetAlerts(ctx, operationID, &types.AlertFilter{Limit: 10000})
	if err != nil {
		return nil, WrapError("GetAlertStatistics", err, "failed to fetch alerts")
	}

	// Calculate statistics
	stats := map[string]int{
		"total":      len(allAlerts),
		"unresolved": 0,
		"resolved":   0,
		"critical":   0,
		"high":       0,
		"medium":     0,
		"low":        0,
		"info":       0,
		"opsec":      0,
		"error":      0,
		"warning":    0,
	}

	for _, alert := range allAlerts {
		// Resolution status
		if alert.Resolved {
			stats["resolved"]++
		} else {
			stats["unresolved"]++
		}

		// Severity counts
		switch alert.Severity {
		case types.AlertSeverityCritical:
			stats["critical"]++
		case types.AlertSeverityHigh:
			stats["high"]++
		case types.AlertSeverityMedium:
			stats["medium"]++
		case types.AlertSeverityLow:
			stats["low"]++
		case types.AlertSeverityInfo:
			stats["info"]++
		}

		// Type counts
		switch alert.AlertType {
		case types.AlertTypeOPSEC:
			stats["opsec"]++
		case types.AlertTypeError:
			stats["error"]++
		case types.AlertTypeWarning:
			stats["warning"]++
		}
	}

	return stats, nil
}

// SubscribeToAlerts creates a real-time subscription for alert events.
//
// This is a convenience wrapper around Subscribe() specifically for alerts.
// Alerts will be delivered through the subscription's Events channel.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - operationID: ID of the operation (0 for current operation)
//   - handler: Function called for each alert event
//
// Returns:
//   - *types.Subscription: Active subscription
//   - error: Error if subscription creation fails
//
// Example:
//
//	sub, err := client.SubscribeToAlerts(ctx, 0, func(event *types.SubscriptionEvent) error {
//	    alertType, _ := event.GetDataField("alert")
//	    message, _ := event.GetDataField("message")
//	    severity, _ := event.GetDataField("severity")
//	    fmt.Printf("[ALERT] %s (severity %v): %s\n", alertType, severity, message)
//	    return nil
//	})
//	if err != nil {
//	    return err
//	}
//	defer sub.Close()
//
//	// Process events
//	for {
//	    select {
//	    case event := <-sub.Events:
//	        // Event handled by handler function
//	    case err := <-sub.Errors:
//	        log.Printf("Subscription error: %v", err)
//	    case <-sub.Done:
//	        return nil
//	    }
//	}
func (c *Client) SubscribeToAlerts(ctx context.Context, operationID int, handler func(*types.SubscriptionEvent) error) (*types.Subscription, error) {
	// Use current operation if not specified
	if operationID == 0 {
		currentOp := c.GetCurrentOperation()
		if currentOp == nil {
			return nil, WrapError("SubscribeToAlerts", ErrNotAuthenticated, "no current operation set")
		}
		operationID = *currentOp
	}

	config := &types.SubscriptionConfig{
		Type:        types.SubscriptionTypeAlert,
		Handler:     handler,
		BufferSize:  100,
		OperationID: operationID,
	}

	return c.Subscribe(ctx, config)
}
