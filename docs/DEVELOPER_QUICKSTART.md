# Developer Quick-Start Guide

Quick-start guide for developers working on the Mythic SDK for Go.

## Project Structure

```
mythic-sdk-go/
├── pkg/mythic/           # Core SDK implementation
│   ├── auth.go           # Authentication (Login, RefreshAccessToken, CreateAPIToken)
│   ├── callbacks.go      # Callback operations
│   ├── files.go          # File upload/download
│   ├── tasks.go          # Task management
│   ├── client.go         # Client initialization and GraphQL execution
│   ├── config.go         # Client configuration
│   ├── errors.go         # Error types and wrapping
│   └── types/            # Type definitions
├── tests/
│   ├── unit/             # Unit tests (no live Mythic required)
│   └── integration/      # Integration tests (requires live Mythic)
├── scripts/utils/
│   └── introspect_schema.py  # GraphQL schema introspection tool
└── .github/workflows/
    ├── test.yml          # Unit tests (fast, matrix: 3 Go versions × 3 OSes)
    └── integration.yml   # Integration tests (slow, builds full Mythic)
```

## Setting Up Local Mythic

Integration tests require a running Mythic instance:

```bash
# Clone Mythic
git clone https://github.com/its-a-feature/Mythic.git /root/Mythic
cd /root/Mythic

# Build and start
sudo make
sudo ./mythic-cli start

# Get admin password
grep MYTHIC_ADMIN_PASSWORD .env | cut -d'=' -f2 | tr -d '"'

# Set environment variables
export MYTHIC_URL="https://127.0.0.1:7443"
export MYTHIC_USERNAME="mythic_admin"
export MYTHIC_PASSWORD="<password from above>"
export MYTHIC_SKIP_TLS_VERIFY="true"
```

## Running Tests Locally

### Unit Tests (Fast, No Mythic Required)

```bash
# Run all unit tests
go test -v ./tests/unit/...

# Run specific test
go test -v ./tests/unit/ -run TestCallbackString

# Check coverage
go test -cover ./tests/unit/...
```

### Integration Tests (Requires Live Mythic)

```bash
# Run all integration tests
go test -v -tags=integration ./tests/integration/...

# Run specific test
go test -v -tags=integration ./tests/integration/ -run TestAuthentication_Login

# Skip tests that need callbacks
# (Tests auto-skip if no callbacks exist in Mythic)
```

**Common integration test helpers:**
- `NewTestClient(t)` - Creates unauthenticated client
- `AuthenticateTestClient(t)` - Creates and authenticates client
- `GetTestConfig(t)` - Gets config from environment variables
- `SkipIfNoMythic(t)` - Skips test if Mythic unavailable

## Referencing Mythic Source Code

**Critical:** Always check Mythic source code when implementing new features or debugging.

### Mythic is a REST + GraphQL Hybrid

```
REST Endpoints:
  /auth                    - Login (username/password)
  /refresh                 - Refresh tokens (access_token + refresh_token)
  /api/v1.4/files/upload   - File upload (multipart/form-data)
  /api/v1.4/files/download/:uuid - File download

GraphQL Endpoint:
  /graphql/                - All queries and mutations (except auth/files)
```

### Finding Mythic Implementations

```bash
# Find Mythic source code location
ls -la /root/Mythic/mythic-docker/src/

# Search for endpoint handlers
grep -r "POST.*auth" /root/Mythic/mythic-docker/src/webserver/
grep -r "files/download" /root/Mythic/mythic-docker/src/webserver/

# Find specific function implementations
find /root/Mythic/mythic-docker/src -name "*.go" -exec grep -l "RefreshJWT" {} \;

# Read handler implementation
cat /root/Mythic/mythic-docker/src/webserver/controllers/login.go
cat /root/Mythic/mythic-docker/src/authentication/mythicjwt/jwt.go
```

### Mythic Source Code Locations

```
/root/Mythic/mythic-docker/src/
├── authentication/       # JWT, refresh tokens, middleware
├── webserver/
│   ├── initialize.go     # Route definitions
│   └── controllers/      # Endpoint handlers
├── database/
│   └── structs/          # Database models
└── grpc/                 # gRPC services
```

## GraphQL Schema Introspection

Use the introspection script to discover available types, fields, and mutations:

### List All Mutations

```bash
cd scripts/utils
python3 introspect_schema.py --mutations

# Filter mutations by name
python3 introspect_schema.py --mutations --filter callback
python3 introspect_schema.py --mutations --filter token
```

### Introspect Specific Types

```bash
# Introspect callback type
python3 introspect_schema.py --type callback

# Introspect operator type
python3 introspect_schema.py --type operator

# Introspect task type
python3 introspect_schema.py --type task
```

**Output shows:**
- Field names and types
- Whether fields are nullable
- Nested object structures
- List vs scalar types

### Using Hasura Console (Alternative)

Access Hasura GraphQL console directly:
```
http://localhost:8080/console
```

Use the GraphiQL interface to:
- Test queries/mutations
- View schema documentation
- See query suggestions

## Writing Unit Tests

Unit tests validate types, helpers, and logic without requiring Mythic.

### Test File Structure

```go
package unit

import (
    "testing"
    "github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
    "github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

func TestCallbackString(t *testing.T) {
    cb := &types.Callback{
        User:   "admin",
        Host:   "DC01",
        OS:     "Windows",
        Active: true,
    }

    expected := "admin@DC01 (Windows, active)"
    if cb.String() != expected {
        t.Errorf("Callback.String() = %q, want %q", cb.String(), expected)
    }
}
```

### Table-Driven Tests

```go
func TestCallbackIsHigh(t *testing.T) {
    tests := []struct {
        name           string
        integrityLevel types.CallbackIntegrityLevel
        want           bool
    }{
        {"low integrity", types.IntegrityLevelLow, false},
        {"high integrity", types.IntegrityLevelHigh, true},
        {"system integrity", types.IntegrityLevelSystem, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cb := &types.Callback{IntegrityLevel: tt.integrityLevel}
            if got := cb.IsHigh(); got != tt.want {
                t.Errorf("Callback.IsHigh() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Writing Integration Tests

Integration tests validate SDK against live Mythic instance.

### Test File Structure

```go
//go:build integration

package integration

import (
    "context"
    "testing"
    "time"
)

func TestFeature_Operation(t *testing.T) {
    SkipIfNoMythic(t)  // Skip if Mythic not available

    client := AuthenticateTestClient(t)  // Auto-authenticates

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Test your operation
    result, err := client.SomeOperation(ctx, params)
    if err != nil {
        t.Fatalf("SomeOperation failed: %v", err)
    }

    // Validate result
    if result.Field != expected {
        t.Errorf("Expected %v, got %v", expected, result.Field)
    }

    t.Logf("Operation succeeded: %v", result)
}
```

### Error Case Testing

```go
func TestFeature_Operation_NotFound(t *testing.T) {
    SkipIfNoMythic(t)
    client := AuthenticateTestClient(t)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Test error case
    _, err := client.SomeOperation(ctx, "nonexistent-id")
    if err == nil {
        t.Fatal("Expected error for non-existent resource, got nil")
    }

    // Optionally check error type
    if err != mythic.ErrNotFound {
        t.Errorf("Expected ErrNotFound, got: %v", err)
    }

    t.Logf("Expected error: %v", err)
}
```

## Monitoring GitHub Actions

### View Recent Runs

```bash
# List recent test runs
gh run list --workflow=test.yml --limit 5

# List integration test runs
gh run list --workflow=integration.yml --limit 5
```

### Watch Live Run

```bash
# Start a live watch
gh run watch <run-id> --interval 10

# View run details
gh run view <run-id>

# View specific job
gh run view <run-id> --job=<job-id>
```

### Check Test Results

```bash
# View full logs
gh run view <run-id> --log

# View only test output
gh run view <run-id> --job=<job-id> --log | grep -E "RUN|PASS|FAIL"

# Check for failures
gh run view <run-id> --log | grep "FAIL"

# Count failures
gh run view <run-id> --log | grep -c "FAIL"

# View specific test results
gh run view <run-id> --log | grep -A 10 "TestAuthentication_Login"
```

### Download Test Artifacts

```bash
# List artifacts
gh run view <run-id> --log-failed

# Download test results
gh run download <run-id>
```

## Common Patterns & Gotchas

### 1. REST vs GraphQL

**Mythic uses BOTH:**
- **REST:** Authentication (`/auth`, `/refresh`), file operations
- **GraphQL:** Everything else (queries, mutations)

**Authentication headers differ:**
```go
// For GraphQL
req.Header.Set("Authorization", "Bearer " + accessToken)
// OR
req.Header.Set("apitoken", apiToken)

// For REST (file downloads)
// Requires BOTH headers AND session cookies
httpClient := &http.Client{
    Jar: cookieJar,  // Critical for REST endpoints!
}
```

### 2. Mythic Returns HTTP 200 for Errors

**Problem:** Many Mythic endpoints return HTTP 200 with JSON error messages:

```json
{
  "status": "error",
  "error": "Failed to find file"
}
```

**Solution:** Always check for JSON error responses:

```go
if resp.StatusCode != http.StatusOK {
    return WrapError("Operation", ErrInvalidResponse, "...")
}

// Read body
body, _ := io.ReadAll(resp.Body)

// Check for JSON error
if body[0] == '{' {
    var errorResp struct {
        Status string `json:"status"`
        Error  string `json:"error"`
    }
    if json.Unmarshal(body, &errorResp) == nil && errorResp.Status == "error" {
        return WrapError("Operation", ErrNotFound, errorResp.Error)
    }
}
```

### 3. Session Cookies Required

**Problem:** File upload/download fail with 401 even with valid auth headers.

**Solution:** Add cookie jar to HTTP client:

```go
jar, _ := cookiejar.New(nil)
httpClient := &http.Client{
    Timeout: config.Timeout,
    Jar:     jar,  // Maintains session from /auth endpoint
}
```

### 4. Mutex Deadlocks

**Problem:** Nested lock acquisition causes deadlock:

```go
func (c *Client) Operation(ctx context.Context) error {
    c.authMutex.Lock()
    defer c.authMutex.Unlock()

    // This calls IsAuthenticated() which tries to acquire authMutex again!
    err := c.executeMutation(ctx, &mutation, variables)  // DEADLOCK!
}
```

**Solution:** Call GraphQL directly with manual auth headers:

```go
func (c *Client) Operation(ctx context.Context) error {
    c.authMutex.Lock()
    defer c.authMutex.Unlock()

    // Manually add auth headers to avoid nested lock
    headers := c.getAuthHeaders()
    client := c.graphqlClient.WithRequestModifier(func(req *http.Request) {
        for key, value := range headers {
            req.Header.Set(key, value)
        }
    })

    err := client.Mutate(ctx, &mutation, variables)
}
```

### 5. GraphQL Field Names

**Mythic uses snake_case in GraphQL:**

```go
// Correct
type Callback struct {
    DisplayID int `graphql:"display_id"`
    AgentCallbackID string `graphql:"agent_callback_id"`
}

// Wrong
type Callback struct {
    DisplayID int `graphql:"displayId"`  // Will fail!
}
```

Use introspection script to verify exact field names.

### 6. Pointer Fields for Optional Parameters

**Mutations with optional fields require pointers:**

```go
type CallbackUpdateRequest struct {
    CallbackDisplayID int      // Required (int)
    Active            *bool    // Optional (use pointer)
    Description       *string  // Optional (use pointer)
    IPs               []string // Optional (use nil slice)
}

// Usage
active := true
description := "High-value target"
req := &CallbackUpdateRequest{
    CallbackDisplayID: 1,
    Active:            &active,
    Description:       &description,
    IPs:               []string{"10.0.0.1"},
}
```

### 7. Integration Test Timing

**Integration tests timeout after 30 seconds by default:**

```go
// Adjust timeout for slow operations
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()
```

### 8. Test Cleanup

**Always clean up resources created during tests:**

```go
func TestFiles_Upload(t *testing.T) {
    client := AuthenticateTestClient(t)

    fileID, err := client.UploadFile(ctx, data, "test.txt")
    if err != nil {
        t.Fatalf("Upload failed: %v", err)
    }

    // Clean up
    defer func() {
        if err := client.DeleteFile(ctx, fileID); err != nil {
            t.Logf("Warning: Failed to delete test file: %v", err)
        }
    }()

    // Test operations...
}
```

## Debugging Workflow

1. **Reproduce locally** with unit/integration tests
2. **Check Mythic source code** for endpoint implementation
3. **Use introspection** to verify GraphQL schema
4. **Test with curl** if REST endpoint issue
5. **Check Mythic logs** for server-side errors:
   ```bash
   cd /root/Mythic
   sudo ./mythic-cli logs
   ```
6. **Verify in CI** via GitHub Actions

## Commit Message Format

Use clear, descriptive commit messages with context:

```
Fix RefreshAccessToken to use REST /refresh endpoint

Root cause analysis from Mythic source code:
- File: mythic-docker/src/webserver/initialize.go
  Route: protected.POST("/refresh", webcontroller.RefreshJWT)
- Previous implementation incorrectly used GraphQL mutation
- Mythic refresh is a REST endpoint, not GraphQL

Implementation:
- Changed to POST to /refresh with JSON body
- Request: {"access_token": "...", "refresh_token": "..."}
- Requires authentication header
- Response structure matches login response

Fixes TestAuthentication_RefreshAccessToken

🤖 Generated with Claude Code

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

## Getting Help

- **Mythic Documentation:** https://docs.mythic-c2.net/
- **Mythic Source Code:** /root/Mythic/mythic-docker/src/
- **GraphQL Introspection:** scripts/utils/introspect_schema.py
- **Hasura Console:** http://localhost:8080/console (when Mythic running)
- **GitHub Issues:** Check existing issues for similar problems

## Quick Reference

```bash
# Run unit tests
go test -v ./tests/unit/...

# Run integration tests
go test -v -tags=integration ./tests/integration/...

# Introspect GraphQL
cd scripts/utils && python3 introspect_schema.py --mutations

# Check Mythic source
grep -r "endpoint_name" /root/Mythic/mythic-docker/src/

# Watch CI run
gh run watch <run-id>

# Check test failures
gh run view <run-id> --log | grep "FAIL"
```
