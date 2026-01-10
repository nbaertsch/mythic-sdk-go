package mythic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LoginResponse represents the response from a login mutation.
type LoginResponse struct {
	Login struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"login"`
}

// CreateAPITokenResponse represents the response from creating an API token.
type CreateAPITokenResponse struct {
	CreateAPIToken struct {
		ID         int    `json:"id"`
		TokenValue string `json:"token_value"`
		Active     bool   `json:"active"`
	} `json:"createAPIToken"`
}

// Login authenticates the client with the Mythic server.
// It uses the credentials provided in the client configuration.
// After successful login, the client is marked as authenticated.
func (c *Client) Login(ctx context.Context) error {
	c.authMutex.Lock()
	defer c.authMutex.Unlock()

	// If API token is provided, we don't need to login
	if c.config.APIToken != "" {
		c.authenticated = true
		return nil
	}

	// If access token is already set, we're authenticated
	if c.config.AccessToken != "" {
		c.authenticated = true
		return nil
	}

	// Must have username and password
	if c.config.Username == "" || c.config.Password == "" {
		return WrapError("Login", ErrAuthenticationFailed, "username and password required")
	}

	// Construct auth endpoint URL
	scheme := "https"
	if !c.config.SSL {
		scheme = "http"
	}
	authURL := fmt.Sprintf("%s://%s/auth", scheme, stripScheme(c.config.ServerURL))

	// Prepare login request
	loginReq := map[string]string{
		"username": c.config.Username,
		"password": c.config.Password,
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return WrapError("Login", err, "failed to marshal login request")
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return WrapError("Login", err, "failed to create login request")
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return WrapError("Login", err, "failed to execute login request")
	}
	defer resp.Body.Close() //nolint:errcheck // Response body close error not critical

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WrapError("Login", err, "failed to read login response")
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return WrapError("Login", ErrAuthenticationFailed, fmt.Sprintf("login failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse response
	var authResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ID                     int    `json:"id"`
			UserID                 int    `json:"user_id"`
			Username               string `json:"username"`
			CurrentOperationID     int    `json:"current_operation_id"`
			CurrentOperation       string `json:"current_operation"`
			CurrentOperationBanner string `json:"current_operation_banner_text"`
		} `json:"user"`
	}

	if err := json.Unmarshal(body, &authResp); err != nil {
		return WrapError("Login", err, "failed to parse login response")
	}

	// Check if we got tokens
	if authResp.AccessToken == "" {
		return WrapError("Login", ErrAuthenticationFailed, "no access token returned")
	}

	// Store tokens in config
	c.config.AccessToken = authResp.AccessToken
	c.config.RefreshToken = authResp.RefreshToken
	c.authenticated = true

	// Store current operation ID
	if authResp.User.CurrentOperationID > 0 {
		c.currentOperationID = &authResp.User.CurrentOperationID
	}

	return nil
}

// CreateAPIToken creates a new API token for the authenticated user.
// Returns the token value which should be saved for future use.
func (c *Client) CreateAPIToken(ctx context.Context) (string, error) {
	if !c.IsAuthenticated() {
		return "", ErrNotAuthenticated
	}

	var mutation struct {
		CreateAPIToken struct {
			ID         int    `graphql:"id"`
			TokenValue string `graphql:"token_value"`
		} `graphql:"createAPIToken(token_type: $token_type)"`
	}

	variables := map[string]interface{}{
		"token_type": "User", // User token type for API access
	}

	err := c.executeMutation(ctx, &mutation, variables)
	if err != nil {
		return "", WrapError("CreateAPIToken", err, "failed to create API token")
	}

	if mutation.CreateAPIToken.TokenValue == "" {
		return "", WrapError("CreateAPIToken", ErrInvalidResponse, "no token value returned")
	}

	return mutation.CreateAPIToken.TokenValue, nil
}

// GetMe returns information about the currently authenticated user.
func (c *Client) GetMe(ctx context.Context) (*Operator, error) {
	if !c.IsAuthenticated() {
		return nil, ErrNotAuthenticated
	}

	// Query operator - filter by operation if we have one
	var query struct {
		Operator []struct {
			ID                 int    `graphql:"id"`
			Username           string `graphql:"username"`
			Admin              bool   `graphql:"admin"`
			Active             bool   `graphql:"active"`
			CurrentOperationID int    `graphql:"current_operation_id"`
		} `graphql:"operator(order_by: {id: asc}, limit: 1)"`
	}

	err := c.executeQuery(ctx, &query, nil)
	if err != nil {
		return nil, WrapError("GetMe", err, "failed to get current user")
	}

	if len(query.Operator) == 0 {
		return nil, WrapError("GetMe", ErrNotFound, "current user not found")
	}

	op := query.Operator[0]
	operator := &Operator{
		ID:       op.ID,
		Username: op.Username,
		Admin:    op.Admin,
		Active:   op.Active,
	}

	// Set current operation if available
	if op.CurrentOperationID > 0 {
		operator.CurrentOperation = &Operation{
			ID: op.CurrentOperationID,
		}
		c.SetCurrentOperation(op.CurrentOperationID)
	}

	return operator, nil
}

// Operator represents a Mythic operator (user).
type Operator struct {
	ID               int        `json:"id"`
	Username         string     `json:"username"`
	Admin            bool       `json:"admin"`
	Active           bool       `json:"active"`
	CurrentOperation *Operation `json:"current_operation,omitempty"`
}

// Operation represents a Mythic operation.
type Operation struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
	Webhook  string `json:"webhook,omitempty"`
}

// Logout clears authentication state.
// Note: This does not revoke tokens on the server.
func (c *Client) Logout() {
	c.authMutex.Lock()
	defer c.authMutex.Unlock()

	c.config.AccessToken = ""
	c.config.RefreshToken = ""
	c.authenticated = false
	c.currentOperationID = nil
}

// ClearRefreshToken clears only the refresh token.
// This is primarily for testing scenarios where you need to simulate a missing refresh token.
func (c *Client) ClearRefreshToken() {
	c.authMutex.Lock()
	defer c.authMutex.Unlock()

	c.config.RefreshToken = ""
}

// RefreshAccessToken refreshes the access token using the refresh token.
// This uses the REST /refresh endpoint, not GraphQL.
// Requires both access_token and refresh_token from the initial login.
func (c *Client) RefreshAccessToken(ctx context.Context) error {
	c.authMutex.Lock()
	defer c.authMutex.Unlock()

	if c.config.RefreshToken == "" {
		return WrapError("RefreshAccessToken", ErrAuthenticationFailed, "no refresh token available")
	}

	// Construct refresh endpoint URL
	scheme := "https"
	if !c.config.SSL {
		scheme = "http"
	}
	refreshURL := fmt.Sprintf("%s://%s/refresh", scheme, stripScheme(c.config.ServerURL))

	// Prepare refresh request
	refreshReq := map[string]string{
		"access_token":  c.config.AccessToken,
		"refresh_token": c.config.RefreshToken,
	}

	jsonData, err := json.Marshal(refreshReq)
	if err != nil {
		return WrapError("RefreshAccessToken", err, "failed to marshal refresh request")
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", refreshURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return WrapError("RefreshAccessToken", err, "failed to create refresh request")
	}
	req.Header.Set("Content-Type", "application/json")

	// Add authentication header (refresh endpoint requires authentication)
	if c.config.APIToken != "" {
		req.Header.Set("apitoken", c.config.APIToken)
	} else if c.config.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return WrapError("RefreshAccessToken", err, "failed to execute refresh request")
	}
	defer resp.Body.Close() //nolint:errcheck // Response body close error not critical

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WrapError("RefreshAccessToken", err, "failed to read refresh response")
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return WrapError("RefreshAccessToken", ErrAuthenticationFailed, fmt.Sprintf("refresh failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse response (same structure as login response)
	var refreshResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ID                     int    `json:"id"`
			UserID                 int    `json:"user_id"`
			Username               string `json:"username"`
			CurrentOperationID     int    `json:"current_operation_id"`
			CurrentOperation       string `json:"current_operation"`
			CurrentOperationBanner string `json:"current_operation_banner_text"`
		} `json:"user"`
	}

	if err := json.Unmarshal(body, &refreshResp); err != nil {
		return WrapError("RefreshAccessToken", err, "failed to parse refresh response")
	}

	// Check if we got tokens
	if refreshResp.AccessToken == "" {
		return WrapError("RefreshAccessToken", ErrAuthenticationFailed, "no access token returned")
	}

	// Update tokens in config
	c.config.AccessToken = refreshResp.AccessToken
	c.config.RefreshToken = refreshResp.RefreshToken

	// Update current operation ID if provided
	if refreshResp.User.CurrentOperationID > 0 {
		c.currentOperationID = &refreshResp.User.CurrentOperationID
	}

	return nil
}

// EnsureAuthenticated checks if the client is authenticated and attempts to login if not.
func (c *Client) EnsureAuthenticated(ctx context.Context) error {
	if c.IsAuthenticated() {
		return nil
	}

	return c.Login(ctx)
}

// String returns a string representation of the operator.
func (o *Operator) String() string {
	if o.CurrentOperation != nil {
		return fmt.Sprintf("%s (ID: %d, Admin: %t, Operation: %s)",
			o.Username, o.ID, o.Admin, o.CurrentOperation.Name)
	}
	return fmt.Sprintf("%s (ID: %d, Admin: %t)", o.Username, o.ID, o.Admin)
}
