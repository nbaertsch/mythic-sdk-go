#!/bin/bash

# Integration Test Script
# Starts Mythic via Docker Compose and runs integration tests

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
INTEGRATION_DIR="$PROJECT_ROOT/tests/integration"

echo "=== Mythic Go SDK Integration Tests ==="
echo ""

# Check dependencies
command -v docker >/dev/null 2>&1 || { echo "Error: docker is required but not installed."; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "Error: docker-compose is required but not installed."; exit 1; }

# Navigate to integration test directory
cd "$INTEGRATION_DIR"

echo "Starting Mythic test instance..."
docker-compose up -d

echo "Waiting for Mythic to be ready..."
MYTHIC_URL="https://localhost:7443"
MAX_ATTEMPTS=60
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    if curl -k -s "$MYTHIC_URL" > /dev/null 2>&1; then
        echo "✓ Mythic is ready"
        break
    fi
    echo "Waiting for Mythic... ($ATTEMPT/$MAX_ATTEMPTS)"
    sleep 5
    ATTEMPT=$((ATTEMPT + 1))
done

if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
    echo "✗ Mythic failed to start within timeout"
    docker-compose logs
    docker-compose down -v
    exit 1
fi

# Run integration tests
echo ""
echo "Running integration tests..."
cd "$PROJECT_ROOT"

export MYTHIC_URL="https://localhost:7443"
export MYTHIC_USERNAME="mythic_admin"
export MYTHIC_PASSWORD="mythic_password"
export MYTHIC_SKIP_TLS_VERIFY="true"

if go test -v -tags=integration ./tests/integration/...; then
    echo ""
    echo "✓ Integration tests passed"
    TEST_RESULT=0
else
    echo ""
    echo "✗ Integration tests failed"
    TEST_RESULT=1
fi

# Cleanup
echo ""
echo "Stopping Mythic test instance..."
cd "$INTEGRATION_DIR"
docker-compose down -v

exit $TEST_RESULT
