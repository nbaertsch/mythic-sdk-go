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

## Generic Task Creation Architecture

### Overview

Mythic's command system is **dynamic and extensible**. Each payload type (agent) registers its own commands with parameter definitions. The SDK must **query these definitions** at runtime and adapt parameter formatting accordingly.

**DO NOT hardcode command-specific logic.** Always query command definitions from Mythic.

### Command Registration Flow

```
┌────────────────┐
│ Payload Type   │ (e.g., Poseidon)
│ Container      │
└────────┬───────┘
         │ 1. Container starts
         │ 2. Registers commands with Mythic
         ▼
┌────────────────────┐
│ Mythic Server      │
│                    │
│ ┌────────────────┐ │
│ │  command       │ │ (cmd, description, author, etc.)
│ └────────┬───────┘ │
│          │         │
│ ┌────────▼────────┐│
│ │commandparameters││ (name, type, required, default, etc.)
│ └─────────────────┘│
└────────────────────┘
         ▲
         │ 3. SDK queries command definitions
         │ 4. SDK builds params dynamically
┌────────┴────────┐
│ SDK Client      │
└─────────────────┘
```

### Two Types of Commands

#### Type 1: Parameterized Commands (JSON-based)

**Characteristics:**
- Has `CommandParameters` defined
- Expects params as **JSON object**
- Parameters have names, types, defaults, validation

**Example: `curl` command**

```go
// Poseidon registration (agentfunctions/curl.go)
Command Parameters: [
    {Name: "url", Type: "String", Required: true},
    {Name: "method", Type: "ChooseOne", Choices: ["GET", "POST", ...], Default: "GET"},
    {Name: "headers", Type: "Array"},
    {Name: "body", Type: "String"},
]

// SDK params format:
Params: `{"url": "https://example.com", "method": "POST", "body": "data"}`
```

#### Type 2: Raw String Commands

**Characteristics:**
- NO `CommandParameters` defined
- Expects params as **plain text string**
- No parameter validation or structure

**Example: `shell` command**

```go
// Poseidon registration (agentfunctions/shell.go)
var shell = agentstructs.Command{
    Name: "shell",
    // NO CommandParameters!
    TaskFunctionCreateTasking: shellCreateTasking,
}

// SDK params format:
Params: "whoami"  // Plain string, NOT {"command": "whoami"}
```

### Proper Task Creation Workflow

#### Step 1: Query Available Commands

```go
// For a specific callback (recommended - shows only loaded commands)
commands, err := client.GetLoadedCommandsForCallback(ctx, callbackID)

// Or query by payload type
type CommandsQuery struct {
    Command []struct {
        ID          int       `graphql:"id"`
        Cmd         string    `graphql:"cmd"`
        Description string    `graphql:"description"`
        Version     int       `graphql:"version"`
        Author      string    `graphql:"author"`
        HelpCmd     string    `graphql:"help_cmd"`
        CommandParameters []struct {
            ID           int         `graphql:"id"`
            Name         string      `graphql:"name"`
            DisplayName  string      `graphql:"display_name"`
            Type         string      `graphql:"type"`
            Required     bool        `graphql:"required"`
            DefaultValue interface{} `graphql:"default_value"`
            Choices      interface{} `graphql:"choices"`
            Description  string      `graphql:"description"`
        } `graphql:"commandparameters"`
    } `graphql:"command(where: {payload_type_id: {_eq: $payload_type_id}, deleted: {_eq: false}})"`
}
```

#### Step 2: Determine Command Type

```go
func IsRawStringCommand(command *types.LoadedCommand) bool {
    return len(command.Parameters) == 0
}

func GetParametersHelp(command *types.LoadedCommand) string {
    if IsRawStringCommand(command) {
        return fmt.Sprintf("%s: Plain text command (no structured parameters)", command.CommandName)
    }

    help := fmt.Sprintf("%s parameters:\n", command.CommandName)
    for _, param := range command.Parameters {
        required := ""
        if param.Required {
            required = " (REQUIRED)"
        }
        help += fmt.Sprintf("  - %s (%s)%s: %s\n", 
            param.Name, param.Type, required, param.Description)
    }
    return help
}
```

#### Step 3: Build Params Dynamically

```go
func BuildTaskParams(command *types.LoadedCommand, inputs map[string]interface{}) (string, error) {
    if IsRawStringCommand(command) {
        // Raw string command - params should be plain text
        if rawValue, ok := inputs["raw"]; ok {
            return fmt.Sprintf("%v", rawValue), nil
        }
        return "", fmt.Errorf("raw string command requires 'raw' input")
    }

    // Parameterized command - build JSON object
    paramsMap := make(map[string]interface{})

    for _, param := range command.Parameters {
        value, provided := inputs[param.Name]

        // Check required parameters
        if param.Required && !provided {
            return "", fmt.Errorf("required parameter missing: %s", param.Name)
        }

        // Use provided value or default
        if provided {
            paramsMap[param.Name] = value
        } else if param.DefaultValue != nil {
            paramsMap[param.Name] = param.DefaultValue
        }
    }

    // Marshal to JSON string
    paramsJSON, err := json.Marshal(paramsMap)
    if err != nil {
        return "", fmt.Errorf("failed to marshal params: %w", err)
    }

    return string(paramsJSON), nil
}
```

#### Step 4: Issue Task

```go
func IssueTaskDynamic(ctx context.Context, client *mythic.Client, 
                       callbackID int, commandName string, 
                       inputs map[string]interface{}) (*mythic.Task, error) {
    
    // 1. Get command definition
    commands, err := client.GetLoadedCommandsForCallback(ctx, callbackID)
    if err != nil {
        return nil, err
    }

    var targetCommand *types.LoadedCommand
    for _, cmd := range commands {
        if cmd.CommandName == commandName {
            targetCommand = cmd
            break
        }
    }

    if targetCommand == nil {
        return nil, fmt.Errorf("command %s not loaded for callback", commandName)
    }

    // 2. Build params based on command type
    params, err := BuildTaskParams(targetCommand, inputs)
    if err != nil {
        return nil, fmt.Errorf("invalid parameters: %w", err)
    }

    // 3. Issue task
    return client.IssueTask(ctx, &mythic.TaskRequest{
        CallbackID: &callbackID,
        Command:    commandName,
        Params:     params,
    })
}

// Usage examples:
// Raw string command:
task, err := IssueTaskDynamic(ctx, client, callbackID, "shell", map[string]interface{}{
    "raw": "whoami",
})

// Parameterized command:
task, err := IssueTaskDynamic(ctx, client, callbackID, "curl", map[string]interface{}{
    "url":    "https://example.com",
    "method": "POST",
    "body":   "test data",
})
```

### Parameter Types

Mythic supports these parameter types (from MythicContainer agent_structs):

```go
const (
    COMMAND_PARAMETER_TYPE_STRING           = "String"
    COMMAND_PARAMETER_TYPE_BOOLEAN          = "Boolean"
    COMMAND_PARAMETER_TYPE_NUMBER           = "Number"
    COMMAND_PARAMETER_TYPE_ARRAY            = "Array"
    COMMAND_PARAMETER_TYPE_CHOOSE_ONE       = "ChooseOne"
    COMMAND_PARAMETER_TYPE_CHOOSE_MULTIPLE  = "ChooseMultiple"
    COMMAND_PARAMETER_TYPE_FILE             = "File"
    COMMAND_PARAMETER_TYPE_CREDENTIAL_JSON  = "CredentialJson"
    COMMAND_PARAMETER_TYPE_LINKINFO         = "LinkInfo"
    COMMAND_PARAMETER_TYPE_PAYLOAD_LIST     = "PayloadList"
    COMMAND_PARAMETER_TYPE_AGGREGATE_DATA   = "AggregateBrowserScriptData"
    COMMAND_PARAMETER_TYPE_DATE             = "Date"
)
```

**Validation by Type:**

```go
func ValidateParameterValue(param *types.CommandParameter, value interface{}) error {
    switch param.Type {
    case "String":
        if _, ok := value.(string); !ok {
            return fmt.Errorf("%s must be a string", param.Name)
        }

    case "Boolean":
        if _, ok := value.(bool); !ok {
            return fmt.Errorf("%s must be a boolean", param.Name)
        }

    case "Number":
        switch value.(type) {
        case int, int64, float64:
            // Valid
        default:
            return fmt.Errorf("%s must be a number", param.Name)
        }

    case "Array":
        if _, ok := value.([]interface{}); !ok {
            return fmt.Errorf("%s must be an array", param.Name)
        }

    case "ChooseOne":
        strVal, ok := value.(string)
        if !ok {
            return fmt.Errorf("%s must be a string", param.Name)
        }
        // Validate against choices
        for _, choice := range param.Choices {
            if choice == strVal {
                return nil
            }
        }
        return fmt.Errorf("%s must be one of: %v", param.Name, param.Choices)

    case "ChooseMultiple":
        arr, ok := value.([]interface{})
        if !ok {
            return fmt.Errorf("%s must be an array", param.Name)
        }
        // Validate each choice
        for _, item := range arr {
            strItem, ok := item.(string)
            if !ok {
                return fmt.Errorf("%s items must be strings", param.Name)
            }
            found := false
            for _, choice := range param.Choices {
                if choice == strItem {
                    found = true
                    break
                }
            }
            if !found {
                return fmt.Errorf("%s: invalid choice '%s'", param.Name, strItem)
            }
        }
    }

    return nil
}
```

### Interactive CLI Example

```go
func InteractiveCLI(ctx context.Context, client *mythic.Client, callbackID int) error {
    // 1. List commands
    commands, err := client.GetLoadedCommandsForCallback(ctx, callbackID)
    if err != nil {
        return err
    }

    fmt.Println("\nAvailable commands:")
    for i, cmd := range commands {
        fmt.Printf("%2d. %-15s - %s\n", i+1, cmd.CommandName, cmd.Description)
    }

    // 2. Select command
    fmt.Print("\nSelect command (number): ")
    var cmdIndex int
    fmt.Scanf("%d", &cmdIndex)
    if cmdIndex < 1 || cmdIndex > len(commands) {
        return fmt.Errorf("invalid selection")
    }

    selectedCmd := commands[cmdIndex-1]
    fmt.Printf("\nSelected: %s\n%s\n", selectedCmd.CommandName, selectedCmd.Description)

    // 3. Build parameters
    inputs := make(map[string]interface{})

    if IsRawStringCommand(selectedCmd) {
        // Raw string command
        fmt.Print("\nEnter command: ")
        reader := bufio.NewReader(os.Stdin)
        line, _ := reader.ReadString('\n')
        inputs["raw"] = strings.TrimSpace(line)
    } else {
        // Parameterized command
        fmt.Println("\nEnter parameters (leave blank for default):")

        for _, param := range selectedCmd.Parameters {
            required := ""
            if param.Required {
                required = " [REQUIRED]"
            }

            fmt.Printf("\n%s (%s)%s:\n", param.Name, param.Type, required)
            fmt.Printf("  %s\n", param.Description)

            if param.DefaultValue != nil {
                fmt.Printf("  Default: %v\n", param.DefaultValue)
            }

            if len(param.Choices) > 0 {
                fmt.Printf("  Choices: %v\n", param.Choices)
            }

            fmt.Print("  Value: ")
            reader := bufio.NewReader(os.Stdin)
            line, _ := reader.ReadString('\n')
            line = strings.TrimSpace(line)

            if line == "" && !param.Required {
                continue
            }

            // Parse value based on type
            value, err := parseParameterValue(param.Type, line)
            if err != nil {
                fmt.Printf("  Error: %v\n", err)
                return err
            }

            // Validate
            if err := ValidateParameterValue(param, value); err != nil {
                fmt.Printf("  Error: %v\n", err)
                return err
            }

            inputs[param.Name] = value
        }
    }

    // 4. Build and preview params
    params, err := BuildTaskParams(selectedCmd, inputs)
    if err != nil {
        return fmt.Errorf("failed to build params: %w", err)
    }

    fmt.Printf("\n┌─ Task Preview ──────────────────┐\n")
    fmt.Printf("│ Command: %s\n", selectedCmd.CommandName)
    fmt.Printf("│ Params:  %s\n", params)
    fmt.Printf("└────────────────────────────────┘\n")

    // 5. Confirm
    fmt.Print("\nIssue task? (y/n): ")
    var confirm string
    fmt.Scanf("%s", &confirm)

    if strings.ToLower(confirm) != "y" {
        fmt.Println("Cancelled.")
        return nil
    }

    // 6. Issue task
    task, err := client.IssueTask(ctx, &mythic.TaskRequest{
        CallbackID: &callbackID,
        Command:    selectedCmd.CommandName,
        Params:     params,
    })

    if err != nil {
        return fmt.Errorf("task failed: %w", err)
    }

    fmt.Printf("\n✓ Task created successfully\n")
    fmt.Printf("  Display ID: %d\n", task.DisplayID)
    fmt.Printf("  Status: %s\n", task.Status)

    return nil
}
```

### SDK Implementation Checklist

- [ ] Add `GetCommandsForPayloadType(ctx, payloadTypeID)` method
- [ ] Add `Command` type with `Parameters []*CommandParameter`
- [ ] Add `CommandParameter` type with all fields from schema
- [ ] Add `IsRawStringCommand(command)` helper
- [ ] Add `BuildTaskParams(command, inputs)` helper
- [ ] Add `ValidateParameterValue(param, value)` helper
- [ ] Update `IssueTask()` documentation with param format requirements
- [ ] Add examples for both command types to README
- [ ] Add integration test for parameterized command (e.g., `curl`)
- [ ] Add integration test for raw string command (e.g., `shell`)

### Key Takeaways

1. **NEVER hardcode command logic** - always query command definitions
2. **Check parameter count** to determine if command is raw string or parameterized
3. **The `commandparameters` table is the source of truth** for all parameter metadata
4. **Use `GetLoadedCommandsForCallback()`** to see what's actually available
5. **Build params dynamically** based on command type and parameter definitions
6. **This approach works with ANY payload type** - not just Poseidon!

### References

- Mythic command registration: `/root/Mythic/InstalledServices/poseidon/poseidon/agentfunctions/`
- Mythic agent structs: `github.com/MythicMeta/MythicContainer/agent_structs`
- Hasura schema: `/root/Mythic/hasura-docker/metadata/databases/default/tables/`
- GraphQL introspection: `scripts/utils/introspect_schema.py --type command`
