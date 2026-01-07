# Integration Tests

This directory contains integration tests that run against a live Mythic C2 server instance.

## Overview

Integration tests verify the SDK's behavior against a real Mythic server, ensuring:
- API calls work correctly end-to-end
- Authentication flows function properly
- Data serialization/deserialization is correct
- Error handling works in real scenarios
- GraphQL queries match server expectations

## Running Integration Tests

### Option 1: Using Docker Compose (Recommended)

The easiest way to run integration tests is using the included Docker Compose setup:

```bash
# Start Mythic test instance
cd tests/integration
docker-compose up -d

# Wait for Mythic to be ready (can take 30-60 seconds)
# The server is ready when https://localhost:7443 responds

# Run integration tests
cd ../..
go test -v -tags=integration ./tests/integration/...

# Stop Mythic when done
cd tests/integration
docker-compose down -v
```

### Option 2: Using an Existing Mythic Instance

You can run tests against any Mythic server by setting environment variables:

```bash
export MYTHIC_URL="https://your-mythic-server:7443"
export MYTHIC_USERNAME="your_username"
export MYTHIC_PASSWORD="your_password"
export MYTHIC_SKIP_TLS_VERIFY="true"  # Set to false if using valid certs

go test -v -tags=integration ./tests/integration/...
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MYTHIC_URL` | `https://localhost:7443` | Mythic server URL |
| `MYTHIC_USERNAME` | `mythic_admin` | Username for authentication |
| `MYTHIC_PASSWORD` | `mythic_password` | Password for authentication |
| `MYTHIC_SKIP_TLS_VERIFY` | `true` | Skip TLS certificate verification |

## Test Organization

### helpers.go
Common test utilities:
- `GetTestConfig()` - Load configuration from environment
- `NewTestClient()` - Create unauthenticated test client
- `AuthenticateTestClient()` - Create authenticated test client
- `SkipIfNoMythic()` - Skip test if Mythic unavailable

### auth_test.go
Authentication tests:
- Login with username/password
- API token creation and usage
- GetMe() endpoint
- Authentication state management
- Invalid credentials handling

### callbacks_test.go
Callback operation tests:
- List all callbacks
- Filter active callbacks
- Get callback by ID
- Update callback properties
- Integrity level helpers
- Field validation

## CI/CD Integration

Integration tests run automatically in GitHub Actions:
- **Pull Requests**: On every PR to main/develop
- **Push to Main/Develop**: On direct pushes
- **Manual Trigger**: Via workflow_dispatch
- **Release Tags**: On version tags (v*)

### CI Workflow Details

The GitHub Actions workflow (`.github/workflows/integration.yml`) performs these steps:

1. **Clone Mythic Framework** - Shallow clone from GitHub (~30s)
2. **Build mythic-cli** - Compile the Mythic CLI tool (~45s)
3. **Start Mythic Services** - Launch all containers (~90s)
4. **Wait for Ready** - Poll until services are up (~30s)
5. **Extract Credentials** - Read admin password from `.env`
6. **Run Tests** - Execute integration test suite (~15s)
7. **Cleanup** - Stop containers and prune Docker

**Total CI Time**: ~4 minutes per run

**Key Features**:
- Tests against latest Mythic main branch
- Automatic credential extraction
- Artifact upload for debugging (test logs, Mythic logs)
- Always cleanup, even on failure
- 15-minute timeout to prevent hung jobs

See `.github/workflows/integration.yml` for the complete implementation.

## Writing New Integration Tests

When adding new integration tests:

1. **Use build tags**: Add `//go:build integration` at the top
2. **Check server availability**: Use `SkipIfNoMythic(t)` to skip if server unavailable
3. **Use helpers**: Leverage `AuthenticateTestClient(t)` for authenticated clients
4. **Handle empty state**: Fresh Mythic instances have no callbacks/tasks
5. **Clean up**: Restore state after destructive operations
6. **Use timeouts**: Always use `context.WithTimeout()`
7. **Log useful info**: Use `t.Logf()` for debugging output

Example:

```go
//go:build integration

package integration

import (
    "context"
    "testing"
    "time"
)

func TestYourFeature(t *testing.T) {
    SkipIfNoMythic(t)

    client := AuthenticateTestClient(t)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Your test code here
}
```

## Troubleshooting

### Tests Skip with "Mythic server not available"
- Ensure Mythic is running: `docker-compose ps`
- Check logs: `docker-compose logs mythic`
- Verify connectivity: `curl -k https://localhost:7443`

### Authentication Failures
- Verify credentials in docker-compose.yml match test defaults
- Check `MYTHIC_ADMIN_PASSWORD` environment variable in docker-compose.yml
- Ensure Mythic initialization completed (check logs)

### Slow Test Execution
- Integration tests are slower than unit tests (network calls)
- First run may take longer as Docker pulls images
- Mythic startup can take 30-60 seconds

### Connection Refused
- Mythic may still be starting - wait and retry
- Check port mapping: `docker-compose ps` should show `0.0.0.0:7443->7443/tcp`
- Verify no other service is using port 7443

## Test Coverage

Integration tests complement unit tests by covering:
- ✅ Real GraphQL API interactions
- ✅ Network error handling
- ✅ TLS/SSL configuration
- ✅ Authentication flows
- ✅ Data format compatibility
- ✅ Server response validation

They do NOT replace unit tests, which should still cover:
- Edge cases and error conditions
- Internal logic and algorithms
- Mock/stub scenarios
- Fast-running test cases
