package mythic

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"github.com/hasura/go-graphql-client"
)

// Client is the main Mythic SDK client.
type Client struct {
	// config holds the client configuration
	config *Config

	// graphqlClient is the underlying GraphQL client
	graphqlClient *graphql.Client

	// httpClient is the underlying HTTP client
	httpClient *http.Client

	// authenticated tracks whether the client is authenticated
	authenticated bool

	// authMutex protects authentication state
	authMutex sync.RWMutex

	// currentOperationID is the currently selected operation ID
	currentOperationID *int
}

// NewClient creates a new Mythic client with the provided configuration.
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, WrapError("NewClient", err, "invalid configuration")
	}

	// Create HTTP client with optional TLS skip verification
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	if config.SkipTLSVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// Construct GraphQL endpoint URL
	scheme := "https"
	if !config.SSL {
		scheme = "http"
	}
	graphqlURL := fmt.Sprintf("%s://%s/graphql/", scheme, stripScheme(config.ServerURL))

	// Create GraphQL client
	gqlClient := graphql.NewClient(graphqlURL, httpClient)

	client := &Client{
		config:        config,
		graphqlClient: gqlClient,
		httpClient:    httpClient,
		authenticated: false,
	}

	// If we have an API token or access token, consider authenticated
	if config.APIToken != "" || config.AccessToken != "" {
		client.authenticated = true
	}

	return client, nil
}

// IsAuthenticated returns whether the client is authenticated.
func (c *Client) IsAuthenticated() bool {
	c.authMutex.RLock()
	defer c.authMutex.RUnlock()
	return c.authenticated
}

// GetConfig returns a copy of the client configuration.
func (c *Client) GetConfig() Config {
	return *c.config
}

// Close closes the client and releases resources.
func (c *Client) Close() error {
	// Close any open connections
	c.httpClient.CloseIdleConnections()
	return nil
}

// SetCurrentOperation sets the current operation ID for the client.
// All subsequent API calls will use this operation context.
func (c *Client) SetCurrentOperation(operationID int) {
	c.authMutex.Lock()
	defer c.authMutex.Unlock()
	c.currentOperationID = &operationID
}

// GetCurrentOperation returns the current operation ID, if set.
func (c *Client) GetCurrentOperation() *int {
	c.authMutex.RLock()
	defer c.authMutex.RUnlock()
	if c.currentOperationID == nil {
		return nil
	}
	id := *c.currentOperationID
	return &id
}

// executeQuery executes a GraphQL query with authentication.
func (c *Client) executeQuery(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	if !c.IsAuthenticated() {
		return ErrNotAuthenticated
	}

	// Create a client with authentication headers
	client := c.getAuthenticatedClient()

	// Execute query
	return client.Query(ctx, query, variables)
}

// executeMutation executes a GraphQL mutation with authentication.
func (c *Client) executeMutation(ctx context.Context, mutation interface{}, variables map[string]interface{}) error {
	if !c.IsAuthenticated() {
		return ErrNotAuthenticated
	}

	// Create a client with authentication headers
	client := c.getAuthenticatedClient()

	// Execute mutation
	return client.Mutate(ctx, mutation, variables)
}

// getAuthenticatedClient returns a GraphQL client with authentication headers.
func (c *Client) getAuthenticatedClient() *graphql.Client {
	headers := c.getAuthHeaders()

	return c.graphqlClient.WithRequestModifier(func(req *http.Request) {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	})
}

// getAuthHeaders returns the authentication headers for API requests.
func (c *Client) getAuthHeaders() map[string]string {
	headers := make(map[string]string)

	// Use API token if available (preferred)
	if c.config.APIToken != "" {
		headers["apitoken"] = c.config.APIToken
	} else if c.config.AccessToken != "" {
		// Use JWT access token
		headers["Authorization"] = "Bearer " + c.config.AccessToken
	}

	return headers
}

// stripScheme removes http:// or https:// from a URL if present.
func stripScheme(url string) string {
	if len(url) > 8 && url[:8] == "https://" {
		return url[8:]
	}
	if len(url) > 7 && url[:7] == "http://" {
		return url[7:]
	}
	return url
}
