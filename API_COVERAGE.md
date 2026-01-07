# Mythic SDK Go - API Coverage Report

Generated: 2026-01-07

This document provides a comprehensive overview of all available Mythic APIs and their implementation status in the Go SDK.

## Legend

- ✅ **Tested**: Fully implemented with unit and integration tests
- 🚧 **In Progress**: Implementation started but not complete
- ⏳ **Pending**: Not yet implemented

## Coverage Summary

| Category | Implemented | In Progress | Pending | Total |
|----------|-------------|-------------|---------|-------|
| Authentication | 4 | 0 | 0 | 4 |
| Callbacks | 12 | 0 | 0 | 12 |
| Tasks | 12 | 0 | 0 | 12 |
| Files | 8 | 0 | 0 | 8 |
| Operations | 11 | 0 | 0 | 11 |
| Payloads | 12 | 0 | 0 | 12 |
| Credentials | 5 | 0 | 0 | 5 |
| C2 Profiles | 0 | 0 | 9 | 9 |
| Artifacts | 0 | 0 | 3 | 3 |
| Tags | 0 | 0 | 3 | 3 |
| Tokens | 0 | 0 | 4 | 4 |
| Processes | 6 | 0 | 0 | 6 |
| Keylogs | 3 | 0 | 0 | 3 |
| Browser Scripts | 0 | 0 | 3 | 3 |
| MITRE ATT&CK | 0 | 0 | 3 | 3 |
| Reporting | 0 | 0 | 2 | 2 |
| Eventing/Workflows | 0 | 0 | 15 | 15 |
| Operators | 0 | 0 | 11 | 11 |
| GraphQL Subscriptions | 0 | 0 | 1 | 1 |
| Advanced Features | 0 | 0 | 20 | 20 |
| **TOTAL** | **73** | **0** | **60** | **133** |

**Overall Coverage: 54.9%**

---

## 1. Authentication & Authorization

### ✅ Tested (4/4 - 100%)

- **Login()** - Authenticate with username/password
  - File: `pkg/mythic/auth.go:36`
  - Tests: `tests/integration/auth_test.go`

- **CreateAPIToken()** - Generate API token for current user
  - File: `pkg/mythic/auth.go:138`
  - Tests: `tests/integration/auth_test.go:46`

- **GetMe()** - Get current authenticated user info
  - File: `pkg/mythic/auth.go:167`
  - Tests: `tests/integration/auth_test.go:30`

- **RefreshAccessToken()** - Refresh JWT access token using refresh token
  - File: `pkg/mythic/auth.go:242`
  - Tests: `tests/integration/auth_test.go:248`

---

## 2. Callbacks (Agent Sessions)

### ✅ Tested (12/12 - 100%)

**Note:** This includes 10 Client API methods and 3 helper methods on the Callback type.

**Client API Methods:**

- **GetAllCallbacks()** - List all callbacks with filtering
  - File: `pkg/mythic/callbacks.go:106`
  - Tests: `tests/integration/callbacks_test.go:11`

- **GetCallbackByID()** - Get specific callback by display ID
  - File: `pkg/mythic/callbacks.go:181`
  - Tests: `tests/integration/callbacks_test.go:33`

- **GetAllActiveCallbacks()** - Filter only active callbacks
  - File: `pkg/mythic/callbacks.go:240`
  - Tests: `tests/integration/callbacks_test.go:51`

- **UpdateCallback()** - Update callback properties (description, ips, host, etc.)
  - File: `pkg/mythic/callbacks.go:293`
  - GraphQL: `updateCallback` mutation

- **CreateCallback()** - Manually register a new callback
  - File: `pkg/mythic/callbacks.go:396`
  - Tests: `tests/unit/callbacks_test.go:265`
  - GraphQL: `createCallback` mutation

- **DeleteCallback()** - Remove callback and associated tasks
  - File: `pkg/mythic/callbacks.go:445`
  - GraphQL: `deleteTasksAndCallbacks` mutation

- **AddCallbackGraphEdge()** - Add P2P connection between callbacks
  - File: `pkg/mythic/callbacks.go:501`
  - GraphQL: `callbackgraphedge_add` mutation

- **RemoveCallbackGraphEdge()** - Remove P2P connection
  - File: `pkg/mythic/callbacks.go:549`
  - GraphQL: `callbackgraphedge_remove` mutation

- **ExportCallbackConfig()** - Export callback configuration
  - File: `pkg/mythic/callbacks.go:588`
  - GraphQL: `exportCallbackConfig` query

- **ImportCallbackConfig()** - Import callback configuration
  - File: `pkg/mythic/callbacks.go:625`
  - GraphQL: `importCallbackConfig` mutation

**Helper Methods (on Callback type):**

- **Callback.IsActive()** - Helper to check if callback is active
  - File: `pkg/mythic/callbacks.go:379`

- **Callback.IsDead()** - Helper to check if callback is dead
  - File: `pkg/mythic/callbacks.go:384`

- **Callback.String()** - String representation
  - File: `pkg/mythic/callbacks.go:374`

---

## 3. Tasks (Commands)

### ✅ Tested (12/12 - 100%)

- **IssueTask()** - Issue task to callback(s)
  - File: `pkg/mythic/tasks.go:82`
  - Tests: `tests/integration/tasks_test.go`
  - GraphQL: `createTask` mutation

- **GetTask()** - Get task by ID
  - File: `pkg/mythic/tasks.go:171`
  - Tests: `tests/integration/tasks_test.go`

- **GetTasksForCallback()** - List all tasks for a callback
  - File: `pkg/mythic/tasks.go:239`

- **GetTaskOutput()** - Get task responses/output
  - File: `pkg/mythic/tasks.go:312`
  - Tests: `tests/integration/tasks_test.go`

- **UpdateTask()** - Add/update task comment
  - File: `pkg/mythic/tasks.go:389`

- **ReissueTask()** - Re-issue a task
  - File: `pkg/mythic/tasks.go:425`
  - GraphQL: `reissue_task` mutation

- **ReissueTaskWithHandler()** - Re-issue task with handler
  - File: `pkg/mythic/tasks.go:465`
  - GraphQL: `reissue_task_handler` mutation

- **RequestOpsecBypass()** - Request OPSEC bypass for blocked task
  - File: `pkg/mythic/tasks.go:505`
  - GraphQL: `requestOpsecBypass` mutation

- **AddMITREAttackToTask()** - Tag task with MITRE ATT&CK technique
  - File: `pkg/mythic/tasks.go:550`
  - GraphQL: `addAttackToTask` mutation

- **GetTasksByStatus()** - Filter tasks by status (preprocessing, submitted, etc.)
  - File: `pkg/mythic/tasks.go:597`
  - Database: `task` table with status filter

- **GetTaskArtifacts()** - Get artifacts created by task
  - File: `pkg/mythic/tasks.go:654`
  - Tests: `tests/unit/tasks_test.go:323`
  - Database: `taskartifact` table

- **WaitForTaskComplete()** - Wait for task to complete with timeout
  - File: `pkg/mythic/tasks.go` (helper method)
  - Tests: `tests/integration/tasks_test.go:TestTasks_WaitForTaskComplete_Timeout`
  - Polls task status until completion or timeout

### ⏳ Pending (0/12)

**Note:** Real-time task output subscriptions via WebSocket are a future enhancement and not part of the core 12 task operations.

---

## 4. Files

### ✅ Tested (8/8 - 100%)

- **GetFiles()** - List files with metadata
  - File: `pkg/mythic/files.go:102`
  - Tests: `tests/integration/files_test.go:11`

- **GetFileByID()** - Get specific file by agent_file_id
  - File: `pkg/mythic/files.go:184`
  - Tests: `tests/integration/files_test.go:150`

- **GetDownloadedFiles()** - Filter files downloaded from agents
  - File: `pkg/mythic/files.go:263`
  - Tests: `tests/integration/files_test.go:46`

- **UploadFile()** - Upload file to Mythic
  - File: `pkg/mythic/files.go:346`
  - Tests: `tests/integration/files_test.go:77`
  - REST: `POST /api/v1.4/task_upload_file_webhook`

- **DownloadFile()** - Download file content
  - File: `pkg/mythic/files.go:431`
  - Tests: `tests/integration/files_test.go:168`
  - REST: `GET /api/v1.4/files/download/{id}`

- **DeleteFile()** - Mark file as deleted
  - File: `pkg/mythic/files.go:496`
  - Tests: `tests/integration/files_test.go:237`
  - GraphQL: `update_filemeta` mutation

- **BulkDownloadFiles()** - Download multiple files as ZIP
  - File: `pkg/mythic/files.go:546`
  - Tests: `tests/integration/files_test.go:384`
  - GraphQL: `download_bulk` mutation

- **PreviewFile()** - Get file preview/metadata without downloading
  - File: `pkg/mythic/files.go:596`
  - Tests: `tests/integration/files_test.go:450`, `tests/unit/files_test.go:322`
  - GraphQL: `previewFile` mutation

---

## 5. Operations

### ✅ Tested (11/11 - 100%)

- **GetOperations()** - List all operations
  - File: `pkg/mythic/operations.go:12`
  - Tests: `tests/unit/operations_test.go`
  - Database: `operation` table

- **GetOperationByID()** - Get specific operation
  - File: `pkg/mythic/operations.go:77`
  - Tests: `tests/unit/operations_test.go`
  - Database: `operation` table

- **CreateOperation()** - Create new operation
  - File: `pkg/mythic/operations.go:145`
  - Tests: `tests/unit/operations_test.go`
  - GraphQL: `createOperation` mutation

- **UpdateOperation()** - Update operation details
  - File: `pkg/mythic/operations.go:186`
  - Tests: `tests/unit/operations_test.go`
  - GraphQL: `update_operation` mutation
  - Fields: name, channel, complete, webhook, admin_id, banner_text, banner_color

- **UpdateCurrentOperationForUser()** - Switch user's current operation
  - File: `pkg/mythic/operations.go:249`
  - GraphQL: `updateCurrentOperation` mutation

- **GetOperatorsByOperation()** - List operators in operation
  - File: `pkg/mythic/operations.go:285`
  - Tests: `tests/unit/operations_test.go`
  - Database: `operatoroperation` table

- **UpdateOperatorOperation()** - Add/remove operators from operation
  - File: `pkg/mythic/operations.go:325`
  - Tests: `tests/unit/operations_test.go`
  - GraphQL: `updateOperatorOperation` mutation

- **GetOperationEventLog()** - Get operation event logs
  - File: `pkg/mythic/operations.go:391`
  - Tests: `tests/unit/operations_test.go`
  - Database: `operationeventlog` table

- **CreateOperationEventLog()** - Create event log entry
  - File: `pkg/mythic/operations.go:441`
  - Tests: `tests/unit/operations_test.go`
  - GraphQL: `insert_operationeventlog` mutation

- **GetGlobalSettings()** - Get Mythic global settings
  - File: `pkg/mythic/operations.go:531`
  - GraphQL: `getGlobalSettings` query

- **UpdateGlobalSettings()** - Update global settings
  - File: `pkg/mythic/operations.go:547`
  - GraphQL: `updateGlobalSettings` mutation

---

## 6. Payloads

### ✅ Tested (12/12 - 100%)

**Note:** This includes 10 core Client API methods and 2 helper methods (WaitForPayloadComplete, DownloadPayload).

**Client API Methods:**

- **GetPayloads()** - List all payloads
  - File: `pkg/mythic/payloads.go:15`
  - Tests: `tests/integration/payloads_test.go:30`
  - Database: `payload` table

- **GetPayloadByUUID()** - Get specific payload
  - File: `pkg/mythic/payloads.go:82`
  - Tests: `tests/integration/payloads_test.go:240`
  - Database: `payload` table

- **CreatePayload()** - Build new payload
  - File: `pkg/mythic/payloads.go:156`
  - Tests: `tests/integration/payloads_test.go:78`
  - GraphQL: `createPayload` mutation
  - Input: JSON payload definition with type, commands, C2 profiles, build parameters

- **RebuildPayload()** - Rebuild existing payload
  - File: `pkg/mythic/payloads.go:282`
  - Tests: `tests/integration/payloads_test.go:330`
  - GraphQL: `rebuild_payload` mutation

- **UpdatePayload()** - Update payload settings
  - File: `pkg/mythic/payloads.go:220`
  - Tests: `tests/integration/payloads_test.go:175`
  - GraphQL: `updatePayload` mutation
  - Fields: callback_alert, description, deleted

- **DeletePayload()** - Delete payload
  - File: `pkg/mythic/payloads.go:276`
  - Tests: `tests/integration/payloads_test.go:388`
  - GraphQL: `deleteFile` mutation

- **ExportPayloadConfig()** - Export payload configuration
  - File: `pkg/mythic/payloads.go:317`
  - Tests: `tests/integration/payloads_test.go:531`
  - GraphQL: `exportPayloadConfig` query
  - Returns: JSON configuration string

- **GetPayloadTypes()** - List available payload types
  - File: `pkg/mythic/payloads.go:356`
  - Tests: `tests/integration/payloads_test.go:13`
  - Database: `payloadtype` table

- **GetPayloadCommands()** - Get commands for payload
  - File: `pkg/mythic/payloads.go:400`
  - Tests: `tests/integration/payloads_test.go:255`
  - Database: `payloadcommand` table
  - Input: payload ID (int)

- **GetPayloadOnHost()** - Track payloads deployed on hosts
  - File: `pkg/mythic/payloads.go:435`
  - Tests: `tests/integration/payloads_test.go:291`
  - Database: `payloadonhost` table
  - Input: operation ID

**Helper Methods:**

- **WaitForPayloadComplete()** - Wait for payload build to complete
  - File: `pkg/mythic/payloads.go:490`
  - Tests: `tests/integration/payloads_test.go:562`
  - Polls payload status until ready, failed, or timeout
  - Input: UUID, timeout in seconds

- **DownloadPayload()** - Download built payload file
  - File: `pkg/mythic/payloads.go:520`
  - Tests: `tests/integration/payloads_test.go:579`
  - REST: `GET /api/v1.4/files/download/{uuid}`
  - Returns: Binary payload data

**Helper Methods (on Payload type):**

- **Payload.IsReady()** - Check if payload build succeeded
  - File: `pkg/mythic/types/payload.go:149`
  - Tests: `tests/unit/payloads_test.go:82`

- **Payload.IsFailed()** - Check if payload build failed
  - File: `pkg/mythic/types/payload.go:154`
  - Tests: `tests/unit/payloads_test.go:94`

- **Payload.IsBuilding()** - Check if payload is still building
  - File: `pkg/mythic/types/payload.go:159`
  - Tests: `tests/unit/payloads_test.go:106`

- **Payload.String()** - String representation
  - File: `pkg/mythic/types/payload.go:132`
  - Tests: `tests/unit/payloads_test.go:22`

---

## 7. Credentials

### ✅ Tested (5/5 - 100%)

**Note:** This includes 3 core Client API methods plus 2 additional helper methods (GetCredentialsByOperation, DeleteCredential).

**Client API Methods:**

- **GetCredentials()** - List all credentials (non-deleted)
  - File: `pkg/mythic/credentials.go:11`
  - Tests: `tests/integration/credentials_test.go:12`
  - Database: `credential` table
  - Returns credentials sorted by timestamp (newest first)

- **GetCredentialsByOperation()** - List credentials for specific operation
  - File: `pkg/mythic/credentials.go:61`
  - Tests: `tests/integration/credentials_test.go:295`
  - Database: `credential` table with operation filter
  - Input: operation ID

- **CreateCredential()** - Add new credential
  - File: `pkg/mythic/credentials.go:118`
  - Tests: `tests/integration/credentials_test.go:47`
  - GraphQL: `createCredential` mutation
  - Fields: type, account, realm, credential, comment, task_id, metadata
  - Requires current operation to be set

- **UpdateCredential()** - Update credential
  - File: `pkg/mythic/credentials.go:213`
  - Tests: `tests/integration/credentials_test.go:47` (within create/update test)
  - GraphQL: `update_credential` mutation
  - Fields: type, account, realm, credential, comment, deleted, metadata
  - Supports partial updates (only specified fields)

- **DeleteCredential()** - Mark credential as deleted
  - File: `pkg/mythic/credentials.go:311`
  - Tests: `tests/integration/credentials_test.go:47` (cleanup)
  - Wrapper around UpdateCredential with deleted=true

**Helper Methods (on Credential type):**

- **Credential.String()** - String representation showing realm\account (type)
  - File: `pkg/mythic/types/credential.go:49`
  - Tests: `tests/unit/credentials_test.go:12`

- **Credential.IsDeleted()** - Check if credential is marked as deleted
  - File: `pkg/mythic/types/credential.go:61`
  - Tests: `tests/unit/credentials_test.go:66`

**Supported Credential Types:**
- `plaintext` - Plain text passwords
- `hash` - Password hashes (NTLM, etc.)
- `key` - SSH keys, API keys, etc.
- `ticket` - Kerberos tickets
- `cookie` - Session cookies
- `certificate` - SSL/TLS certificates

---

## 8. C2 Profiles

### ⏳ Pending (9/9)

- **GetC2Profiles()** - List all C2 profiles
  - Database: `c2profile` table

- **GetC2ProfileByID()** - Get specific C2 profile
  - Database: `c2profile` table

- **CreateC2Instance()** - Create C2 profile instance
  - GraphQL: `create_c2_instance` mutation

- **ImportC2Instance()** - Import C2 instance config
  - GraphQL: `import_c2_instance` mutation

- **StartStopProfile()** - Start/stop C2 profile
  - GraphQL: `startStopProfile` mutation

- **GetProfileOutput()** - Get C2 profile output/logs
  - GraphQL: `getProfileOutput` query

- **C2HostFile()** - Host file via C2 profile
  - GraphQL: `c2HostFile` mutation

- **C2SampleMessage()** - Generate sample C2 message
  - GraphQL: `c2SampleMessage` query

- **C2GetIOC()** - Get IOCs for C2 profile
  - GraphQL: `c2GetIOC` query

---

## 9. Artifacts (Indicators)

### ⏳ Pending (3/3)

- **GetArtifacts()** - List task artifacts/IOCs
  - Database: `artifact` table

- **CreateArtifact()** - Create artifact entry
  - GraphQL: `createArtifact` mutation

- **GetTaskArtifacts()** - Get artifacts for specific task
  - Database: `taskartifact` table

---

## 10. Tags

### ⏳ Pending (3/3)

- **GetTagTypes()** - List tag types
  - Database: `tagtype` table

- **CreateTag()** - Tag an object (task, callback, file, etc.)
  - GraphQL: `createTag` mutation

- **DeleteTagType()** - Delete tag type
  - GraphQL: `deleteTagtype` mutation

---

## 11. Tokens

### ⏳ Pending (4/4)

- **GetTokens()** - List tokens (process/user tokens)
  - Database: `token` table

- **GetCallbackTokens()** - Get tokens for callback
  - Database: `callbacktoken` table

- **GetAPITokens()** - List API tokens
  - Database: `apitokens` table

- **DeleteAPIToken()** - Delete API token
  - GraphQL: `deleteAPIToken` mutation

---

## 12. Processes

### ✅ Tested (6/6 - 100%)

**Note:** This includes 2 core Client API methods plus 4 additional helper methods for filtering and tree building.

**Client API Methods:**

- **GetProcesses()** - List all processes (non-deleted)
  - File: `pkg/mythic/processes.go:10`
  - Tests: `tests/integration/processes_test.go:13`
  - Database: `process` table
  - Returns processes sorted by timestamp (newest first)

- **GetProcessTree()** - Get process tree for callback
  - File: `pkg/mythic/processes.go:252`
  - Tests: `tests/integration/processes_test.go:160`
  - Database: `process` table with parent relationships
  - Returns hierarchical ProcessTree structure
  - Automatically builds parent-child relationships

**Helper Methods:**

- **GetProcessesByOperation()** - Filter processes by operation
  - File: `pkg/mythic/processes.go:76`
  - Tests: `tests/integration/processes_test.go:52`
  - Database: `process` table with operation filter

- **GetProcessesByCallback()** - Filter processes by callback
  - File: `pkg/mythic/processes.go:147`
  - Tests: `tests/integration/processes_test.go:95`
  - Database: `process` table with callback filter
  - Returns processes sorted by PID (ascending)

- **GetProcessesByHost()** - Filter processes by host
  - File: `pkg/mythic/processes.go:297`
  - Tests: `tests/integration/processes_test.go:234`
  - Database: `process` table with host filter
  - Returns processes sorted by PID (ascending)

- **buildProcessTree()** - Internal helper to build hierarchical tree
  - File: `pkg/mythic/processes.go:263`
  - Constructs parent-child relationships from flat process list

**Helper Methods (on Process type):**

- **Process.String()** - String representation showing name (PID)
  - File: `pkg/mythic/types/process.go:39`
  - Tests: `tests/unit/processes_test.go:12`

- **Process.IsDeleted()** - Check if process is marked as deleted
  - File: `pkg/mythic/types/process.go:51`
  - Tests: `tests/unit/processes_test.go:44`

- **Process.HasParent()** - Check if process has a parent process
  - File: `pkg/mythic/types/process.go:56`
  - Tests: `tests/unit/processes_test.go:60`

- **Process.GetIntegrityLevelString()** - Get human-readable integrity level
  - File: `pkg/mythic/types/process.go:61`
  - Tests: `tests/unit/processes_test.go:77`
  - Returns: Untrusted, Low, Medium, High, System, or Unknown

**Process Tree Structure:**
The ProcessTree type provides a hierarchical view of processes with automatic parent-child relationship building based on ProcessID and ParentProcessID fields.

---

## 13. Keylogs

### ✅ Tested (3/3 - 100%)

**Note:** This includes 2 core Client API methods plus 1 additional helper method (GetKeylogsByOperation).

**Client API Methods:**

- **GetKeylogs()** - List all keylog entries
  - File: `pkg/mythic/keylogs.go:10`
  - Tests: `tests/integration/keylogs_test.go:13`
  - Database: `keylog` table
  - Returns keylogs sorted by timestamp (newest first)

- **GetKeylogsByCallback()** - Filter keylogs by callback
  - File: `pkg/mythic/keylogs.go:103`
  - Tests: `tests/integration/keylogs_test.go:103`
  - Database: `keylog` table with callback filter
  - Returns keylogs sorted by timestamp (newest first)

**Helper Methods:**

- **GetKeylogsByOperation()** - Filter keylogs by operation
  - File: `pkg/mythic/keylogs.go:48`
  - Tests: `tests/integration/keylogs_test.go:54`
  - Database: `keylog` table with operation filter

**Helper Methods (on Keylog type):**

- **Keylog.String()** - String representation showing timestamp, window, and user
  - File: `pkg/mythic/types/keylog.go:18`
  - Tests: `tests/unit/keylogs_test.go:11`

- **Keylog.HasKeystrokes()** - Check if keylog has captured keystrokes
  - File: `pkg/mythic/types/keylog.go:30`
  - Tests: `tests/unit/keylogs_test.go:53`

---

## 14. Browser Scripts

### ⏳ Pending (3/3)

- **GetBrowserScripts()** - List browser scripts
  - Database: `browserscript` table

- **GetBrowserScriptsByOperation()** - Filter by operation
  - Database: `browserscriptoperation` table

- **CustomBrowserExport()** - Export browser data
  - GraphQL: `custombrowserExportFunction` mutation

---

## 15. MITRE ATT&CK

### ⏳ Pending (3/3)

- **GetAttackTechniques()** - List MITRE ATT&CK mappings
  - Database: `attack` table

- **GetAttackByTask()** - Get ATT&CK tags for task
  - Database: `attacktask` table

- **GetAttackByCommand()** - Get ATT&CK tags for command
  - Database: `attackcommand` table

---

## 16. Reporting

### ⏳ Pending (2/2)

- **GenerateReport()** - Generate operation report
  - GraphQL: `generateReport` mutation
  - Options: MITRE coverage, output format, filters

- **GetRedirectRules()** - Get C2 redirect rules
  - GraphQL: `redirect_rules` query

---

## 17. Eventing & Workflows

### ⏳ Pending (15/15)

- **EventingTriggerManual()** - Manually trigger event group
  - GraphQL: `eventingTriggerManual` mutation

- **EventingTriggerManualBulk()** - Bulk trigger on multiple objects
  - GraphQL: `eventingTriggerManualBulk` mutation

- **EventingTriggerKeyword()** - Trigger by keyword
  - GraphQL: `eventingTriggerKeyword` mutation

- **EventingTriggerCancel()** - Cancel running event
  - GraphQL: `eventingTriggerCancel` mutation

- **EventingTriggerRetry()** - Retry failed event
  - GraphQL: `eventingTriggerRetry` mutation

- **EventingTriggerRetryFromStep()** - Retry from specific step
  - GraphQL: `eventingTriggerRetryFromStep` mutation

- **EventingTriggerRunAgain()** - Re-run completed event
  - GraphQL: `eventingTriggerRunAgain` mutation

- **EventingTriggerUpdate()** - Update event group config
  - GraphQL: `eventingTriggerUpdate` mutation

- **EventingExportWorkflow()** - Export workflow definition
  - GraphQL: `eventingExportWorkflow` query

- **EventingImportContainerWorkflow()** - Import workflow
  - GraphQL: `eventingImportContainerWorkflow` mutation

- **EventingTestFile()** - Test workflow file
  - GraphQL: `eventingTestFile` query

- **UpdateEventGroupApproval()** - Approve/reject event execution
  - GraphQL: `updateEventGroupApproval` mutation

- **SendExternalWebhook()** - Send webhook notification
  - GraphQL: `sendExternalWebhook` mutation

- **ConsumingServicesTestWebhook()** - Test webhook service
  - GraphQL: `consumingServicesTestWebhook` mutation

- **ConsumingServicesTestLog()** - Test logging service
  - GraphQL: `consumingServicesTestLog` mutation

---

## 18. Operators (Users)

### ⏳ Pending (11/11)

- **GetOperators()** - List all operators
  - Database: `operator` table

- **GetOperatorByID()** - Get specific operator
  - Database: `operator` table

- **CreateOperator()** - Create new operator
  - GraphQL: `createOperator` mutation

- **UpdateOperatorStatus()** - Update operator status
  - GraphQL: `updateOperatorStatus` mutation
  - Fields: active, admin, deleted

- **UpdatePasswordAndEmail()** - Update credentials
  - GraphQL: `updatePasswordAndEmail` mutation

- **GetOperatorPreferences()** - Get UI preferences
  - GraphQL: `getOperatorPreferences` query

- **UpdateOperatorPreferences()** - Update preferences
  - GraphQL: `updateOperatorPreferences` mutation

- **GetOperatorSecrets()** - Get operator secrets
  - GraphQL: `getOperatorSecrets` query

- **UpdateOperatorSecrets()** - Update secrets
  - GraphQL: `updateOperatorSecrets` mutation

- **GetInviteLinks()** - List invite links
  - GraphQL: `getInviteLinks` query

- **CreateInviteLink()** - Create invite link for new operators
  - GraphQL: `createInviteLink` mutation

---

## 19. GraphQL Subscriptions

### ⏳ Pending (1/1)

- **Real-time subscriptions** - WebSocket-based real-time updates
  - Requires WebSocket transport for GraphQL client
  - Useful for: Task output streaming, callback status changes, new files

---

## 20. Advanced Features

### ⏳ Pending (20/20)

**Dynamic Queries:**
- **DynamicQueryFunction()** - Dynamic parameter queries
  - GraphQL: `dynamic_query_function` mutation

- **DynamicQueryBuildParameter()** - Build parameter queries
  - GraphQL: `dynamicQueryBuildParameterFunction` mutation

- **TypedarrayParseFunction()** - Parse typed arrays
  - GraphQL: `typedarray_parse_function` mutation

**Container Management:**
- **ContainerListFiles()** - List files in container
  - GraphQL: `containerListFiles` query

- **ContainerDownloadFile()** - Download from container
  - GraphQL: `containerDownloadFile` query

- **ContainerWriteFile()** - Write file to container
  - GraphQL: `containerWriteFile` mutation

- **ContainerRemoveFile()** - Remove container file
  - GraphQL: `containerRemoveFile` mutation

**Proxy Operations:**
- **ToggleProxy()** - Enable/disable SOCKS proxy
  - GraphQL: `toggleProxy` mutation

- **TestProxy()** - Test proxy connection
  - GraphQL: `testProxy` mutation

**File Browser:**
- **GetFileBrowserObjects()** - Browse file system
  - Database: `filebrowserobj` table

**Build Parameters:**
- **GetBuildParameters()** - List build parameters
  - Database: `buildparameter` table

- **GetBuildParameterInstances()** - Get parameter instances
  - Database: `buildparameterinstance` table

**Staging:**
- **GetStagingInfo()** - Get payload staging info
  - Database: `staginginfo` table

**Block Lists:**
- **DeleteBlockList()** - Delete block list
  - GraphQL: `deleteBlockList` mutation

- **DeleteBlockListEntry()** - Remove block list entries
  - GraphQL: `deleteBlockListEntry` mutation

**Commands:**
- **GetCommands()** - List available commands
  - Database: `command` table

- **GetCommandParameters()** - Get command parameters
  - Database: `commandparameters` table

- **GetLoadedCommands()** - Commands loaded in callback
  - Database: `loadedcommands` table

**Miscellaneous:**
- **CreateRandom()** - Generate random string with format
  - GraphQL: `createRandom` mutation

- **ConfigCheck()** - Check configuration
  - GraphQL: `config_check` query

---

## Implementation Priority Recommendations

### High Priority (Core Functionality)
1. ✅ **Operations Management** - Essential for multi-operation environments
2. ✅ **Payloads** - Critical for agent deployment
3. ✅ **Credentials** - Important for tracking compromised accounts
4. **C2 Profiles** - Needed for agent communication management
5. ✅ **Processes** - Important for situational awareness

### Medium Priority (Enhanced Features)
6. **Artifacts/IOCs** - Useful for tracking indicators
7. **Tags** - Organization and categorization
8. **Keylogs** - Credential harvesting operations
9. **MITRE ATT&CK** - Threat intelligence integration
10. **Reporting** - Operation documentation

### Low Priority (Advanced Features)
11. **Eventing/Workflows** - Automation for advanced users
12. **Browser Scripts** - Custom UI functionality
13. **Container Management** - Development/debugging
14. **Dynamic Queries** - Advanced parameter handling
15. **Proxy Operations** - Specialized networking

---

## Notes

- REST API endpoints are used for file upload/download operations
- GraphQL is used for all other operations
- Some features require WebSocket support for real-time subscriptions
- The Mythic API is under active development; new endpoints may be added

## Related Documentation

- **Mythic Documentation**: https://docs.mythic-c2.net/
- **GraphQL Schema**: Available via introspection at `{MYTHIC_URL}/graphql/`
- **REST API**: Documented in Mythic source at `/mythic-docker/src/webserver/controllers/`

---

*Last updated: 2026-01-07*
*SDK Version: In Development*
*Mythic API Version: v3.4.x*
