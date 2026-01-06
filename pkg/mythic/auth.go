package mythic

import (
	"context"
	"fmt"
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

	// Perform login mutation
	var mutation struct {
		Login struct {
			AccessToken  string `graphql:"access_token"`
			RefreshToken string `graphql:"refresh_token"`
			User         struct {
				ID       int    `graphql:"id"`
				Username string `graphql:"username"`
			} `graphql:"user"`
		} `graphql:"login(username: $username, password: $password)"`
	}

	variables := map[string]interface{}{
		"username": c.config.Username,
		"password": c.config.Password,
	}

	// Execute login mutation without authentication (special case)
	err := c.graphqlClient.Mutate(ctx, &mutation, variables)
	if err != nil {
		return WrapError("Login", err, "login mutation failed")
	}

	// Check if we got tokens
	if mutation.Login.AccessToken == "" {
		return WrapError("Login", ErrAuthenticationFailed, "no access token returned")
	}

	// Store tokens in config
	c.config.AccessToken = mutation.Login.AccessToken
	c.config.RefreshToken = mutation.Login.RefreshToken
	c.authenticated = true

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
			Active     bool   `graphql:"active"`
		} `graphql:"createAPIToken"`
	}

	err := c.executeMutation(ctx, &mutation, nil)
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

	var query struct {
		Operator []struct {
			ID              int    `graphql:"id"`
			Username        string `graphql:"username"`
			Admin           bool   `graphql:"admin"`
			Active          bool   `graphql:"active"`
			CurrentOperation struct {
				ID   int    `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"current_operation"`
		} `graphql:"operator(limit: 1)"`
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
		CurrentOperation: &Operation{
			ID:   op.CurrentOperation.ID,
			Name: op.CurrentOperation.Name,
		},
	}

	// Update current operation ID
	c.SetCurrentOperation(op.CurrentOperation.ID)

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

// RefreshAccessToken refreshes the access token using the refresh token.
// This is called automatically when a request fails with an authentication error.
func (c *Client) RefreshAccessToken(ctx context.Context) error {
	c.authMutex.Lock()
	defer c.authMutex.Unlock()

	if c.config.RefreshToken == "" {
		return WrapError("RefreshAccessToken", ErrAuthenticationFailed, "no refresh token available")
	}

	var mutation struct {
		RefreshToken struct {
			AccessToken  string `graphql:"access_token"`
			RefreshToken string `graphql:"refresh_token"`
		} `graphql:"refreshToken(refresh_token: $refresh_token)"`
	}

	variables := map[string]interface{}{
		"refresh_token": c.config.RefreshToken,
	}

	// Execute without authentication since we're refreshing
	err := c.graphqlClient.Mutate(ctx, &mutation, variables)
	if err != nil {
		return WrapError("RefreshAccessToken", err, "failed to refresh token")
	}

	if mutation.RefreshToken.AccessToken == "" {
		return WrapError("RefreshAccessToken", ErrInvalidResponse, "no access token returned")
	}

	// Update tokens
	c.config.AccessToken = mutation.RefreshToken.AccessToken
	c.config.RefreshToken = mutation.RefreshToken.RefreshToken

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
