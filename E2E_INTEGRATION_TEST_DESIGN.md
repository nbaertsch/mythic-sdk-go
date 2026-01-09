# End-to-End Integration Test Design

**Author:** Claude Code
**Date:** 2026-01-09
**Purpose:** Design comprehensive E2E integration tests for Mythic Go SDK
**Target:** 100% functional API coverage through realistic workflows

---

## Table of Contents
1. [Implementation Progress](#implementation-progress)
2. [Executive Summary](#executive-summary)
3. [Current State Analysis](#current-state-analysis)
4. [E2E Test Philosophy](#e2e-test-philosophy)
5. [Agent Selection](#agent-selection)
6. [Test Infrastructure](#test-infrastructure)
7. [E2E Test Workflows](#e2e-test-workflows)
8. [Implementation Plan](#implementation-plan)
9. [Tests to Remove](#tests-to-remove)
10. [Success Criteria](#success-criteria)

---

## Implementation Progress

**Status:** IN PROGRESS
**Started:** 2026-01-09
**Current Phase:** Phase 2 - Core Workflows
**Completion:** 2/16 workflows (12.5%)
**API Coverage:** 18/204 methods (8.8%)

### Completed Workflows

#### Phase 1: Infrastructure ✅ COMPLETE
- [CI Run #20867791754](https://github.com/nbaertsch/mythic-sdk-go/actions/runs/20867791754)
- **File:** `tests/integration/e2e_helpers.go`
- **E2ETestSetup struct** with agent management, callback waiting, task execution helpers
- **Cleanup functions** for full resource teardown
- **Status:** All infrastructure ready for E2E tests

#### ✅ Workflow 1: Authentication & Session Management - COMPLETE
- [CI Run #20867791754](https://github.com/nbaertsch/mythic-sdk-go/actions/runs/20867791754)
- **File:** `tests/integration/e2e_auth_test.go`
- **Tests:**
  * `TestE2E_AuthenticationLifecycle` - Complete auth workflow (15 sub-tests)
  * `TestE2E_AuthenticationErrorHandling` - Error scenarios (3 sub-tests)
- **APIs Covered (7/204):**
  * Login ✓
  * Logout ✓
  * IsAuthenticated ✓
  * GetMe ✓
  * CreateAPIToken ✓
  * DeleteAPIToken ✓
  * RefreshAccessToken ✓
- **Duration:** ~45 seconds
- **Skip Rate:** 0%
- **Status:** All tests passing in CI

#### ✅ Workflow 2: Operations & Global Settings - COMPLETE
- [CI Run #20868077719](https://github.com/nbaertsch/mythic-sdk-go/actions/runs/20868077719)
- **File:** `tests/integration/e2e_operations_test.go`
- **Tests:**
  * `TestE2E_OperationsManagement` - Complete operations workflow (16 sub-tests)
  * `TestE2E_OperationsErrorHandling` - Error scenarios (4 sub-tests)
- **APIs Covered (11/204):**
  * GetOperations ✓
  * GetOperationByID ✓
  * CreateOperation ✓
  * UpdateOperation ✓
  * SetCurrentOperation ✓
  * GetCurrentOperation ✓
  * GetOperatorsByOperation ✓
  * CreateOperationEventLog ✓
  * GetOperationEventLog ✓
  * GetGlobalSettings ✓
  * UpdateGlobalSettings ✓
- **Duration:** ~50 seconds
- **Skip Rate:** 0%
- **Status:** All tests passing in CI

### In Progress

None

### Pending Workflows
- ⏳ Workflow 3: File Management
- ⏳ Workflow 4: Credentials & Artifacts
- ⏳ Workflow 5: Tags & Categorization
- ⏳ Workflow 6: Operator & User Management
- ⏳ Workflow 7: C2 Profile Management
- ⏳ Workflow 8: MITRE ATT&CK Framework
- ⏳ Workflow 9: Payload Build & Deployment (Requires Poseidon)
- ⏳ Workflow 10: Callback & Task Execution (Requires Workflow 9)
- ⏳ Workflow 11: Advanced Features (Requires Workflow 10)
- ⏳ Workflow 12: Real-time Monitoring (Requires Workflow 10)
- ⏳ Workflow 13: P2P Networking (Optional, 2 agents)
- ⏳ Workflow 14: Eventing & Workflows
- ⏳ Workflow 15: Container Operations
- ⏳ Workflow 16: Utility Functions

---

## Executive Summary

### Problem Statement
Current integration tests contain **208 skip statements** and test APIs in isolation without realistic workflows. Tests skip when dependencies (callbacks, payloads, tasks) don't exist, making coverage unreliable.

### Solution
Implement **12 comprehensive E2E test workflows** that:
- Install and deploy Poseidon agent (Go-based, 5s compile time)
- Create full workflows from payload build → callback → task execution → cleanup
- Cover all 204 SDK methods through realistic operational scenarios
- Run reliably in CI without manual setup

### Coverage Impact
- **Before:** 415 test functions, ~50% skip in clean environment
- **After:** ~80-100 test functions, 0% skip, 100% API coverage through E2E workflows
- **Time:** ~5-10 minutes per full E2E test suite

---

## Current State Analysis

### SDK API Coverage (204 Methods)
Based on comprehensive analysis:

**Categories:**
- Authentication & Session (7 methods)
- Operations Management (11 methods)
- Callbacks (14 methods)
- Payloads (12 methods)
- Tasks & Responses (12 methods)
- Files & Downloads (10 methods)
- Credentials & Artifacts (10 methods)
- C2 Profiles (8 methods)
- MITRE ATT&CK (6 methods)
- Operators & Users (11 methods)
- Tags & Categories (9 methods)
- Advanced Features (94 remaining methods across 20+ categories)

### Current Test Problems

**1. Tests That Always Skip in Clean Environment:**
```go
// Callback tests skip when no callbacks exist
if len(callbacks) == 0 {
    t.Skip("No callbacks available")
}

// Task tests skip when no active callbacks
if len(activeCallbacks) == 0 {
    t.Skip("No active callbacks for task testing")
}

// Process tests skip when no processes enumerated
if len(processes) == 0 {
    t.Skip("No processes found")
}
```

**2. Tests Requiring Manual Setup:**
- Payload container deployment
- Agent installation on target hosts
- C2 profile configuration
- Callback establishment
- Command execution

**3. Incomplete Workflows:**
- Test payload creation but not deployment
- Test task creation but not execution
- Test callback updates but not lifecycle
- Test file operations but not agent downloads

---

## E2E Test Philosophy

### Principles

**1. No Manual Setup**
Tests must:
- Install required infrastructure automatically
- Deploy agents programmatically
- Create all dependencies in test setup
- Clean up completely in teardown

**2. Realistic Workflows**
Tests should mirror actual operator workflows:
- Build payload → Deploy → Get callback → Execute tasks → Collect results
- Not: Test each API in isolation

**3. Fail Fast**
Tests must FAIL (not skip) when:
- Mythic server unavailable
- Agent containers not running
- Infrastructure setup fails
- API calls fail unexpectedly

**4. Comprehensive Coverage**
Every SDK method must be tested in at least one E2E workflow.

**5. Reproducible**
Tests must produce identical results on every run in CI.

---

## Agent Selection

### Chosen Agent: **Poseidon** (Golang)

**Rationale:**
1. **Fast Compile:** ~5 seconds build time
2. **Cross-platform:** Linux & macOS support
3. **Rich Feature Set:** All major capabilities (shell, upload, download, process, etc.)
4. **Active Development:** Well-maintained by MythicAgents
5. **Docker Support:** Official container available

**Installation in Tests:**
```bash
# Install Poseidon agent
sudo ./mythic-cli install github https://github.com/MythicAgents/poseidon

# Wait for agent container
sudo ./mythic-cli status
```

**Capabilities Covered:**
- Shell command execution (for tasks)
- File upload/download (for file operations)
- Process enumeration (for process tests)
- Screenshot capture (for screenshot tests)
- Keylogging (for keylog tests)
- Token enumeration (for token tests)
- Network connections (for host/network tests)

**Sources:**
- [MythicAgents/poseidon on GitHub](https://github.com/MythicAgents/poseidon)
- [Introduction to Mythic C2 - RedSiege](https://redsiege.com/blog/2023/06/introduction-to-mythic-c2/)
- [Mythic C2 Framework Documentation](https://docs.mythic-c2.net/home)

---

## Test Infrastructure

### Test Environment Setup

**CI Workflow Enhancement (.github/workflows/integration.yml):**
```yaml
- name: Install Poseidon Agent
  run: |
    cd /tmp/mythic
    sudo ./mythic-cli install github https://github.com/MythicAgents/poseidon
    sudo ./mythic-cli status

- name: Wait for Poseidon Ready
  run: |
    echo "Waiting for Poseidon container..."
    timeout=120
    attempt=0
    while [ $attempt -lt $timeout ]; do
      if sudo docker ps | grep -q poseidon; then
        echo "✓ Poseidon ready after $attempt seconds"
        break
      fi
      sleep 5
      attempt=$((attempt + 5))
    done
```

### Test Helpers (tests/integration/e2e_helpers.go)

**New helper functions:**
```go
// E2ETestSetup performs full environment setup for E2E tests
type E2ETestSetup struct {
    Client        *mythic.Client
    OperationID   int
    PayloadUUID   string
    CallbackID    int
    AgentProcess  *os.Process
}

func SetupE2ETest(t *testing.T) *E2ETestSetup
func (s *E2ETestSetup) Cleanup()
func (s *E2ETestSetup) WaitForCallback(timeout time.Duration) error
func (s *E2ETestSetup) ExecuteCommand(cmd string) (*types.Task, error)
func (s *E2ETestSetup) WaitForTaskOutput(taskID int) (string, error)
```

### Payload Build Configuration

**Default Poseidon Build Config:**
```go
payloadConfig := &types.CreatePayloadRequest{
    PayloadType: "poseidon",
    C2Profiles: []types.C2ProfileSelection{
        {
            Name: "http",
            Parameters: map[string]interface{}{
                "callback_host":     "http://127.0.0.1:80",
                "callback_interval": 5,
                "callback_jitter":   10,
            },
        },
    },
    BuildParameters: map[string]interface{}{
        "mode":         "default",
        "architecture": "AMD_x64",
        "output_type":  "Executable",
    },
}
```

### Agent Execution

**Start agent in background:**
```go
func (s *E2ETestSetup) StartAgent() error {
    payloadPath := fmt.Sprintf("/tmp/payload_%s", s.PayloadUUID)

    cmd := exec.Command(payloadPath)
    err := cmd.Start()
    if err != nil {
        return fmt.Errorf("failed to start agent: %w", err)
    }

    s.AgentProcess = cmd.Process
    return nil
}
```

---

## E2E Test Workflows

### Workflow 1: Core Authentication & Session Management
**File:** `tests/integration/e2e_auth_test.go`
**Duration:** ~30 seconds
**Dependencies:** None (Mythic only)

**Test: TestE2E_AuthenticationLifecycle**
```
Setup:
  - None required

Workflow:
  1. Login with username/password (Login)
  2. Get current user info (GetMe)
  3. Verify authenticated (IsAuthenticated)
  4. Get current operation (GetCurrentOperation)
  5. Create API token (CreateAPIToken)
  6. Create new client with API token
  7. Login with API token
  8. Verify new client authenticated (IsAuthenticated)
  9. Call authenticated endpoint with token client (GetMe)
  10. Refresh access token (RefreshAccessToken)
  11. Verify refresh worked (GetMe)
  12. Delete API token (DeleteAPIToken)
  13. Logout both clients (Logout)
  14. Verify not authenticated (IsAuthenticated)

Cleanup:
  - Logout all clients

APIs Covered (7):
  - Login, Logout, IsAuthenticated, GetMe, CreateAPIToken,
    DeleteAPIToken, RefreshAccessToken

Assertions:
  - Login succeeds
  - GetMe returns valid operator
  - API token login works
  - Refresh token works
  - Logout clears auth
```

---

### Workflow 2: Operations & Global Settings
**File:** `tests/integration/e2e_operations_test.go`
**Duration:** ~45 seconds
**Dependencies:** None

**Test: TestE2E_OperationsManagement**
```
Setup:
  - Authenticate

Workflow:
  1. Get all operations (GetOperations)
  2. Get current operation (GetCurrentOperation)
  3. Get operation by ID (GetOperationByID)
  4. Create new operation (CreateOperation)
  5. Switch to new operation (SetCurrentOperation)
  6. Update operation settings (UpdateOperation)
  7. Get operators in operation (GetOperatorsByOperation)
  8. Create event log entry (CreateOperationEventLog)
  9. Get event log (GetOperationEventLog)
  10. Get global settings (GetGlobalSettings)
  11. Update global settings (UpdateGlobalSettings)
  12. Verify settings updated (GetGlobalSettings)
  13. Restore original settings
  14. Switch back to original operation
  15. Delete test operation (DeleteOperation if exists)

Cleanup:
  - Restore global settings
  - Switch to original operation

APIs Covered (11):
  - GetOperations, GetOperationByID, CreateOperation,
    UpdateOperation, SetCurrentOperation, GetCurrentOperation,
    GetOperatorsByOperation, CreateOperationEventLog,
    GetOperationEventLog, GetGlobalSettings, UpdateGlobalSettings

Assertions:
  - Operations list is not empty
  - New operation created successfully
  - Can switch operations
  - Settings persist after update
```

---

### Workflow 3: File Management
**File:** `tests/integration/e2e_files_test.go`
**Duration:** ~40 seconds
**Dependencies:** None

**Test: TestE2E_FileOperations**
```
Setup:
  - Authenticate
  - Create test files (1KB, 1MB, 10MB)

Workflow:
  1. Get files before upload (GetFiles) - baseline
  2. Upload small file (UploadFile)
  3. Upload medium file (UploadFile)
  4. Upload large file (UploadFile)
  5. Get all files (GetFiles) - verify 3 new files
  6. Get file by ID (GetFileByID) for each
  7. Verify file metadata (size, hash, timestamp)
  8. Preview small file (PreviewFile)
  9. Download small file (DownloadFile)
  10. Verify downloaded content matches original
  11. Download medium file (DownloadFile)
  12. Bulk download all 3 files (BulkDownloadFiles)
  13. Verify bulk download contents
  14. Get downloaded files filter (GetDownloadedFiles)
  15. Delete files one by one (DeleteFile)
  16. Verify files marked deleted (GetFiles)

Cleanup:
  - Delete all uploaded files
  - Remove temp files

APIs Covered (8):
  - GetFiles, GetFileByID, GetDownloadedFiles, UploadFile,
    DownloadFile, DeleteFile, BulkDownloadFiles, PreviewFile

Assertions:
  - Upload succeeds for various sizes
  - Download matches upload
  - Bulk download works
  - Delete marks files deleted
```

---

### Workflow 4: Credentials & Artifacts
**File:** `tests/integration/e2e_credentials_artifacts_test.go`
**Duration:** ~35 seconds
**Dependencies:** None

**Test: TestE2E_CredentialManagement**
```
Setup:
  - Authenticate

Workflow:
  1. Get credentials baseline (GetCredentials)
  2. Create credential - password (CreateCredential)
  3. Create credential - SSH key (CreateCredential)
  4. Create credential - API token (CreateCredential)
  5. Get all credentials (GetCredentials) - verify 3 new
  6. Get credentials by operation (GetCredentialsByOperation)
  7. Update credential (UpdateCredential) - change comment
  8. Verify update (GetCredentialsByOperation)
  9. Delete credentials (DeleteCredential)
  10. Verify deletion (GetCredentials)

Cleanup:
  - Delete test credentials

APIs Covered (5):
  - GetCredentials, GetCredentialsByOperation, CreateCredential,
    UpdateCredential, DeleteCredential

Assertions:
  - Credentials created with correct types
  - Can query by operation
  - Updates persist
  - Deletion works
```

**Test: TestE2E_ArtifactManagement**
```
Setup:
  - Authenticate

Workflow:
  1. Get artifacts baseline (GetArtifacts)
  2. Create file artifact (CreateArtifact)
  3. Create registry artifact (CreateArtifact)
  4. Create process artifact (CreateArtifact)
  5. Get all artifacts (GetArtifacts)
  6. Get artifacts by operation (GetArtifactsByOperation)
  7. Get artifacts by host (GetArtifactsByHost)
  8. Get artifacts by type (GetArtifactsByType)
  9. Update artifact (UpdateArtifact) - mark reviewed
  10. Verify update (GetArtifacts)
  11. Delete artifacts (DeleteArtifact)
  12. Verify deletion (GetArtifacts)

Cleanup:
  - Delete test artifacts

APIs Covered (7):
  - GetArtifacts, GetArtifactsByOperation, GetArtifactsByHost,
    GetArtifactsByType, CreateArtifact, UpdateArtifact, DeleteArtifact

Assertions:
  - Artifacts created with metadata
  - Filtering works correctly
  - Updates persist
```

---

### Workflow 5: Tags & Categorization
**File:** `tests/integration/e2e_tags_test.go`
**Duration:** ~30 seconds
**Dependencies:** None

**Test: TestE2E_TagManagement**
```
Setup:
  - Authenticate

Workflow:
  1. Get tag types baseline (GetTagTypes)
  2. Create tag type - "Priority" (CreateTagType)
  3. Create tag type - "Target Type" (CreateTagType)
  4. Get all tag types (GetTagTypes)
  5. Get tag types by operation (GetTagTypesByOperation)
  6. Update tag type (UpdateTagType) - add description
  7. Create tags with types (CreateTag)
     - Priority: High, Medium, Low
     - Target Type: Domain Controller, Workstation, Server
  8. Get all tags (GetTags)
  9. Get tags by operation (GetTagsByOperation)
  10. Assign tags to resources (when available in later workflows)
  11. Delete tags (DeleteTag)
  12. Delete tag types (DeleteTagType)

Cleanup:
  - Delete test tags and tag types

APIs Covered (8):
  - CreateTag, GetTags, GetTagsByOperation, DeleteTag,
    CreateTagType, GetTagTypes, GetTagTypesByOperation, UpdateTagType

Assertions:
  - Tag types created successfully
  - Tags associated with types
  - Filtering by operation works
```

---

### Workflow 6: Operator & User Management
**File:** `tests/integration/e2e_operators_test.go`
**Duration:** ~40 seconds
**Dependencies:** Admin privileges

**Test: TestE2E_OperatorManagement**
```
Setup:
  - Authenticate as admin
  - Save original operator preferences

Workflow:
  1. Get all operators (GetOperators)
  2. Get current operator (GetOperatorByID)
  3. Get operator preferences (GetOperatorPreferences)
  4. Update preferences (UpdateOperatorPreferences)
     - Change theme, timezone, etc.
  5. Verify preferences updated (GetOperatorPreferences)
  6. Get operator secrets (GetOperatorSecrets)
  7. Update secrets (UpdateOperatorSecrets)
  8. Verify secrets updated (GetOperatorSecrets)
  9. Create invite link (CreateInviteLink)
  10. Get invite links (GetInviteLinks)
  11. Create new operator (CreateOperator)
  12. Get operators (GetOperators) - verify new operator
  13. Update operator status (UpdateOperatorStatus) - deactivate
  14. Verify status change (GetOperatorByID)
  15. Update operator status (UpdateOperatorStatus) - reactivate
  16. Update operator operation (UpdateOperatorOperation) - add to current op
  17. Verify operation assignment (GetOperatorsByOperation)
  18. Update password/email (UpdatePasswordAndEmail)
  19. Restore original preferences

Cleanup:
  - Restore operator preferences
  - Deactivate test operator
  - Restore password if changed

APIs Covered (11):
  - GetOperators, GetOperatorByID, CreateOperator,
    UpdateOperatorStatus, UpdateOperatorOperation,
    GetOperatorPreferences, UpdateOperatorPreferences,
    GetOperatorSecrets, UpdateOperatorSecrets,
    CreateInviteLink, GetInviteLinks, UpdatePasswordAndEmail

Assertions:
  - Operator creation works
  - Preferences persist
  - Status changes work
  - Operation assignment works
```

---

### Workflow 7: C2 Profile Management
**File:** `tests/integration/e2e_c2profiles_test.go`
**Duration:** ~50 seconds
**Dependencies:** C2 containers running

**Test: TestE2E_C2ProfileOperations**
```
Setup:
  - Authenticate
  - Verify HTTP C2 profile available

Workflow:
  1. Get all C2 profiles (GetC2Profiles)
  2. Verify HTTP profile exists
  3. Get HTTP profile by ID (GetC2ProfileByID)
  4. Get profile parameters
  5. Create C2 instance for testing (CreateC2Instance)
  6. Get profile output/logs (GetProfileOutput)
  7. Test sample message (C2SampleMessage)
  8. Get profile IOCs (C2GetIOC)
  9. Export profile config (if supported)
  10. Import profile config (ImportC2Instance)
  11. Start/stop profile (StartStopProfile)
  12. Verify profile state changes

Cleanup:
  - Restore original profile state
  - Remove test C2 instance

APIs Covered (8):
  - GetC2Profiles, GetC2ProfileByID, CreateC2Instance,
    ImportC2Instance, StartStopProfile, GetProfileOutput,
    C2HostFile, C2SampleMessage, C2GetIOC

Assertions:
  - C2 profiles available
  - Can create instances
  - IOC generation works
  - State management works
```

---

### Workflow 8: MITRE ATT&CK Framework
**File:** `tests/integration/e2e_attack_test.go`
**Duration:** ~25 seconds
**Dependencies:** None (uses built-in mappings)

**Test: TestE2E_MITREAttackMapping**
```
Setup:
  - Authenticate

Workflow:
  1. Get all MITRE techniques (GetAttackTechniques)
  2. Verify techniques loaded (should be ~600+)
  3. Get technique by ID (GetAttackTechniqueByID)
  4. Get technique by T-number (GetAttackTechniqueByTNum) - "T1059"
  5. Verify technique details (name, description, tactics)
  6. Get techniques by tactic (filter)
  7. Search for command execution techniques
  8. Get command mappings (GetAttackByCommand)
  9. Verify command->technique mappings exist
  10. Get operation attack statistics (GetAttacksByOperation)
  11. When task available (from Workflow 10):
      - Add MITRE mapping to task (AddMITREAttackToTask)
      - Get attack by task (GetAttackByTask)
      - Verify mapping persisted

Cleanup:
  - None (read-only except task mapping)

APIs Covered (6):
  - GetAttackTechniques, GetAttackTechniqueByID,
    GetAttackTechniqueByTNum, GetAttackByCommand, GetAttackByTask,
    AddMITREAttackToTask, GetAttacksByOperation

Assertions:
  - Techniques loaded in database
  - Can query by ID and T-number
  - Command mappings exist
  - Can add custom mappings
```

---

### Workflow 9: Payload Build & Deployment (Poseidon)
**File:** `tests/integration/e2e_payload_test.go`
**Duration:** ~2 minutes (includes build time)
**Dependencies:** Poseidon container running

**Test: TestE2E_PayloadLifecycle**
```
Setup:
  - Authenticate
  - Verify Poseidon payload type available
  - Verify HTTP C2 profile available

Workflow:
  1. Get all payload types (GetPayloadTypes)
  2. Verify "poseidon" in list
  3. Get build parameters for Poseidon (GetBuildParametersByPayloadType)
  4. Get available commands (GetPayloadCommands)
  5. Create payload build config
     - Type: poseidon
     - OS: linux (for CI environment)
     - C2: http (localhost callback)
     - Commands: shell, download, upload, ps, etc.
  6. Submit payload build (CreatePayload)
  7. Wait for build completion (WaitForPayloadComplete) - up to 60s
  8. Verify build succeeded
  9. Get payload by UUID (GetPayloadByUUID)
  10. Verify payload metadata (commands, c2, status)
  11. Download payload binary (DownloadPayload)
  12. Verify binary downloaded (file exists, size > 0)
  13. Save payload to /tmp for agent execution
  14. Export payload config (ExportPayloadConfig)
  15. Verify config JSON valid
  16. Update payload (UpdatePayload) - add description
  17. Get payload info (GetPayloadOnHost) - empty initially
  18. Create second payload for rebuild testing
  19. Rebuild payload (RebuildPayload)
  20. Wait for rebuild (WaitForPayloadComplete)
  21. Delete second payload (DeletePayload)
  22. Verify deletion (GetPayloads)

Cleanup:
  - Delete test payloads
  - Remove temp payload binaries

APIs Covered (12):
  - GetPayloadTypes, GetBuildParameters,
    GetBuildParametersByPayloadType, GetBuildParameterInstances,
    GetBuildParameterInstancesByPayload, CreatePayload, UpdatePayload,
    DeletePayload, GetPayloadByUUID, GetPayloads, DownloadPayload,
    ExportPayloadConfig, RebuildPayload, WaitForPayloadComplete,
    GetPayloadCommands, GetPayloadOnHost

Assertions:
  - Poseidon build succeeds in <60s
  - Binary is valid executable
  - Config export is valid JSON
  - Rebuild works
  - Deletion works
```

---

### Workflow 10: Agent Callback & Command Execution (FULL E2E)
**File:** `tests/integration/e2e_callback_task_test.go`
**Duration:** ~3-5 minutes
**Dependencies:** Workflow 9 (payload built)

**Test: TestE2E_CallbackTaskLifecycle**
```
Setup:
  - Run Workflow 9 to get payload
  - Make payload executable (chmod +x)

Workflow:
  Part 1: Callback Establishment
  1. Get callbacks baseline (GetAllCallbacks)
  2. Start Poseidon agent in background
  3. Wait for callback (up to 60s)
     - Poll GetAllActiveCallbacks every 5s
  4. Verify callback created
  5. Get callback by ID (GetCallbackByID)
  6. Verify callback details (IP, hostname, user, process)
  7. Update callback (UpdateCallback) - add description
  8. Verify update persisted (GetCallbackByID)
  9. Export callback config (ExportCallbackConfig)
  10. Get loaded commands (GetLoadedCommands)

  Part 2: Task Execution - Shell Command
  11. Issue shell task (IssueTask) - "whoami"
  12. Get task immediately (GetTask)
  13. Verify task created (status: submitted)
  14. Wait for completion (WaitForTaskComplete) - up to 30s
  15. Get task again (GetTask) - verify completed
  16. Get task output (GetTaskOutput)
  17. Verify output contains expected result
  18. Get responses by task (GetResponsesByTask)
  19. Verify response data
  20. Get latest responses (GetLatestResponses)
  21. Search responses (SearchResponses) - search for "whoami"
  22. Get response by ID (GetResponseByID)
  23. Get response statistics (GetResponseStatistics)

  Part 3: Task Execution - File Upload
  24. Create test file locally
  25. Issue upload task (IssueTask) - upload test file
  26. Wait for completion (WaitForTaskComplete)
  27. Get task artifacts (GetTaskArtifacts)
  28. Verify file artifact created
  29. Get files (GetFiles)
  30. Find uploaded file in list
  31. Download file (DownloadFile)
  32. Verify downloaded = original

  Part 4: Task Execution - File Download
  33. Issue download task (IssueTask) - download /etc/hostname
  34. Wait for completion (WaitForTaskComplete)
  35. Get files (GetFiles)
  36. Find downloaded file
  37. Download from Mythic (DownloadFile)
  38. Verify file content

  Part 5: Task Execution - Process List
  39. Issue ps task (IssueTask) - list processes
  40. Wait for completion (WaitForTaskComplete)
  41. Get processes (GetProcesses)
  42. Verify process list populated
  43. Get process tree (GetProcessTree) for callback
  44. Get processes by callback (GetProcessesByCallback)
  45. Get processes by operation (GetProcessesByOperation)

  Part 6: Host Enumeration
  46. Get hosts (GetHosts)
  47. Find callback's host
  48. Get host by ID (GetHostByID)
  49. Get host by hostname (GetHostByHostname)
  50. Get callbacks for host (GetCallbacksForHost)
  51. Get host network map (GetHostNetworkMap)
  52. Get processes by host (GetProcessesByHost)

  Part 7: Task Management
  53. Issue failing task (IssueTask) - invalid command
  54. Wait for completion (WaitForTaskComplete)
  55. Verify task error
  56. Update task (UpdateTask) - add comment
  57. Reissue task (ReissueTask)
  58. Wait for completion (WaitForTaskComplete)
  59. Test reissue with handler (ReissueTaskWithHandler)
  60. Get tasks by status (GetTasksByStatus) - completed/error
  61. Get tasks for callback (GetTasksForCallback)

  Part 8: Callback Cleanup
  62. Get callback tokens (GetCallbackTokens)
  63. Get callback tokens by callback (GetCallbackTokensByCallback)
  64. Delete callback (DeleteCallback)
  65. Verify deletion (GetAllCallbacks)
  66. Kill agent process

Cleanup:
  - Kill agent process
  - Delete callback
  - Delete test files
  - Delete payload

APIs Covered (45):
  Callbacks: GetAllCallbacks, GetAllActiveCallbacks, GetCallbackByID,
             UpdateCallback, DeleteCallback, ExportCallbackConfig,
             GetCallbacksForHost, GetLoadedCommands, CreateCallback (manual)
  Tasks: IssueTask, GetTask, UpdateTask, GetTasksByStatus,
         GetTasksForCallback, WaitForTaskComplete, ReissueTask,
         ReissueTaskWithHandler, GetTaskArtifacts, GetTaskOutput
  Responses: GetResponsesByTask, GetResponsesByCallback, GetResponseByID,
             GetLatestResponses, SearchResponses, GetResponseStatistics
  Files: UploadFile, DownloadFile, GetFiles, GetFileByID
  Processes: GetProcesses, GetProcessTree, GetProcessesByCallback,
             GetProcessesByOperation, GetProcessesByHost
  Hosts: GetHosts, GetHostByID, GetHostByHostname,
         GetCallbacksForHost, GetHostNetworkMap
  Tokens: GetCallbackTokens, GetCallbackTokensByCallback

Assertions:
  - Callback established within timeout
  - Tasks execute successfully
  - File upload/download works
  - Process enumeration works
  - Host enumeration works
  - Task reissue works
  - Cleanup successful
```

---

### Workflow 11: Advanced Features (Screenshots, Keylogs, Tokens)
**File:** `tests/integration/e2e_advanced_test.go`
**Duration:** ~2 minutes
**Dependencies:** Workflow 10 (active callback)

**Test: TestE2E_AdvancedFeatures**
```
Setup:
  - Run Workflow 10 to get active callback
  - Ensure callback has advanced commands loaded

Workflow:
  Part 1: Screenshots
  1. Issue screenshot task (IssueTask)
  2. Wait for completion (WaitForTaskComplete)
  3. Get screenshots (GetScreenshots)
  4. Verify screenshot created
  5. Get screenshot by ID (GetScreenshotByID)
  6. Get screenshot thumbnail (GetScreenshotThumbnail)
  7. Get screenshot timeline (GetScreenshotTimeline)
  8. Download screenshot (DownloadScreenshot)
  9. Verify screenshot file valid
  10. Delete screenshot (DeleteScreenshot)

  Part 2: Keylogs
  11. Issue keylog start task (IssueTask)
  12. Simulate keyboard activity
  13. Wait for keylog data
  14. Get keylogs (GetKeylogs)
  15. Get keylogs by callback (GetKeylogsByCallback)
  16. Get keylogs by operation (GetKeylogsByOperation)
  17. Verify keylog data captured
  18. Issue keylog stop task (IssueTask)

  Part 3: Token Enumeration
  19. Issue token enumeration task (IssueTask)
  20. Wait for completion (WaitForTaskComplete)
  21. Get tokens (GetTokens)
  22. Get tokens by operation (GetTokensByOperation)
  23. Get token by ID (GetTokenByID)
  24. Verify token details

  Part 4: File Browser
  25. Issue file browser task (IssueTask) - browse /tmp
  26. Wait for completion (WaitForTaskComplete)
  27. Get file browser objects (GetFileBrowserObjects)
  28. Get file browser by callback (GetFileBrowserObjectsByCallback)
  29. Get file browser by host (GetFileBrowserObjectsByHost)
  30. Verify directory listing

  Part 5: Network Proxy/RPFWD
  31. Create RPFWD (CreateRPFWD) - forward port 8080
  32. Get RFPWDs (GetRPFWDs)
  33. Get RPFWD status (GetRPFWDStatus)
  34. Test proxy (TestProxy)
  35. Toggle proxy (ToggleProxy)
  36. Delete RPFWD (DeleteRPFWD)

Cleanup:
  - Stop keylogs
  - Delete screenshots
  - Delete RPFWD
  - Clean up test data

APIs Covered (23):
  Screenshots: GetScreenshots, GetScreenshotByID, GetScreenshotThumbnail,
               GetScreenshotTimeline, DownloadScreenshot, DeleteScreenshot
  Keylogs: GetKeylogs, GetKeylogsByCallback, GetKeylogsByOperation
  Tokens: GetTokens, GetTokensByOperation, GetTokenByID
  File Browser: GetFileBrowserObjects, GetFileBrowserObjectsByCallback,
                GetFileBrowserObjectsByHost
  Network: CreateRPFWD, GetRPFWDs, GetRPFWDStatus, DeleteRPFWD,
           TestProxy, ToggleProxy

Assertions:
  - Screenshots captured and downloadable
  - Keylogs record activity
  - Tokens enumerated correctly
  - File browser works
  - RPFWD creates forwarding rules
```

---

### Workflow 12: Alerts, Subscriptions & Real-time Updates
**File:** `tests/integration/e2e_realtime_test.go`
**Duration:** ~1-2 minutes
**Dependencies:** Workflow 10 (active callback for events)

**Test: TestE2E_RealtimeMonitoring**
```
Setup:
  - Run Workflow 10 to get active callback
  - Prepare to capture events

Workflow:
  Part 1: Subscription Setup
  1. Create WebSocket subscription (Subscribe)
     - Type: all events
     - Operation: current
  2. Verify subscription active
  3. Subscribe to alerts (SubscribeToAlerts)

  Part 2: Generate Events
  4. Issue task to trigger events (IssueTask) - shell command
  5. Wait for task output event in subscription
  6. Verify event received via WebSocket
  7. Verify event data matches task

  Part 3: Alert Management
  8. Create custom alert (CreateCustomAlert)
  9. Wait for alert in subscription
  10. Get unresolved alerts (GetUnresolvedAlerts)
  11. Verify alert in list
  12. Get alert by ID (GetAlertByID)
  13. Get all alerts (GetAlerts)
  14. Get alert statistics (GetAlertStatistics)
  15. Resolve alert (ResolveAlert)
  16. Verify alert resolved (GetUnresolvedAlerts)

  Part 4: Subscription Management
  17. Unsubscribe from alerts (Unsubscribe)
  18. Unsubscribe from events (Unsubscribe)
  19. Verify subscriptions closed

  Part 5: Event Reporting
  20. Generate comprehensive report (GenerateReport)
  21. Verify report contains test data
  22. Custom browser export (CustomBrowserExport)
  23. Get browser scripts (GetBrowserScripts)
  24. Get browser scripts by operation (GetBrowserScriptsByOperation)

Cleanup:
  - Close all subscriptions
  - Resolve test alerts
  - Clean up generated reports

APIs Covered (13):
  Subscriptions: Subscribe, Unsubscribe, SubscribeToAlerts
  Alerts: GetAlerts, GetAlertByID, GetUnresolvedAlerts,
          ResolveAlert, CreateCustomAlert, GetAlertStatistics
  Reporting: GenerateReport, GetRedirectRules, CustomBrowserExport
  Browser Scripts: GetBrowserScripts, GetBrowserScriptsByOperation

Assertions:
  - Subscriptions receive events in real-time
  - WebSocket connection stable
  - Alerts created and resolved
  - Events correlate to actions
  - Reports generated successfully
```

---

### Workflow 13: P2P Networking & Multi-Callback Operations
**File:** `tests/integration/e2e_p2p_test.go`
**Duration:** ~4 minutes
**Dependencies:** 2 active callbacks from Workflow 10

**Test: TestE2E_P2PNetworking** (Optional - requires 2 agents)
```
Setup:
  - Build 2 Poseidon payloads (different ports)
  - Start 2 agents
  - Wait for 2 callbacks

Workflow:
  1. Get all callbacks (GetAllCallbacks) - verify 2
  2. Select parent and child callbacks
  3. Add callback graph edge (AddCallbackGraphEdge)
     - Parent: callback1
     - Child: callback2
     - C2: http
  4. Verify edge created
  5. Issue task on parent
  6. Verify task proxied to child
  7. Get task output from child
  8. Issue task on child directly
  9. Compare routing behavior
  10. Remove callback graph edge (RemoveCallbackGraphEdge)
  11. Verify edge removed
  12. Test direct communication blocked

Cleanup:
  - Remove graph edges
  - Delete both callbacks
  - Kill both agents

APIs Covered (4):
  - AddCallbackGraphEdge, RemoveCallbackGraphEdge,
    GetAllCallbacks (filtering), IssueTask (P2P routing)

Assertions:
  - P2P connection established
  - Tasks route through parent
  - Edge removal works
  - Cleanup successful
```

---

### Workflow 14: Eventing & Workflow Automation
**File:** `tests/integration/e2e_eventing_test.go`
**Duration:** ~2 minutes
**Dependencies:** Workflow 10 (callback for triggers)

**Test: TestE2E_EventingWorkflows**
```
Setup:
  - Authenticate
  - Create test workflow config

Workflow:
  1. Export workflow template (EventingExportWorkflow)
  2. Modify workflow for testing
  3. Import workflow (EventingImportContainerWorkflow)
  4. Test workflow file (EventingTestFile)
  5. Trigger manual event (EventingTriggerManual)
  6. Wait for workflow execution
  7. Verify workflow completed
  8. Trigger bulk manual events (EventingTriggerManualBulk)
  9. Update event trigger (EventingTriggerUpdate)
  10. Handle event approval (UpdateEventGroupApproval)
  11. Test event retry (EventingTriggerRetry)
  12. Test retry from step (EventingTriggerRetryFromStep)
  13. Test run again (EventingTriggerRunAgain)
  14. Cancel event (EventingTriggerCancel)
  15. Test webhook trigger (SendExternalWebhook)
  16. Test consuming service webhook (ConsumingServicesTestWebhook)
  17. Test consuming service log (ConsumingServicesTestLog)
  18. Get event logs
  19. Verify all events tracked

Cleanup:
  - Delete test workflows
  - Clean up event logs

APIs Covered (13):
  - EventingExportWorkflow, EventingImportContainerWorkflow,
    EventingTestFile, EventingTriggerManual, EventingTriggerManualBulk,
    EventingTriggerUpdate, EventingTriggerRetry, EventingTriggerRetryFromStep,
    EventingTriggerRunAgain, EventingTriggerCancel,
    UpdateEventGroupApproval, SendExternalWebhook,
    ConsumingServicesTestWebhook, ConsumingServicesTestLog

Assertions:
  - Workflows import successfully
  - Manual triggers work
  - Event approval works
  - Retry logic works
  - Webhooks fire correctly
```

---

### Workflow 15: Container & Build Operations
**File:** `tests/integration/e2e_containers_test.go`
**Duration:** ~1 minute
**Dependencies:** Payload containers running

**Test: TestE2E_ContainerOperations**
```
Setup:
  - Authenticate
  - Verify Poseidon container running

Workflow:
  1. List container files (ContainerListFiles) - poseidon
  2. Verify expected files present
  3. Create test file locally
  4. Write file to container (ContainerWriteFile)
  5. List files again (ContainerListFiles)
  6. Verify new file present
  7. Download file from container (ContainerDownloadFile)
  8. Verify content matches
  9. Remove file from container (ContainerRemoveFile)
  10. List files final check (ContainerListFiles)
  11. Verify file removed

Cleanup:
  - Remove test files from container

APIs Covered (4):
  - ContainerListFiles, ContainerDownloadFile,
    ContainerWriteFile, ContainerRemoveFile

Assertions:
  - Container accessible
  - File operations work
  - Cleanup successful
```

---

### Workflow 16: Utility Functions & Edge Cases
**File:** `tests/integration/e2e_utility_test.go`
**Duration:** ~1 minute
**Dependencies:** None

**Test: TestE2E_UtilityFunctions**
```
Setup:
  - Authenticate

Workflow:
  1. Get client config (GetConfig)
  2. Verify config values
  3. Run config check (ConfigCheck)
  4. Verify check results
  5. Generate random data (CreateRandom)
  6. Verify randomness
  7. Test dynamic query (DynamicQueryFunction)
  8. Test build parameter query (DynamicBuildParameter)
  9. Test typed array parsing (TypedArrayParseFunction)
  10. Get staging info (GetStagingInfo)
  11. Verify staging data
  12. Close client connection (Close)
  13. Verify connection closed
  14. Reconnect and verify works

Cleanup:
  - Close client

APIs Covered (8):
  - GetConfig, ConfigCheck, CreateRandom,
    DynamicQueryFunction, DynamicBuildParameter,
    TypedArrayParseFunction, GetStagingInfo, Close

Assertions:
  - Config accessible
  - Utility functions work
  - Client lifecycle works
```

---

## Implementation Plan

### Phase 1: Infrastructure Setup (Week 1)
**Tasks:**
1. Update CI workflow with Poseidon installation
2. Create e2e_helpers.go with setup/teardown functions
3. Test Poseidon build in CI environment
4. Verify agent callback establishment
5. Document infrastructure requirements

**Files to Create:**
- `tests/integration/e2e_helpers.go`
- Updated `.github/workflows/integration.yml`

**Validation:**
- Poseidon builds in <10 seconds
- Agent establishes callback
- Cleanup works reliably

---

### Phase 2: Core Workflows (Week 2)
**Implement in order:**
1. Workflow 1: Authentication ✅ (no dependencies)
2. Workflow 2: Operations ✅ (no dependencies)
3. Workflow 3: Files ✅ (no dependencies)
4. Workflow 4: Credentials & Artifacts ✅ (no dependencies)
5. Workflow 5: Tags ✅ (no dependencies)
6. Workflow 6: Operators ✅ (no dependencies)
7. Workflow 7: C2 Profiles ✅ (no dependencies)
8. Workflow 8: MITRE ATT&CK ✅ (no dependencies)

**Validation:**
- All tests pass in CI
- 0% skip rate
- <5 minutes total execution
- Cleanup leaves no artifacts

---

### Phase 3: Agent-Dependent Workflows (Week 3)
**Implement in order:**
1. Workflow 9: Payload Build ⚠️ (requires Poseidon)
2. Workflow 10: Callback & Task Execution ⚠️⚠️ (requires Workflow 9)
3. Workflow 11: Advanced Features ⚠️⚠️ (requires Workflow 10)
4. Workflow 12: Real-time Monitoring ⚠️⚠️ (requires Workflow 10)

**Validation:**
- Agent builds and runs reliably
- Callbacks establish consistently
- Tasks execute without failures
- WebSocket subscriptions stable

---

### Phase 4: Complex Workflows (Week 4)
**Implement in order:**
1. Workflow 14: Eventing & Workflows ⚠️ (complex setup)
2. Workflow 15: Container Operations ⚠️ (requires containers)
3. Workflow 16: Utility Functions ✅ (no dependencies)
4. Workflow 13: P2P Networking ⚠️⚠️⚠️ (optional, 2 agents)

**Validation:**
- All workflows pass
- Complex scenarios work
- Multi-agent tests optional but working

---

### Phase 5: Cleanup & Documentation (Week 5)
**Tasks:**
1. Remove old integration tests with skip statements
2. Update integration test README
3. Document E2E test running procedures
4. Add troubleshooting guides
5. Create developer documentation

**Deliverables:**
- Updated test documentation
- E2E test architecture guide
- Troubleshooting playbook
- CI/CD integration guide

---

## Tests to Remove

### Tests with Skip Statements (Remove these entire tests):

**From callbacks_test.go:**
- `TestCallbacks_CreateCallback` - Has skip, replaced by Workflow 10
- `TestCallbacks_DeleteCallback` - Has skip, replaced by Workflow 10
- All tests that skip when no callbacks available

**From tasks_test.go:**
- Tests that skip when no active callbacks
- Tests that skip when command not loaded
- Replaced entirely by Workflow 10

**From payloads_test.go:**
- Tests that skip when payload containers unavailable
- Tests that skip on build failures
- Replaced by Workflow 9

**From files_test.go:**
- Tests that skip when no agent files available
- Replaced by Workflow 3 + Workflow 10

**From processes_test.go:**
- Tests that skip when no processes enumerated
- Replaced by Workflow 10

**From hosts_test.go:**
- Tests that skip when no hosts available
- Replaced by Workflow 10

**From screenshots_test.go:**
- Tests that skip when no screenshots available
- Replaced by Workflow 11

**From keylogs_test.go:**
- Tests that skip when no keylogs available
- Replaced by Workflow 11

**From tokens_test.go:**
- Tests that skip when no tokens enumerated
- Replaced by Workflow 11

**From c2profiles_test.go:**
- Tests that skip when profiles not available
- Replaced by Workflow 7

**From browserscripts_test.go:**
- Tests with undefined client variable
- Fixed in Workflow 12

**From eventing_test.go:**
- Tests that skip when workflows not configured
- Replaced by Workflow 14

**From subscriptions_test.go:**
- Tests that skip when no live events
- Replaced by Workflow 12

### Total Tests to Remove: ~150-200 test functions

### Tests to Keep (Read-only APIs that work on empty state):
- Authentication tests (already work)
- Operation listing tests (work on default operation)
- Operator listing tests (work with default admin user)
- Command/Payload type listing (read-only metadata)
- MITRE ATT&CK technique queries (built-in data)
- Global settings queries (configuration data)

---

## Success Criteria

### Quantitative Metrics:
1. **API Coverage:** 100% (204/204 methods tested in E2E workflows)
2. **Skip Rate:** 0% (no tests skip due to missing dependencies)
3. **Pass Rate:** 100% (all tests pass in clean CI environment)
4. **Execution Time:** <10 minutes for full E2E suite
5. **Lines of Test Code:** ~5,000-8,000 lines (replacing 14,688 lines)
6. **Test Functions:** ~80-100 E2E tests (replacing 415 unit-style tests)

### Qualitative Goals:
1. ✅ Tests mirror realistic operator workflows
2. ✅ No manual infrastructure setup required
3. ✅ Tests fail fast with clear error messages
4. ✅ Comprehensive logging for debugging
5. ✅ Reproducible results in CI and locally
6. ✅ Self-contained cleanup (no artifacts left)
7. ✅ Documentation clear for contributors
8. ✅ Easy to add new E2E workflows

### Definition of Done:
- [ ] All 16 E2E workflows implemented
- [ ] CI runs successfully with 100% pass rate
- [ ] Old tests with skip statements removed
- [ ] Documentation updated
- [ ] README has E2E test instructions
- [ ] Troubleshooting guide created
- [ ] Code reviewed and approved
- [ ] All tests pass on main branch

---

## Appendix A: File Structure

```
tests/integration/
├── e2e_helpers.go               # E2E test infrastructure
├── e2e_auth_test.go            # Workflow 1
├── e2e_operations_test.go      # Workflow 2
├── e2e_files_test.go           # Workflow 3
├── e2e_credentials_artifacts_test.go  # Workflow 4
├── e2e_tags_test.go            # Workflow 5
├── e2e_operators_test.go       # Workflow 6
├── e2e_c2profiles_test.go      # Workflow 7
├── e2e_attack_test.go          # Workflow 8
├── e2e_payload_test.go         # Workflow 9
├── e2e_callback_task_test.go   # Workflow 10 (MAIN)
├── e2e_advanced_test.go        # Workflow 11
├── e2e_realtime_test.go        # Workflow 12
├── e2e_p2p_test.go             # Workflow 13 (optional)
├── e2e_eventing_test.go        # Workflow 14
├── e2e_containers_test.go      # Workflow 15
├── e2e_utility_test.go         # Workflow 16
├── helpers.go                  # Keep for backward compatibility
└── README.md                   # Update with E2E instructions

Files to REMOVE after E2E implementation:
- All test files with extensive skip statements
- Tests that don't work in clean environment
```

---

## Appendix B: CI Workflow Changes

### Updated .github/workflows/integration.yml

```yaml
name: Integration Tests

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]
  workflow_dispatch:

jobs:
  e2e-tests:
    name: E2E Integration Tests with Poseidon
    runs-on: ubuntu-latest
    timeout-minutes: 20

    steps:
      - name: Checkout SDK
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache: true

      - name: Clone Mythic Framework
        run: |
          git clone --depth 1 https://github.com/its-a-feature/Mythic.git /tmp/mythic
          cd /tmp/mythic
          echo "Mythic version: $(cat VERSION)"

      - name: Build Mythic CLI
        run: |
          cd /tmp/mythic
          sudo make
          sudo chmod +x mythic-cli

      - name: Start Mythic
        run: |
          cd /tmp/mythic
          sudo ./mythic-cli start

      - name: Install Poseidon Agent
        run: |
          cd /tmp/mythic
          echo "Installing Poseidon agent..."
          sudo ./mythic-cli install github https://github.com/MythicAgents/poseidon
          echo "Waiting for Poseidon container..."
          timeout=120
          attempt=0
          while [ $attempt -lt $timeout ]; do
            if sudo docker ps | grep -q poseidon; then
              echo "✓ Poseidon ready after $attempt seconds"
              break
            fi
            sleep 5
            attempt=$((attempt + 5))
          done
          if [ $attempt -ge $timeout ]; then
            echo "✗ Poseidon failed to start"
            sudo ./mythic-cli status
            exit 1
          fi

      - name: Wait for Mythic Ready
        run: |
          echo "Waiting for Mythic to be ready..."
          timeout=180
          attempt=0
          while [ $attempt -lt $timeout ]; do
            if curl -k -s https://127.0.0.1:7443 > /dev/null 2>&1; then
              echo "✓ Mythic is ready after $attempt seconds"
              break
            fi
            sleep 5
            attempt=$((attempt + 5))
          done
          if [ $attempt -ge $timeout ]; then
            echo "✗ Mythic failed to start within timeout"
            cd /tmp/mythic
            sudo ./mythic-cli status
            exit 1
          fi

      - name: Extract Mythic Credentials
        id: mythic-creds
        run: |
          PASSWORD=$(grep MYTHIC_ADMIN_PASSWORD /tmp/mythic/.env | cut -d'=' -f2 | tr -d '"')
          echo "::add-mask::$PASSWORD"
          echo "password=$PASSWORD" >> $GITHUB_OUTPUT

      - name: Run E2E Integration Tests
        env:
          MYTHIC_URL: "https://127.0.0.1:7443"
          MYTHIC_USERNAME: "mythic_admin"
          MYTHIC_PASSWORD: ${{ steps.mythic-creds.outputs.password }}
          MYTHIC_SKIP_TLS_VERIFY: "true"
        run: |
          set -o pipefail
          go test -v -tags=integration ./tests/integration/... \
            -timeout 15m \
            2>&1 | tee integration-test-output.log

      - name: Upload Test Results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: integration-test-results
          path: integration-test-output.log
          retention-days: 7

      - name: Collect Mythic Logs on Failure
        if: failure()
        run: |
          cd /tmp/mythic
          sudo ./mythic-cli logs > mythic-logs.txt 2>&1 || true
          sudo docker ps -a > docker-containers.txt || true

      - name: Upload Mythic Logs
        if: failure()
        uses: actions/upload-artifact@v4
        with:
          name: mythic-failure-logs
          path: |
            /tmp/mythic/mythic-logs.txt
            /tmp/mythic/docker-containers.txt
          retention-days: 7

      - name: Cleanup Mythic
        if: always()
        run: |
          cd /tmp/mythic
          sudo ./mythic-cli stop || true
          sudo docker system prune -af || true
```

---

## Appendix C: References

### Mythic Documentation
- [Mythic C2 Framework](https://docs.mythic-c2.net/home)
- [Mythic GitHub Repository](https://github.com/its-a-feature/Mythic)

### Agent Documentation
- [MythicAgents Organization](https://github.com/MythicAgents)
- [Poseidon Agent](https://github.com/MythicAgents/poseidon)

### Articles & Guides
- [Introduction to Mythic C2 - RedSiege Blog](https://redsiege.com/blog/2023/06/introduction-to-mythic-c2/)
- [C2 and the Docker Dance - SpecterOps](https://posts.specterops.io/c2-and-the-docker-dance-mythic-3-0s-marvelous-microservice-moves-f6e6e91356e2)
- [A Change of Mythic Proportions - SpecterOps](https://specterops.io/blog/2020/08/13/a-change-of-mythic-proportions/)

---

**END OF DESIGN DOCUMENT**
