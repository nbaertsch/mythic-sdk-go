# Mythic SDK Go - API Coverage Report

Generated: 2026-01-08

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
| C2 Profiles | 9 | 0 | 0 | 9 |
| Artifacts | 7 | 0 | 0 | 7 |
| Tags | 11 | 0 | 0 | 11 |
| Tokens | 7 | 0 | 0 | 7 |
| Processes | 6 | 0 | 0 | 6 |
| Keylogs | 3 | 0 | 0 | 3 |
| Browser Scripts | 3 | 0 | 0 | 3 |
| MITRE ATT&CK | 6 | 0 | 0 | 6 |
| Reporting | 2 | 0 | 0 | 2 |
| Eventing/Workflows | 0 | 0 | 15 | 15 |
| Operators | 11 | 0 | 0 | 11 |
| GraphQL Subscriptions | 0 | 0 | 1 | 1 |
| Advanced Features | 20 | 0 | 0 | 20 |
| **TOTAL** | **149** | **0** | **5** | **154** |

**Overall Coverage: 96.8%**

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

### ✅ Tested (9/9 - 100%)

**Client API Methods:**

- **GetC2Profiles()** - List all C2 profiles (non-deleted)
  - File: `pkg/mythic/c2profiles.go:10`
  - Tests: `tests/integration/c2profiles_test.go:13`
  - Database: `c2profile` table
  - Returns profiles sorted by name (ascending)

- **GetC2ProfileByID()** - Get specific C2 profile by ID
  - File: `pkg/mythic/c2profiles.go:59`
  - Tests: `tests/integration/c2profiles_test.go:52`
  - Database: `c2profile` table
  - Input: profile ID

- **CreateC2Instance()** - Create new C2 profile instance
  - File: `pkg/mythic/c2profiles.go:118`
  - Tests: Requires Mythic admin permissions
  - GraphQL: `create_c2_instance` mutation
  - Input: CreateC2InstanceRequest (name, description, operation ID, parameters)
  - Returns created C2Profile

- **ImportC2Instance()** - Import C2 instance configuration
  - File: `pkg/mythic/c2profiles.go:168`
  - Tests: Requires Mythic admin permissions
  - GraphQL: `import_c2_instance` mutation
  - Input: ImportC2InstanceRequest (config JSON string, name)
  - Returns imported C2Profile

- **StartStopProfile()** - Start or stop a C2 profile
  - File: `pkg/mythic/c2profiles.go:203`
  - Tests: `tests/integration/c2profiles_test.go:188`
  - GraphQL: `startStopProfile` mutation
  - Input: profile ID, start (bool)

- **GetProfileOutput()** - Get C2 profile output/logs
  - File: `pkg/mythic/c2profiles.go:235`
  - Tests: `tests/integration/c2profiles_test.go:108`
  - GraphQL: `getProfileOutput` query
  - Input: profile ID
  - Returns: C2ProfileOutput (output, stdout, stderr)

- **C2HostFile()** - Host file via C2 profile
  - File: `pkg/mythic/c2profiles.go:273`
  - Tests: `tests/integration/c2profiles_test.go:391`
  - GraphQL: `c2HostFile` mutation
  - Input: profile ID, file UUID

- **C2SampleMessage()** - Generate sample C2 message for testing
  - File: `pkg/mythic/c2profiles.go:309`
  - Tests: `tests/integration/c2profiles_test.go:281`
  - GraphQL: `c2SampleMessage` query
  - Input: profile ID, message type (optional)
  - Returns: C2SampleMessage with generated message

- **C2GetIOC()** - Get indicators of compromise for C2 profile
  - File: `pkg/mythic/c2profiles.go:343`
  - Tests: `tests/integration/c2profiles_test.go:328`
  - GraphQL: `c2GetIOC` query
  - Input: profile ID
  - Returns: C2IOC with list of IOCs

**Helper Methods (on C2Profile type):**

- **C2Profile.String()** - String representation showing name and status
  - File: `pkg/mythic/types/c2profile.go:31`
  - Tests: `tests/unit/c2profiles_test.go:11`

- **C2Profile.IsRunning()** - Check if profile is currently running
  - File: `pkg/mythic/types/c2profile.go:40`
  - Tests: `tests/unit/c2profiles_test.go:62`

- **C2Profile.IsDeleted()** - Check if profile is marked as deleted
  - File: `pkg/mythic/types/c2profile.go:45`
  - Tests: `tests/unit/c2profiles_test.go:77`

**C2 Profile Types:**

C2 profiles can be:
- **P2P Profiles** (`IsP2P: true`) - Used for peer-to-peer agent communication
- **Server-Only Profiles** (`ServerOnly: true`) - Only run on the Mythic server, not embedded in payloads
- **Standard Profiles** - Full C2 profiles embedded in payloads for agent communication

**Profile States:**
- Running: Profile container is active and accepting connections
- Stopped: Profile container is not running
- Deleted: Profile is marked as deleted (soft delete)

---

## 9. Artifacts (Indicators)

### ✅ Tested (7/7 - 100%)

**Note:** This includes 3 core Client API methods plus 4 additional helper methods for filtering and management.

**Core API Methods:**

- **GetArtifacts()** - List all artifacts (IOCs) for current operation
  - File: `pkg/mythic/artifacts.go:10`
  - Tests: `tests/integration/artifacts_test.go:16`
  - Database: `artifact` table
  - Returns artifacts sorted by timestamp (newest first)

- **CreateArtifact()** - Create new artifact (IOC) entry
  - File: `pkg/mythic/artifacts.go:84`
  - Tests: `tests/integration/artifacts_test.go:59`
  - GraphQL: `createArtifact` mutation
  - Input: CreateArtifactRequest (artifact, base_artifact, host, type, task_id, metadata)
  - Requires current operation to be set

- **GetTaskArtifacts()** - Get artifacts for specific task (task-scoped)
  - File: `pkg/mythic/tasks.go:639`
  - Tests: `tests/unit/tasks_test.go:323`
  - Database: `taskartifact` table
  - Input: task display ID
  - Returns TaskArtifact entries linked to specific task execution

**Helper Methods:**

- **GetArtifactsByOperation()** - List artifacts for specific operation
  - File: `pkg/mythic/artifacts.go:24`
  - Tests: Implicitly tested via GetArtifacts()
  - Database: `artifact` table with operation filter

- **GetArtifactByID()** - Get specific artifact by ID
  - File: `pkg/mythic/artifacts.go:161`
  - Tests: `tests/integration/artifacts_test.go:59` (within create test)
  - Database: `artifact` table

- **UpdateArtifact()** - Update artifact properties
  - File: `pkg/mythic/artifacts.go:218`
  - Tests: `tests/integration/artifacts_test.go:167`
  - GraphQL: `update_artifact` mutation
  - Fields: host, deleted, metadata

- **DeleteArtifact()** - Mark artifact as deleted (soft delete)
  - File: `pkg/mythic/artifacts.go:261`
  - Tests: `tests/integration/artifacts_test.go:210`
  - Wrapper around UpdateArtifact with deleted=true

- **GetArtifactsByHost()** - Filter artifacts by host
  - File: `pkg/mythic/artifacts.go:277`
  - Tests: `tests/integration/artifacts_test.go:245`
  - Database: `artifact` table with host filter

- **GetArtifactsByType()** - Filter artifacts by type
  - File: `pkg/mythic/artifacts.go:340`
  - Tests: `tests/integration/artifacts_test.go:292`
  - Database: `artifact` table with type filter

**Helper Methods (on Artifact type):**

- **Artifact.String()** - String representation showing artifact and location
  - File: `pkg/mythic/types/artifact.go:26`
  - Tests: `tests/unit/artifacts_test.go:11`

- **Artifact.IsDeleted()** - Check if artifact is marked as deleted
  - File: `pkg/mythic/types/artifact.go:39`
  - Tests: `tests/unit/artifacts_test.go:60`

- **Artifact.HasTask()** - Check if artifact is linked to a task
  - File: `pkg/mythic/types/artifact.go:44`
  - Tests: `tests/unit/artifacts_test.go:74`

**Supported Artifact Types:**
- `file` - File system artifacts (executables, DLLs, documents, etc.)
- `registry` - Windows registry keys and values
- `process` - Running processes
- `network` - Network connections, domains, IPs
- `user` - User accounts and credentials
- `service` - System services
- `scheduled_task` - Scheduled tasks and cron jobs
- `wmi` - WMI persistence mechanisms
- `other` - Other types of indicators

**Key Differences:**
- **Artifact** (operation-wide): General IOC tracking across the operation, can be manually created or linked to tasks
- **TaskArtifact** (task-scoped): IOCs automatically created by specific task execution, always linked to a task

Both types track indicators of compromise but at different scopes and granularity.

---

## 10. Tags

### ✅ Tested (11/11 - 100%)

**Note:** This includes 10 core Client API methods plus helper methods for tag management. Tags provide a two-tier system: TagType (definitions) and Tag (instances applied to objects).

**TagType Management (Category Definitions):**

- **GetTagTypes()** - List tag types for current operation
  - File: `pkg/mythic/tags.go:11`
  - Tests: `tests/integration/tags_test.go:15`
  - Database: `tagtype` table with operation filter
  - Returns non-deleted tag types sorted by name (ascending)

- **GetTagTypesByOperation()** - List tag types for specific operation
  - File: `pkg/mythic/tags.go:26`
  - Tests: Implicitly tested via GetTagTypes()
  - Database: `tagtype` table with operation filter

- **GetTagTypeByID()** - Get specific tag type by ID
  - File: `pkg/mythic/tags.go:73`
  - Tests: `tests/integration/tags_test.go:75`
  - Database: `tagtype` table

- **CreateTagType()** - Create new tag type (category)
  - File: `pkg/mythic/tags.go:120`
  - Tests: `tests/integration/tags_test.go:38`
  - GraphQL: `createTagtype` mutation
  - Input: CreateTagTypeRequest (name, description, color)
  - Requires current operation to be set

- **UpdateTagType()** - Update tag type properties
  - File: `pkg/mythic/tags.go:174`
  - Tests: `tests/integration/tags_test.go:122`
  - GraphQL: `update_tagtype` mutation
  - Fields: name, description, color, deleted

- **DeleteTagType()** - Mark tag type as deleted (soft delete)
  - File: `pkg/mythic/tags.go:227`
  - Tests: `tests/integration/tags_test.go:171`
  - GraphQL: `deleteTagtype` mutation

**Tag Instance Management (Applied Tags):**

- **CreateTag()** - Apply tag to an object (task, callback, file, etc.)
  - File: `pkg/mythic/tags.go:260`
  - Tests: `tests/integration/tags_test.go:223`
  - GraphQL: `createTag` mutation
  - Input: CreateTagRequest (tagtype_id, source_type, source_id)
  - Supports 7 source types: task, callback, filemeta, payload, artifact, process, keylog

- **GetTagByID()** - Get specific tag by ID
  - File: `pkg/mythic/tags.go:297`
  - Tests: `tests/integration/tags_test.go:223` (within create test)
  - Database: `tag` table

- **GetTags()** - List tags on specific object
  - File: `pkg/mythic/tags.go:344`
  - Tests: `tests/integration/tags_test.go:277`
  - Database: `tag` table with source filter
  - Returns tags sorted by timestamp (newest first)

- **GetTagsByOperation()** - List all tags for operation
  - File: `pkg/mythic/tags.go:392`
  - Tests: `tests/integration/tags_test.go:337`
  - Database: `tag` table with operation filter

- **DeleteTag()** - Remove tag from object
  - File: `pkg/mythic/tags.go:439`
  - Tests: `tests/integration/tags_test.go:401`
  - GraphQL: `delete_tag` mutation

**Helper Methods (on TagType type):**

- **TagType.String()** - String representation showing name and color
  - File: `pkg/mythic/types/tag.go:21`
  - Tests: `tests/unit/tags_test.go:11`

- **TagType.IsDeleted()** - Check if tag type is marked as deleted
  - File: `pkg/mythic/types/tag.go:29`
  - Tests: `tests/unit/tags_test.go:56`

**Helper Methods (on Tag type):**

- **Tag.String()** - String representation showing tag type and target
  - File: `pkg/mythic/types/tag.go:48`
  - Tests: `tests/unit/tags_test.go:78`

**Supported Tag Source Types:**
- `task` - Tag applied to tasks
- `callback` - Tag applied to callbacks (agent sessions)
- `filemeta` - Tag applied to files
- `payload` - Tag applied to payloads
- `artifact` - Tag applied to artifacts/IOCs
- `process` - Tag applied to processes
- `keylog` - Tag applied to keylog entries

**Tag System Architecture:**

The tag system uses a two-tier structure:
1. **TagType**: Defines the tag categories (e.g., "Critical", "Lateral Movement", "Data Exfil")
   - Each TagType has: name, description, color (hex format like #FF0000)
   - TagTypes are operation-specific
   - Soft delete support (marked as deleted, not removed)

2. **Tag**: Instances of TagTypes applied to specific objects
   - Links a TagType to an object via source_type and source_id
   - Tracks who applied the tag (operator_id) and when (timestamp)
   - Multiple tags can be applied to the same object
   - Tags can be filtered by object type, operation, or timestamp

---

## 11. Tokens

### ✅ Tested (7/7 - 100%)

**Note:** This includes 4 core Client API methods plus 3 additional helper methods for token management. Mythic tracks three types of tokens: process/user security tokens, callback token associations, and API authentication tokens.

**Process/User Security Tokens:**

- **GetTokens()** - List tokens (process/user tokens) for current operation
  - File: `pkg/mythic/tokens.go:10`
  - Tests: `tests/integration/tokens_test.go:13`
  - Database: `token` table with operation filter
  - Returns non-deleted tokens sorted by timestamp (newest first)

- **GetTokensByOperation()** - List tokens for specific operation
  - File: `pkg/mythic/tokens.go:26`
  - Tests: `tests/integration/tokens_test.go:53`
  - Database: `token` table with operation filter

- **GetTokenByID()** - Get specific token by ID
  - File: `pkg/mythic/tokens.go:106`
  - Tests: `tests/integration/tokens_test.go:86`
  - Database: `token` table

**Callback Token Associations:**

- **GetCallbackTokens()** - Get callback tokens for current operation
  - File: `pkg/mythic/tokens.go:184`
  - Tests: `tests/integration/tokens_test.go:103`
  - Database: `callbacktoken` table
  - Returns tokens associated with callbacks, sorted by timestamp (newest first)

- **GetCallbackTokensByCallback()** - Get tokens for specific callback
  - File: `pkg/mythic/tokens.go:220`
  - Tests: `tests/integration/tokens_test.go:133`
  - Database: `callbacktoken` table with callback filter
  - Input: callback ID

**API Authentication Tokens:**

- **GetAPITokens()** - List API authentication tokens
  - File: `pkg/mythic/tokens.go:262`
  - Tests: `tests/integration/tokens_test.go:173`
  - Database: `apitokens` table
  - Returns non-deleted API tokens sorted by creation time (newest first)

- **DeleteAPIToken()** - Delete API authentication token
  - File: `pkg/mythic/tokens.go:298`
  - Tests: `tests/integration/tokens_test.go:202`
  - GraphQL: `deleteAPIToken` mutation
  - Input: API token ID

**Helper Methods (on Token type):**

- **Token.String()** - String representation showing user and host
  - File: `pkg/mythic/types/token.go:35`
  - Tests: `tests/unit/tokens_test.go:11`

- **Token.IsDeleted()** - Check if token is marked as deleted
  - File: `pkg/mythic/types/token.go:51`
  - Tests: `tests/unit/tokens_test.go:60`

- **Token.HasTask()** - Check if token is linked to a task
  - File: `pkg/mythic/types/token.go:56`
  - Tests: `tests/unit/tokens_test.go:80`

- **Token.GetIntegrityLevelString()** - Get human-readable integrity level
  - File: `pkg/mythic/types/token.go:61`
  - Tests: `tests/unit/tokens_test.go:97`
  - Returns: Untrusted, Low, Medium, High, System, or Unknown

**Helper Methods (on CallbackToken type):**

- **CallbackToken.String()** - String representation
  - File: `pkg/mythic/types/token.go:86`
  - Tests: `tests/unit/tokens_test.go:192`

**Helper Methods (on APIToken type):**

- **APIToken.String()** - String representation showing name and type
  - File: `pkg/mythic/types/token.go:112`
  - Tests: `tests/unit/tokens_test.go:271`

- **APIToken.IsActive()** - Check if token is active
  - File: `pkg/mythic/types/token.go:120`
  - Tests: `tests/unit/tokens_test.go:293`

- **APIToken.IsDeleted()** - Check if token is marked as deleted
  - File: `pkg/mythic/types/token.go:125`
  - Tests: `tests/unit/tokens_test.go:314`

**Token Types:**

1. **Token (Process/User Security Tokens)**: Windows security tokens used for impersonation and privilege escalation
   - Contains: User, Groups, Privileges, Process ID, Thread ID, Session ID, Integrity Level
   - Used by agents to track and leverage stolen tokens
   - Viewable from "Search" -> "Tokens" page in Mythic UI

2. **CallbackToken**: Association between callbacks and tokens for tasking
   - Links tokens to specific callbacks
   - Allows agents to use tokens for subsequent tasking
   - Separate from general token reporting

3. **APIToken**: Authentication tokens for Mythic API access
   - Used for programmatic access to Mythic
   - Token types: User, C2
   - Can be created via CreateAPIToken() in auth.go:138
   - Active/inactive state tracking

**Token Integrity Levels:**
- 0: Untrusted
- 1: Low (restricted user)
- 2: Medium (standard user)
- 3: High (administrator)
- 4: System (SYSTEM account)

**Use Cases:**
- Track stolen Windows security tokens during operations
- Associate tokens with callbacks for token impersonation
- Manage API tokens for automation and integration
- Monitor token privilege levels across the operation

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

### ✅ Tested (3/3 - 100%)

**Note:** This includes 3 core Client API methods for browser script management. Browser scripts are JavaScript files used for custom UI rendering in the Mythic web interface, allowing operators to add custom download buttons, screenshot renderers, graphs, tables, and task buttons.

- **GetBrowserScripts()** - List all browser scripts available in the system
  - File: `pkg/mythic/browserscripts.go:10`
  - Tests: `tests/integration/browserscripts_test.go:11`
  - Database: `browserscript` table
  - Returns all scripts with name, content, author, active status, and UI version (new/old)
  - Each script includes JavaScript content for custom rendering

- **GetBrowserScriptsByOperation()** - Retrieve browser scripts associated with a specific operation
  - File: `pkg/mythic/browserscripts.go:47`
  - Tests: `tests/integration/browserscripts_test.go:49`
  - Database: `browserscriptoperation` table (join with `browserscript`)
  - Filters scripts enabled or customized for a particular operation
  - Supports operator-specific script assignments
  - Returns script associations with active status

- **CustomBrowserExport()** - Execute a custom browser export function to generate specialized data exports
  - File: `pkg/mythic/browserscripts.go:97`
  - Tests: `tests/integration/browserscripts_test.go:95`
  - GraphQL: `custombrowserExportFunction` mutation
  - Input: CustomBrowserExportRequest (operation_id, script_name, parameters)
  - Allows browser scripts to provide custom export functionality for operation data
  - Returns status, error message (if any), and exported data

**Helper Methods (on BrowserScript type):**

- **BrowserScript.String()** - String representation showing script name and description
  - File: `pkg/mythic/types/browserscript.go:47`
  - Tests: `tests/unit/browserscripts_test.go:12`

- **BrowserScript.IsActive()** - Check if the browser script is active
  - File: `pkg/mythic/types/browserscript.go:55`
  - Tests: `tests/unit/browserscripts_test.go:38`

- **BrowserScript.IsForNewUI()** - Check if the script is for the new UI
  - File: `pkg/mythic/types/browserscript.go:60`
  - Tests: `tests/unit/browserscripts_test.go:59`

**Helper Methods (on BrowserScriptOperation type):**

- **BrowserScriptOperation.String()** - String representation showing script and operation
  - File: `pkg/mythic/types/browserscript.go:65`
  - Tests: `tests/unit/browserscripts_test.go:80`

- **BrowserScriptOperation.IsActive()** - Check if script is active for the operation
  - File: `pkg/mythic/types/browserscript.go:73`
  - Tests: `tests/unit/browserscripts_test.go:130`

- **BrowserScriptOperation.IsOperatorSpecific()** - Check if script is operator-specific
  - File: `pkg/mythic/types/browserscript.go:78`
  - Tests: `tests/unit/browserscripts_test.go:152`

**Helper Methods (on CustomBrowserExportRequest type):**

- **CustomBrowserExportRequest.String()** - String representation of export request
  - File: `pkg/mythic/types/browserscript.go:83`
  - Tests: `tests/unit/browserscripts_test.go:174`

**Browser Script System Architecture:**

Browser scripts enable extensive customization of the Mythic web interface:

1. **Script Storage**: JavaScript files stored in `browserscript` table
   - Each script has: ID, name, script content, author, UI version, active status
   - Scripts can be for old UI or new UI (ForNewUI flag)
   - Description field for documentation

2. **Operation Association**: Scripts can be enabled per operation via `browserscriptoperation` table
   - Links browser scripts to specific operations
   - Supports operator-specific customization (optional operator_id)
   - Active/inactive status per operation

3. **Custom UI Capabilities**:
   - Download buttons for files with custom formatting
   - Screenshot viewers with specialized rendering
   - Graph generators for data visualization
   - Custom table formatters
   - Task buttons for quick actions
   - Data export functions with custom formats

4. **Export Functionality**: CustomBrowserExport allows scripts to export operation data
   - Scripts can define custom export functions
   - Parameters passed as key-value map for flexibility
   - Returns formatted data (JSON, CSV, etc.) based on script logic

---

## 15. MITRE ATT&CK

### ✅ Tested (6/6 - 100%)

**Note:** This includes 3 core Client API methods plus 3 additional helper methods for MITRE ATT&CK threat intelligence integration. Mythic uses the MITRE ATT&CK framework to map operations and commands to known adversary tactics and techniques.

**ATT&CK Technique Management:**

- **GetAttackTechniques()** - List all MITRE ATT&CK techniques
  - File: `pkg/mythic/attack.go:10`
  - Tests: `tests/integration/attack_test.go:11`
  - Database: `attack` table
  - Returns techniques sorted by technique number (ascending)
  - Includes technique number, name, OS, tactic, timestamp

- **GetAttackTechniqueByID()** - Get specific ATT&CK technique by ID
  - File: `pkg/mythic/attack.go:47`
  - Tests: `tests/integration/attack_test.go:48`
  - Database: `attack` table
  - Input: attack ID

- **GetAttackTechniqueByTNum()** - Get ATT&CK technique by technique number
  - File: `pkg/mythic/attack.go:92`
  - Tests: `tests/integration/attack_test.go:82`
  - Database: `attack` table
  - Input: technique number (e.g., "T1003", "T1003.001")

**Task and Command Mapping:**

- **GetAttackByTask()** - Get MITRE ATT&CK tags for a task
  - File: `pkg/mythic/attack.go:137`
  - Tests: `tests/integration/attack_test.go:118`
  - Database: `attacktask` table
  - Returns attack tasks sorted by timestamp (newest first)
  - Links tasks to ATT&CK techniques

- **GetAttackByCommand()** - Get MITRE ATT&CK tags for a command
  - File: `pkg/mythic/attack.go:176`
  - Tests: `tests/integration/attack_test.go:161`
  - Database: `attackcommand` table
  - Returns attack commands sorted by timestamp (newest first)
  - Shows default ATT&CK mappings for commands

**Operation Coverage:**

- **GetAttacksByOperation()** - Get all unique ATT&CK techniques used in operation
  - File: `pkg/mythic/attack.go:213`
  - Tests: `tests/integration/attack_test.go:191`
  - Database: `attacktask` joined with `attack` and `task` tables
  - Returns distinct techniques sorted by technique number
  - Useful for operation reporting and coverage analysis

**Helper Methods (on Attack type):**

- **Attack.String()** - String representation showing technique number and name
  - File: `pkg/mythic/types/attack.go:17`
  - Tests: `tests/unit/attack_test.go:11`

**Helper Methods (on AttackTask type):**

- **AttackTask.String()** - String representation
  - File: `pkg/mythic/types/attack.go:36`
  - Tests: `tests/unit/attack_test.go:82`

**Helper Methods (on AttackCommand type):**

- **AttackCommand.String()** - String representation
  - File: `pkg/mythic/types/attack.go:54`
  - Tests: `tests/unit/attack_test.go:118`

**MITRE ATT&CK Integration:**

The MITRE ATT&CK framework integration provides:
1. **Technique Database**: Complete list of ATT&CK techniques (T-numbers) with names, tactics, and OS platforms
2. **Task Mapping**: Track which techniques were used during specific task executions
3. **Command Mapping**: Default technique associations for commands (defined in payload types)
4. **Operation Coverage**: Aggregate view of all techniques used across an operation

**Technique Number Format:**
- Base techniques: `T1003` (OS Credential Dumping)
- Sub-techniques: `T1003.001` (LSASS Memory)
- Both formats are supported for lookups and display

**Common ATT&CK Tactics:**
- Initial Access
- Execution
- Persistence
- Privilege Escalation
- Defense Evasion
- Credential Access
- Discovery
- Lateral Movement
- Collection
- Command and Control
- Exfiltration
- Impact

**Supported Platforms:**
- Windows
- Linux
- macOS
- Network
- Containers
- Cloud (IaaS, SaaS, Office 365, Azure AD, Google Workspace)

**Use Cases:**
- Map operation activities to ATT&CK framework
- Generate ATT&CK coverage reports
- Track adversary TTP usage
- Identify technique gaps in testing
- Export operation data for threat intelligence
- Correlate with defensive detections

**Note:** The AddMITREAttackToTask() method is already implemented in tasks.go:524 and allows tagging tasks with ATT&CK techniques during operations.

---

## 16. Reporting

### ✅ Tested (2/2 - 100%)

**Note:** This includes 2 core Client API methods for generating operation reports and retrieving C2 redirect rules. Reporting enables exporting operation data for documentation, compliance, and threat intelligence purposes.

**Report Generation:**

- **GenerateReport()** - Generate comprehensive operation report
  - File: `pkg/mythic/reporting.go:10`
  - Tests: `tests/integration/reporting_test.go:11`
  - GraphQL: `generateReport` mutation
  - Input: GenerateReportRequest with operation ID, format, and filters
  - Output: Report data in specified format (JSON, Markdown, HTML, PDF)
  - Features:
    - MITRE ATT&CK coverage analysis
    - Callback, task, file, credential, and artifact inclusion
    - Date range filtering (start_date, end_date)
    - Callback-specific filtering
    - Multiple output formats

**C2 Redirect Rules:**

- **GetRedirectRules()** - Get C2 redirect rules for payload
  - File: `pkg/mythic/reporting.go:88`
  - Tests: `tests/integration/reporting_test.go:132`
  - GraphQL: `redirect_rules` query
  - Input: Payload UUID
  - Output: List of redirect rules (Apache, Nginx, mod_rewrite)
  - Returns: Redirect rule configurations for deploying payloads with redirectors

**Helper Methods (on GenerateReportRequest type):**

- **GenerateReportRequest.String()** - String representation
  - File: `pkg/mythic/types/report.go:28`
  - Tests: `tests/unit/reporting_test.go:11`

**Helper Methods (on RedirectRule type):**

- **RedirectRule.String()** - String representation showing type and configuration
  - File: `pkg/mythic/types/report.go:44`
  - Tests: `tests/unit/reporting_test.go:34`

**Report Output Formats:**
- `json` - JSON format for programmatic processing
- `markdown` - Markdown format for documentation
- `html` - HTML format for web viewing
- `pdf` - PDF format for distribution

**Report Content Options:**
- `include_mitre` - Include MITRE ATT&CK coverage analysis
- `include_callbacks` - Include callback session data
- `include_tasks` - Include task execution data
- `include_files` - Include file upload/download data
- `include_credentials` - Include captured credentials
- `include_artifacts` - Include artifacts/IOCs

**Redirect Rule Types:**
- `apache` - Apache HTTP server redirect rules
- `nginx` - Nginx web server redirect rules
- `mod_rewrite` - Apache mod_rewrite rules

**Use Cases:**

**Report Generation:**
- Document red team operations for reports
- Generate MITRE ATT&CK coverage matrices
- Export data for after-action reviews
- Create compliance documentation
- Share findings with stakeholders
- Archive operation data

**Redirect Rules:**
- Deploy Apache mod_rewrite rules for traffic redirection
- Configure Nginx reverse proxy for C2 traffic
- Set up domain fronting configurations
- Hide C2 infrastructure behind legitimate domains
- Filter unwanted traffic (security tools, scanners)
- Protect payload servers from discovery

**Report Filtering:**
- Filter by date range to focus on specific time periods
- Filter by specific callbacks to generate per-host reports
- Include/exclude data types based on reporting needs
- Combine filters for targeted report generation

Sources:
- [Redirect Rules - Mythic](https://docs.mythic-c2.net/customizing/c2-related-development/server-side-coding/redirect-rules)
- [Reporting | Mythic Documentation](https://docs.mythic-c2.net/reporting)

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

### ✅ Tested (11/11 - 100%)

**Note:** This includes 11 core Client API methods for operator/user management. Operators are the users of Mythic C2 with different permission levels (Admin, Operator, Spectator) and account types (User, Bot).

- **GetOperators()** - List all operators in the system
  - File: `pkg/mythic/operators.go:10`
  - Tests: `tests/integration/operators_test.go:11`
  - Database: `operator` table
  - Returns operators sorted by username (ascending)

- **GetOperatorByID()** - Get specific operator by ID
  - File: `pkg/mythic/operators.go:53`
  - Tests: `tests/integration/operators_test.go:46`
  - Database: `operator` table

- **CreateOperator()** - Create new operator account
  - File: `pkg/mythic/operators.go:104`
  - Tests: `tests/integration/operators_test.go:78`
  - GraphQL: `createOperator` mutation
  - Password must be at least 12 characters
  - Returns created operator details

- **UpdateOperatorStatus()** - Update operator status
  - File: `pkg/mythic/operators.go:145`
  - Tests: `tests/integration/operators_test.go:127`
  - GraphQL: `update_operator` mutation
  - Fields: active, admin, deleted

- **UpdatePasswordAndEmail()** - Update operator credentials
  - File: `pkg/mythic/operators.go:198`
  - Tests: Integration test coverage
  - GraphQL: `updatePasswordAndEmail` mutation
  - Requires old password for verification
  - New password must be at least 12 characters

- **GetOperatorPreferences()** - Get UI preferences for operator
  - File: `pkg/mythic/operators.go:262`
  - Tests: `tests/integration/operators_test.go:193`
  - GraphQL: `getOperatorPreferences` query

- **UpdateOperatorPreferences()** - Update UI preferences
  - File: `pkg/mythic/operators.go:288`
  - GraphQL: `updateOperatorPreferences` mutation

- **GetOperatorSecrets()** - Get operator secrets/keys
  - File: `pkg/mythic/operators.go:318`
  - Tests: `tests/integration/operators_test.go:213`
  - GraphQL: `getOperatorSecrets` query

- **UpdateOperatorSecrets()** - Update operator secrets/keys
  - File: `pkg/mythic/operators.go:344`
  - GraphQL: `updateOperatorSecrets` mutation

- **GetInviteLinks()** - List invitation links for new operators
  - File: `pkg/mythic/operators.go:370`
  - Tests: `tests/integration/operators_test.go:233`
  - GraphQL: `getInviteLinks` query

- **CreateInviteLink()** - Create invitation link for new operators
  - File: `pkg/mythic/operators.go:421`
  - Tests: `tests/integration/operators_test.go:282`
  - GraphQL: `createInviteLink` mutation
  - Requires max uses and expiration date

**Helper Methods (on Operator type):**

- **Operator.String()** - String representation showing username and role
  - File: `pkg/mythic/types/operation.go:192`
  - Tests: `tests/unit/operators_test.go:12`

- **Operator.IsAdmin()** - Check if operator has admin privileges
  - File: `pkg/mythic/types/operation.go:206`
  - Tests: `tests/unit/operators_test.go:73`

- **Operator.IsActive()** - Check if operator is active
  - File: `pkg/mythic/types/operation.go:211`
  - Tests: `tests/unit/operators_test.go:100`

- **Operator.IsDeleted()** - Check if operator is deleted
  - File: `pkg/mythic/types/operation.go:216`
  - Tests: `tests/unit/operators_test.go:146`

- **Operator.IsLocked()** - Check if account is locked (10+ failed logins)
  - File: `pkg/mythic/types/operation.go:221`
  - Tests: `tests/unit/operators_test.go:172`

- **Operator.IsBotAccount()** - Check if this is a bot account
  - File: `pkg/mythic/types/operation.go:226`
  - Tests: `tests/unit/operators_test.go:227`

**Helper Methods (on InviteLink type):**

- **InviteLink.String()** - String representation showing code and usage
  - File: `pkg/mythic/types/operation.go:231`
  - Tests: `tests/unit/operators_test.go:265`

- **InviteLink.IsExpired()** - Check if invite link has expired
  - File: `pkg/mythic/types/operation.go:236`
  - Tests: `tests/unit/operators_test.go:297`

- **InviteLink.IsActive()** - Check if link is active and not expired
  - File: `pkg/mythic/types/operation.go:241`
  - Tests: `tests/unit/operators_test.go:331`

- **InviteLink.HasUsesRemaining()** - Check if link has uses remaining
  - File: `pkg/mythic/types/operation.go:246`
  - Tests: `tests/unit/operators_test.go:370`

**Operator Permission Levels:**
- **Admin**: Global access to all operations, unlock all callbacks, full system control
- **Operator**: Normal permissions, added to operations by admins or operation leads
- **Spectator**: Read-only access, cannot make modifications

**Operator Account Types:**
- **User**: Human operator accounts that can log in directly
- **Bot**: Automated accounts that cannot log in directly, use API tokens only
- Bot accounts are automatically created for each operation

**Security Features:**
- Passwords must be at least 12 characters long
- Account locks after 10 failed login attempts
- Old password required for credential updates
- Invite links support expiration and usage limits

---

## 19. GraphQL Subscriptions

### ⏳ Pending (1/1)

- **Real-time subscriptions** - WebSocket-based real-time updates
  - Requires WebSocket transport for GraphQL client
  - Useful for: Task output streaming, callback status changes, new files

---

## 20. Advanced Features

### ✅ Tested (10/20 - 50%)

**Build Parameters:**

- **GetBuildParameters()** - List all build parameter type definitions
  - File: `pkg/mythic/buildparameters.go:9`
  - Tests: `tests/integration/buildparameters_test.go:9`
  - Database: `buildparameter` table
  - Returns non-deleted parameter definitions sorted by payload type, then name
  - Defines what parameters are available when building payloads
  - Includes parameter schema, type, validation, default values

- **GetBuildParametersByPayloadType(payloadTypeID)** - Get build parameters for specific payload type
  - File: `pkg/mythic/buildparameters.go:68`
  - Tests: `tests/integration/buildparameters_test.go:78`
  - Database: `buildparameter` table with payload type filter
  - Validates payload type ID (non-zero)
  - Returns empty list for nonexistent payload types (not an error)
  - Sorted alphabetically by parameter name

- **GetBuildParameterInstances()** - Get all build parameter instances for current operation
  - File: `pkg/mythic/buildparameters.go:133`
  - Tests: `tests/integration/buildparameters_test.go:162`
  - Database: `buildparameterinstance` table
  - Returns actual parameter values used when creating payloads
  - Operation-scoped (only payloads in current operation)
  - Sorted by payload ID, then build parameter ID

- **GetBuildParameterInstancesByPayload(payloadID)** - Get parameter instances for specific payload
  - File: `pkg/mythic/buildparameters.go:179`
  - Tests: `tests/integration/buildparameters_test.go:212`
  - Database: `buildparameterinstance` table with payload filter
  - Validates payload ID (non-zero)
  - Includes nested BuildParameter details (name, type, description)
  - Returns empty list for nonexistent payloads (not an error)

**Build Parameter System:**

Build parameters define the configuration options available when building agent payloads. The system consists of two components:

1. **BuildParameterType** (definitions): Defines what parameters are available
   - Schema and validation rules for each parameter
   - Type information (String, Boolean, Number, ChooseOne, etc.)
   - Default values and randomization options
   - Requirement status and description

2. **BuildParameterInstance** (values): Stores actual values used in specific payloads
   - Links parameter definitions to specific payloads
   - Supports encrypted storage for sensitive parameters
   - Tracks when parameters were set

**BuildParameterType Fields:**

- **Identity**: ID, Name, PayloadTypeID
- **Type System**: ParameterType (String, Boolean, Number, ChooseOne, ChooseMultiple, File, Array, Date)
- **Validation**: VerifierRegex, Required, Parameter (JSON schema)
- **Defaults**: DefaultValue, Randomize, FormatString
- **Security**: IsCryptoType (for encryption keys, etc.)
- **Organization**: ParameterGroupName, Description
- **Status**: Deleted, CreationTime

**BuildParameterInstance Fields:**

- **Links**: ID, PayloadID, BuildParameterID
- **Values**: Value (plaintext), EncValue (encrypted value for sensitive params)
- **Metadata**: CreationTime
- **Relationships**: BuildParameter (nested definition), Payload (nested payload)

**Helper Methods - BuildParameterType:**

- **String()**: Human-readable representation
  - Format: "parameter_name (Type) (required)"
  - Indicates type and requirement status

- **IsRequired()**: Check if parameter is required
  - Returns Required field value
  - Required parameters must be provided when building payloads

- **IsCrypto()**: Check if parameter is a cryptographic type
  - Returns IsCryptoType field value
  - Crypto parameters (keys, secrets) should be handled securely

- **ShouldRandomize()**: Check if parameter should be auto-randomized
  - Returns Randomize field value
  - Randomized parameters get automatic values during build

- **IsDeleted()**: Check if parameter definition has been deleted
  - Returns Deleted field value
  - Deleted parameters no longer available for new payloads

**Helper Methods - BuildParameterInstance:**

- **String()**: Human-readable representation
  - Format: "parameter_name = value"
  - Shows parameter assignment for payload

- **IsEncrypted()**: Check if parameter value is encrypted
  - Returns true if EncValue is set and non-empty
  - Encrypted values stored separately for security

- **GetValue()**: Get the parameter value
  - Returns EncValue if encrypted, otherwise Value
  - Provides unified access to parameter value

**Usage Example:**

```go
// Get all build parameter definitions
parameters, err := client.GetBuildParameters(ctx)
if err != nil {
    return err
}

fmt.Printf("Found %d build parameter definitions\n", len(parameters))

// Group by payload type
paramsByType := make(map[int][]*types.BuildParameterType)
for _, param := range parameters {
    paramsByType[param.PayloadTypeID] = append(paramsByType[param.PayloadTypeID], param)
}

// Get parameters for a specific payload type
payloadTypeID := 1
typeParams, err := client.GetBuildParametersByPayloadType(ctx, payloadTypeID)
if err != nil {
    return err
}

fmt.Printf("Payload type %d has %d parameters:\n", payloadTypeID, len(typeParams))
for _, param := range typeParams {
    fmt.Printf("  %s\n", param.String())
    if param.IsRequired() {
        fmt.Printf("    Required\n")
    }
    if param.IsCrypto() {
        fmt.Printf("    Cryptographic parameter\n")
    }
    if param.ShouldRandomize() {
        fmt.Printf("    Will be randomized\n")
    }
    if param.DefaultValue != "" {
        fmt.Printf("    Default: %s\n", param.DefaultValue)
    }
}

// Get parameter instances for all payloads
instances, err := client.GetBuildParameterInstances(ctx)
if err != nil {
    return err
}

fmt.Printf("Found %d parameter instances across all payloads\n", len(instances))

// Get instances for a specific payload
payloadID := 5
payloadInstances, err := client.GetBuildParameterInstancesByPayload(ctx, payloadID)
if err != nil {
    return err
}

fmt.Printf("Payload %d has %d parameters configured:\n", payloadID, len(payloadInstances))
for _, inst := range payloadInstances {
    if inst.BuildParameter != nil {
        fmt.Printf("  %s = ", inst.BuildParameter.Name)
    }

    if inst.IsEncrypted() {
        fmt.Printf("[ENCRYPTED]\n")
    } else {
        fmt.Printf("%s\n", inst.GetValue())
    }
}
```

**Notes:**

- Build parameters are payload-type-specific; each payload type has its own parameters
- Parameter definitions are versioned with payload types
- Sensitive parameters (crypto keys, passwords) can be encrypted in storage
- Randomization allows automatic generation of unique identifiers
- Parameter groups organize related parameters (e.g., "C2 Configuration", "Obfuscation")
- VerifierRegex validates parameter values before payload build
- Parameter schema (in Parameter field) defines structure for complex parameters
- Default values used when parameter not explicitly set
- Deleted parameters remain in database but are filtered from queries

---

**File Browser:**

- **GetFileBrowserObjects()** - List all file browser objects for current operation
  - File: `pkg/mythic/filebrowser.go:10`
  - Tests: `tests/integration/filebrowser_test.go:11`
  - Database: `filebrowserobj` table
  - Returns non-deleted file/directory objects sorted by full path
  - Filters by current operation automatically
  - Includes files and directories from all callbacks in operation

- **GetFileBrowserObjectsByHost(host)** - Get file browser objects filtered by host
  - File: `pkg/mythic/filebrowser.go:83`
  - Tests: `tests/integration/filebrowser_test.go:83`
  - Database: `filebrowserobj` table with host filter
  - Validates host parameter (non-empty)
  - Returns empty list for nonexistent hosts (not an error)
  - Sorted by full path

- **GetFileBrowserObjectsByCallback(callbackID)** - Get file browser objects for specific callback
  - File: `pkg/mythic/filebrowser.go:164`
  - Tests: `tests/integration/filebrowser_test.go:132`
  - Database: `filebrowserobj` table with callback filter
  - Validates callback ID (non-zero)
  - Returns empty list for nonexistent callbacks (not an error)
  - Sorted by full path

**File Browser System:**

The file browser provides a unified, persistent interface for viewing and tracking files discovered through file browsing commands across all callbacks. Files are tracked with metadata including permissions, timestamps, and deletion status.

Key features:
- **Persistent tracking**: Files remain visible even after callbacks disconnect
- **Deletion tracking**: Files can be marked as deleted without removing history
- **Host-based organization**: Files grouped by hostname for multi-target operations
- **Callback association**: Track which callback discovered each file
- **Path normalization**: Full path construction from parent path and name

**FileBrowserObject Fields:**

- **Identity**: ID, Host, Name, ParentPath, FullPathText
- **Type**: IsFile (true for files, false for directories)
- **Metadata**: Permissions, Size, AccessTime, ModifyTime
- **Status**: Success (listing succeeded), Deleted (marked as deleted)
- **Context**: TaskID (command that discovered it), CallbackID, OperatorID, OperationID
- **Tracking**: Timestamp (when discovered), Comment, UpdateDeleted flag

**Helper Methods:**

- **FileBrowserObject.String()**: Human-readable representation
  - Format: "[file/dir] /full/path (deleted)"
  - Indicates type and deletion status

- **FileBrowserObject.IsDirectory()**: Check if object is a directory
  - Returns true if IsFile is false
  - Directories typically have Size=0

- **FileBrowserObject.IsDeleted()**: Check if object has been deleted
  - Returns Deleted field value
  - Deleted objects can still be queried (not removed from database)

- **FileBrowserObject.GetFullPath()**: Get the complete file path
  - Returns FullPathText if available
  - Otherwise constructs from ParentPath + "/" + Name
  - Handles root path and empty parent path correctly

**Usage Example:**

```go
// Get all file browser objects for current operation
objects, err := client.GetFileBrowserObjects(ctx)
if err != nil {
    return err
}

fmt.Printf("Found %d files and directories\n", len(objects))

// Organize by host
hostMap := make(map[string][]*types.FileBrowserObject)
for _, obj := range objects {
    hostMap[obj.Host] = append(hostMap[obj.Host], obj)
}

for host, hostObjects := range hostMap {
    fmt.Printf("\nHost: %s (%d items)\n", host, len(hostObjects))

    // Get objects for specific host using filter method
    hostObjs, err := client.GetFileBrowserObjectsByHost(ctx, host)
    if err != nil {
        return err
    }

    for _, obj := range hostObjs {
        fmt.Printf("  %s\n", obj.String())

        if obj.IsDirectory() {
            fmt.Printf("    Directory\n")
        } else {
            fmt.Printf("    File (%d bytes)\n", obj.Size)
        }

        if obj.IsDeleted() {
            fmt.Printf("    [DELETED]\n")
        }
    }
}

// Get file browser objects for a specific callback
callbackObjs, err := client.GetFileBrowserObjectsByCallback(ctx, callbackID)
if err != nil {
    return err
}

fmt.Printf("Callback %d has discovered %d files/directories\n",
    callbackID, len(callbackObjs))
```

**Notes:**

- File browser objects persist after callback disconnection for historical reference
- Objects are marked as deleted rather than removed, preserving operation history
- The Success field indicates whether the file listing command succeeded
- Permissions format is OS-specific (e.g., "rwxr-xr-x" on Unix, ACLs on Windows)
- AccessTime and ModifyTime are sourced from the target system
- File browser data is populated by file listing commands (ls, dir, etc.)
- Empty parent path is treated as root ("/")

Sources:
- [File Browser - Mythic](https://docs.mythic-c2.net/customizing/hooking-features/file-browser)
- [Mythic v3.2 Highlights: Interactive Tasking, Push C2, and Dynamic File Browser](https://specterops.io/blog/2023/11/29/mythic-v3-2-highlights-interactive-tasking-push-c2-and-dynamic-file-browser/)

---

**Commands:**

- **GetCommands()** - List all available commands from all payload types
  - File: `pkg/mythic/commands.go:10`
  - Tests: `tests/integration/commands_test.go:13`
  - Database: `command` table
  - Returns commands sorted alphabetically by command name
  - Each command includes: ID, name, payload type, version, support status, script-only flag
  - Helper methods: `IsSupported()`, `IsScriptOnly()`, `String()`

- **GetCommandParameters()** - Get all parameters for all commands
  - File: `pkg/mythic/commands.go:57`
  - Tests: `tests/integration/commands_test.go:73`
  - Database: `commandparameters` table
  - Returns parameters sorted by command ID
  - Parameter types: String, Boolean, Number, ChooseOne, ChooseMultiple, File, Array, Credential, LinkInfo
  - Supports static choices, dynamic queries, all commands, and loaded commands filtering
  - Helper methods: `IsRequired()`, `HasChoices()`, `IsDynamic()`, `String()`

- **GetLoadedCommands()** - Get commands loaded in a specific callback
  - File: `pkg/mythic/commands.go:110`
  - Tests: `tests/integration/commands_test.go:280`
  - Database: `loadedcommands` table
  - Requires callback ID (validates non-zero)
  - Returns commands loaded in callback, sorted by command name
  - Includes nested Command details (name, description, version)
  - Returns empty list for nonexistent callbacks (not an error)

**Command Type System:**

Commands in Mythic represent the available actions that can be executed through agents/payloads. The SDK provides three related types:

1. **Command**: Represents a command definition available in a payload type
   - Contains metadata: name, version, description, help text, author
   - Status flags: `Supported` (UI-supported), `ScriptOnly` (scripting-only)
   - MITRE ATT&CK mappings for threat intelligence correlation
   - Attributes (JSON) for custom metadata and filtering

2. **CommandParameter**: Defines parameters accepted by commands
   - Typed parameters with validation rules (required, default values)
   - Choice systems:
     - Static choices (JSON array)
     - All commands (reference other commands)
     - Loaded commands (only commands in current callback)
     - Dynamic query functions (runtime-generated choices)
   - Agent filtering (parameter support varies by agent type)
   - Build parameter filtering (parameter support varies by build config)

3. **LoadedCommand**: Tracks which commands are available in specific callbacks
   - Links callbacks to their loaded commands with version tracking
   - Operator tracking (who loaded the command)
   - Used for validating task submissions (only loaded commands can execute)
   - Enables filtered parameter choices (`choices_are_loaded_commands`)

**Command Parameters Design:**

Command parameters support sophisticated choice systems for building command arguments:

- **Static Choices**: Fixed list of values defined in parameter definition
  - Example: `["read", "write", "execute"]` for permission parameter
  - Parsed from `choices` field (JSON array)

- **All Commands**: Parameter accepts any command name from payload type
  - Enabled via `choices_are_all_commands` flag
  - Useful for command chaining (e.g., "execute this command after X")
  - Can be filtered by command attributes via `choice_filter_by_command_attributes`

- **Loaded Commands**: Parameter accepts only commands loaded in current callback
  - Enabled via `choices_are_loaded_commands` flag
  - Used for callback-specific operations (e.g., "unload command X")
  - Validates command availability at task submission time

- **Dynamic Queries**: Parameter choices generated at runtime via function
  - Specified in `dynamic_query_function` field
  - Examples: file browser (list files), process list, credential selector
  - Allows context-aware parameter building

**Helper Methods:**

The Commands API includes utility methods for common operations:

- **Command.String()**: Human-readable command representation
  - Format: "command_name vX (unsupported) (script-only)"
  - Indicates support status and execution restrictions

- **Command.IsSupported()**: Check if command is UI-supported
  - Returns `Supported` field value
  - Non-supported commands may still be usable via API/scripting

- **Command.IsScriptOnly()**: Check if command is script-only
  - Returns `ScriptOnly` field value
  - Script-only commands require special handling

- **CommandParameter.String()**: Human-readable parameter representation
  - Format: "parameter_name (Type) (required)"
  - Shows type and requirement status

- **CommandParameter.IsRequired()**: Check if parameter is required
  - Returns `Required` field value
  - Required parameters must be provided in task submission

- **CommandParameter.HasChoices()**: Check if parameter has predefined choices
  - Returns true if any choice system is configured
  - Includes static, all commands, or loaded commands

- **CommandParameter.IsDynamic()**: Check if parameter uses dynamic queries
  - Returns true if `DynamicQueryFunction` is set
  - Dynamic parameters fetch choices at runtime

- **LoadedCommand.String()**: Human-readable loaded command representation
  - Format: "command_name vX (Callback Y)"
  - Links command to specific callback

**Usage Example:**

```go
// Get all available commands
commands, err := client.GetCommands(ctx)
if err != nil {
    return err
}

// Find a specific command
for _, cmd := range commands {
    if cmd.Cmd == "download" && cmd.IsSupported() {
        fmt.Printf("Found: %s\n", cmd.String())

        // Get parameters for all commands
        params, err := client.GetCommandParameters(ctx)
        if err != nil {
            return err
        }

        // Filter to this command's parameters
        for _, param := range params {
            if param.CommandID == cmd.ID {
                fmt.Printf("  Parameter: %s\n", param.String())
                if param.IsRequired() {
                    fmt.Printf("    (required)\n")
                }
                if param.HasChoices() {
                    fmt.Printf("    (has predefined choices)\n")
                }
            }
        }
        break
    }
}

// Check what commands are loaded in a callback
loadedCmds, err := client.GetLoadedCommands(ctx, callbackID)
if err != nil {
    return err
}

fmt.Printf("Callback %d has %d loaded commands\n", callbackID, len(loadedCmds))
for _, lc := range loadedCmds {
    fmt.Printf("  %s\n", lc.String())
}
```

**Notes:**

- Commands are payload-type-specific; different payload types have different commands
- Not all commands are UI-supported; some are script-only or deprecated
- LoadedCommands validation ensures only available commands can be tasked
- Parameter choices enable rich command building UIs
- Dynamic parameters allow context-aware parameter selection
- Command versioning tracks command evolution across payload updates

---

**Container Management:**

- **ContainerListFiles(containerName, path)** - List files in a Docker container directory
  - File: `pkg/mythic/containers.go:27`
  - Tests: `tests/integration/containers_test.go:105`
  - GraphQL: `containerListFiles` query
  - Input: Container name (e.g., "mythic_athena", "http"), directory path
  - Returns: List of ContainerFileInfo with name, size, type (file/dir), permissions, mod time
  - Validates container name and path are non-empty
  - Useful for browsing payload type and C2 profile container filesystems during development

- **ContainerDownloadFile(containerName, path)** - Download a file from a Docker container
  - File: `pkg/mythic/containers.go:89`
  - Tests: `tests/integration/containers_test.go:54`
  - GraphQL: `containerDownloadFile` query
  - Input: Container name, file path within container
  - Returns: File content as bytes (base64 decoded)
  - Allows retrieving files from payload type and C2 profile containers for backup or analysis

- **ContainerWriteFile(containerName, path, content)** - Write a file to a Docker container
  - File: `pkg/mythic/containers.go:138`
  - Tests: `tests/integration/containers_test.go:121`
  - GraphQL: `containerWriteFile` mutation
  - Input: Container name, destination path, file content (base64 encoded)
  - Allows updating configuration files, adding scripts, or modifying containers during development
  - Content is automatically base64 encoded for transmission

- **ContainerRemoveFile(containerName, path)** - Remove a file from a Docker container
  - File: `pkg/mythic/containers.go:183`
  - Tests: `tests/integration/containers_test.go:121`
  - GraphQL: `containerRemoveFile` mutation
  - Input: Container name, file path to remove
  - Allows cleaning up temporary files or removing old configurations from containers

**Container Management System:**

The container management APIs provide direct file system access to Docker containers running payload types and C2 profiles. This is primarily useful for:
- **Development**: Modifying agent source code, build scripts, or configuration files
- **Debugging**: Retrieving log files or examining container state
- **Backup**: Downloading important files before container updates
- **Configuration**: Updating C2 profile settings or payload type parameters

**ContainerFileInfo Structure:**
- Name: File or directory name
- Size: File size in bytes (0 for directories)
- IsDir: Boolean indicating if entry is a directory
- ModTime: Last modification timestamp
- Permission: Unix-style permissions (e.g., "rwxr-xr-x")

**Helper Methods:**
- **ContainerFileInfo.String()**: Human-readable representation showing name, type, and size
- **ContainerFileInfo.IsDirectory()**: Returns true if entry is a directory
- **ContainerListFilesRequest.String()**: Display list request details
- **ContainerDownloadFileRequest.String()**: Display download request details
- **ContainerWriteFileRequest.String()**: Display write request with content size
- **ContainerRemoveFileRequest.String()**: Display remove request details

**Common Container Names:**
- `mythic_server`: Main Mythic server container
- `mythic_postgres`: PostgreSQL database container
- `mythic_rabbitmq`: RabbitMQ message broker container
- Payload type containers: `mythic_athena`, `poseidon`, `apfell`, etc.
- C2 profile containers: `http`, `websocket`, `smb`, `tcp`, etc.

**Common Use Cases:**
```go
// List files in a payload type container
files, err := client.ContainerListFiles(ctx, "mythic_athena", "/Mythic/agent_code")
for _, file := range files {
    if !file.IsDirectory() {
        fmt.Printf("Found source file: %s (%d bytes)\n", file.Name, file.Size)
    }
}

// Download a configuration file for backup
config, err := client.ContainerDownloadFile(ctx, "http", "/srv/config.json")
err = os.WriteFile("backup_config.json", config, 0644)

// Update a configuration file
newConfig, err := os.ReadFile("new_config.json")
err = client.ContainerWriteFile(ctx, "http", "/srv/config.json", newConfig)

// Clean up temporary files
err = client.ContainerRemoveFile(ctx, "mythic_athena", "/tmp/build_cache")
```

**Security Considerations:**
- Container operations require authenticated admin access to Mythic
- File operations execute with container user permissions
- Modifying container files can affect stability - use with caution
- Changes to containers may be lost on container restart unless persisted in volumes
- Downloading files may expose sensitive data - handle securely

**Notes:**
- File content is transferred using base64 encoding for binary safety
- Directory listing is not recursive - only immediate children are returned
- Path separators are Unix-style (/) regardless of host OS
- Empty files can be created by writing zero-length content
- Permissions format follows Unix conventions (rwxr-xr-x)
- ModTime format is ISO 8601 timestamp from container's perspective

---

**Proxy Operations:**

- **ToggleProxy(taskID, port, enable)** - Enable or disable a SOCKS proxy on a callback
  - File: `pkg/mythic/proxy.go:31`
  - Tests: `tests/integration/proxy_test.go:7`
  - GraphQL: `toggleProxy` mutation
  - Input: Task ID that started/will stop the proxy, port number (1-65535), enable flag (true/false)
  - Returns: ProxyInfo with proxy state details (ID, callback, port, active status, remote endpoint)
  - Validates task ID is positive and port is in valid range (1-65535)
  - Enables routing network traffic through compromised systems for lateral movement
  - Returns nil ProxyInfo when disabling (no active proxy state)

- **TestProxy(callbackID, port, targetURL)** - Test a SOCKS proxy connection
  - File: `pkg/mythic/proxy.go:111`
  - Tests: `tests/integration/proxy_test.go:58`
  - GraphQL: `testProxy` mutation
  - Input: Callback ID hosting the proxy, SOCKS port, target URL to test connectivity
  - Returns: TestProxyResponse with status, message, and error details
  - Validates callback ID, port range, and non-empty target URL
  - Tests connectivity by attempting to reach target URL through the proxy
  - Useful for validating proxy functionality before operational use

**Proxy System:**

SOCKS proxies in Mythic enable operators to route network traffic through compromised systems, providing access to internal networks and bypassing network segmentation. The proxy operations support:
- **Lateral Movement**: Access internal systems not directly reachable from the operator's network
- **Pivoting**: Chain proxies through multiple compromised systems
- **Network Reconnaissance**: Scan and interact with internal network resources
- **Operational Security**: Hide operator traffic by routing through victim infrastructure

**ProxyInfo Structure:**
- ID: Unique proxy identifier
- CallbackID: Callback hosting the SOCKS proxy
- Port: SOCKS proxy port number
- PortType: Proxy type (typically "socks")
- Active: Boolean indicating if proxy is currently running
- LocalPort: Local port on the Mythic server for proxy forwarding
- RemoteIP: Remote IP address being proxied to
- RemotePort: Remote port being proxied to
- OperationID: Associated operation
- Deleted: Soft delete flag
- ProxyCallback: Optional callback ID that provides the proxy (for chained proxies)

**Helper Methods:**
- **ProxyInfo.String()**: Human-readable representation showing type, port, and status
- **ProxyInfo.IsActive()**: Returns true if proxy is active and not deleted
- **ProxyInfo.IsDeleted()**: Returns true if proxy has been marked as deleted
- **ToggleProxyRequest.String()**: Display toggle request details
- **TestProxyRequest.String()**: Display test request details
- **TestProxyResponse.String()**: Display test result
- **TestProxyResponse.IsSuccessful()**: Returns true if proxy test succeeded

**Common Use Cases:**
```go
// Start a SOCKS proxy on callback
taskID := 123 // Task that starts the proxy
proxy, err := client.ToggleProxy(ctx, taskID, 1080, true)
if err != nil {
    return err
}
fmt.Printf("Proxy started: %s\n", proxy.String())
fmt.Printf("Forward to: %s:%d\n", proxy.RemoteIP, proxy.RemotePort)

// Test the proxy connection
result, err := client.TestProxy(ctx, proxy.CallbackID, 1080, "https://www.google.com")
if err != nil {
    return err
}
if result.IsSuccessful() {
    fmt.Printf("Proxy is working: %s\n", result.Message)
} else {
    fmt.Printf("Proxy test failed: %s\n", result.Error)
}

// Test internal network connectivity through proxy
internalResult, err := client.TestProxy(ctx, proxy.CallbackID, 1080, "http://10.0.0.1")
if internalResult.IsSuccessful() {
    fmt.Println("Can reach internal network through proxy")
}

// Stop the proxy when done
proxy, err = client.ToggleProxy(ctx, taskID, 1080, false)
if err != nil {
    return err
}
fmt.Println("Proxy stopped")
```

**Proxy Workflow:**
1. **Issue SOCKS command**: Task callback to start SOCKS proxy listener
2. **Wait for task completion**: Proxy starts on specified port
3. **Toggle proxy**: Call ToggleProxy to register the proxy with Mythic
4. **Test connectivity**: Use TestProxy to validate proxy is working
5. **Configure tools**: Set SOCKS proxy settings in tools (proxychains, browsers, etc.)
6. **Operational use**: Route traffic through the proxy for lateral movement
7. **Stop proxy**: Toggle proxy off when finished or callback dies

**Port Selection:**
- **Common SOCKS ports**: 1080, 9050, 8080
- **Avoid conflicts**: Check for existing services on target system
- **Privileged ports**: Ports 1-1023 may require elevated privileges on target
- **Firewall considerations**: Choose ports allowed by target's firewall rules

**Security Considerations:**
- SOCKS proxies expose compromised systems to network traffic routing
- Test proxies before operational use to avoid alerting defenders
- Monitor proxy performance - slow proxies may indicate detection or issues
- Clean up proxies when finished to reduce detection surface
- Consider network flow impact - heavy proxy traffic may trigger alerts
- Proxy chains may be logged by intermediate systems

**Limitations:**
- Proxies require active callback to function
- If callback dies, proxy stops working
- Network latency increases with proxy hops
- Some protocols may not work well through SOCKS proxies
- Target system must have network access to desired destinations
- Firewall rules on target may block proxy traffic

**Notes:**
- ToggleProxy requires a task ID from a command that starts/stops the proxy
- The task must complete successfully before proxy becomes active
- TestProxy performs actual connectivity test - not just a status check
- Proxy port must not conflict with existing services on target system
- LocalPort is assigned by Mythic for server-side proxy forwarding
- ProxyCallback field supports chained proxies (proxy through another proxy)
- Active proxies are automatically cleaned up when callbacks exit

---

**Utility Functions:**

- **CreateRandom(format, length)** - Generate a random string based on a format specification
  - File: `pkg/mythic/utility.go:20`
  - Tests: `tests/integration/utility_test.go:12`
  - GraphQL: `createRandom` mutation
  - Input: Format string with specifiers (%s, %S, %d, %x, %X), optional length
  - Returns: CreateRandomResponse with generated random string
  - Validates format is non-empty
  - Useful for generating random identifiers, callback IDs, payload names, or test data

  Format specifiers:
  - `%s`: Random lowercase letters (a-z)
  - `%S`: Random uppercase letters (A-Z)
  - `%d`: Random digits (0-9)
  - `%x`: Random lowercase hexadecimal (0-9, a-f)
  - `%X`: Random uppercase hexadecimal (0-9, A-F)
  - Literal characters preserved (e.g., "callback-%s" → "callback-xyzab")

- **ConfigCheck()** - Check Mythic configuration validity and status
  - File: `pkg/mythic/utility.go:93`
  - Tests: `tests/integration/utility_test.go:117`
  - GraphQL: `config_check` query
  - Returns: ConfigCheckResponse with validation results, errors, and config details
  - No input parameters required
  - Validates database, RabbitMQ, Redis, containers, environment variables, and permissions
  - Useful for debugging configuration issues and validating setup before operations

**CreateRandomResponse Structure:**
- Status: Operation status ("success" or "error")
- RandomString: Generated random string
- Error: Error message if generation failed

**ConfigCheckResponse Structure:**
- Status: Operation status
- Valid: Boolean indicating if configuration is valid
- Errors: List of configuration error messages
- Config: Map of configuration details (database, services, containers)
- Message: Human-readable status message

**Helper Methods:**
- **CreateRandomRequest.String()**: Display request details
- **CreateRandomResponse.String()**: Display generated string or error
- **CreateRandomResponse.IsSuccessful()**: Returns true if generation succeeded
- **ConfigCheckResponse.String()**: Display validation result
- **ConfigCheckResponse.IsValid()**: Returns true if configuration is valid with no errors
- **ConfigCheckResponse.HasErrors()**: Returns true if there are configuration errors
- **ConfigCheckResponse.GetErrors()**: Returns list of error messages

**Common Use Cases:**
```go
// Generate random callback ID
result, err := client.CreateRandom(ctx, "callback-%s-%d", 8)
if err != nil {
    return err
}
callbackID := result.RandomString
fmt.Printf("Generated callback ID: %s\n", callbackID)

// Generate random hex payload name
result, err = client.CreateRandom(ctx, "payload_%x", 16)
if err != nil {
    return err
}
payloadName := result.RandomString

// Generate random test data
testData, err := client.CreateRandom(ctx, "%S%d%s", 6)
if testData.IsSuccessful() {
    fmt.Printf("Test data: %s\n", testData.RandomString)
}

// Check configuration before starting operations
config, err := client.ConfigCheck(ctx)
if err != nil {
    return err
}

if !config.IsValid() {
    fmt.Printf("Configuration is invalid:\n")
    for _, err := range config.GetErrors() {
        fmt.Printf("  - %s\n", err)
    }
    return fmt.Errorf("fix configuration errors before proceeding")
}

fmt.Println("Configuration is valid, proceeding with operations")

// Log configuration details
if len(config.Config) > 0 {
    fmt.Println("Configuration details:")
    for key, value := range config.Config {
        fmt.Printf("  %s: %v\n", key, value)
    }
}
```

**CreateRandom Use Cases:**
- **Callback/Agent IDs**: Generate unique identifiers for callbacks and agents
- **Payload Names**: Create random payload filenames to avoid detection
- **Test Data**: Generate random strings for testing and development
- **Session IDs**: Create random session identifiers for C2 communication
- **Obfuscation**: Generate random strings for variable names, function names
- **Operation Names**: Create random but memorable operation codenames

**ConfigCheck Use Cases:**
- **Pre-flight Checks**: Validate configuration before starting operations
- **Troubleshooting**: Identify configuration issues when things aren't working
- **Health Monitoring**: Periodic checks to ensure all services are connected
- **Setup Validation**: Verify new Mythic installations are configured correctly
- **Container Status**: Check that all required containers are running
- **Database Connectivity**: Verify PostgreSQL connection is working
- **Message Queue**: Ensure RabbitMQ is accessible for task distribution

**Configuration Items Checked:**
- **Database (PostgreSQL)**: Connection status, schema version
- **RabbitMQ**: Message queue connectivity for task distribution
- **Redis** (if configured): Cache and session storage connectivity
- **Containers**: Status of payload type and C2 profile containers
- **Environment Variables**: Required variables are set correctly
- **Permissions**: File system and Docker permissions are adequate
- **Network**: Internal service networking is functioning

**Random String Patterns:**
```go
// Lowercase identifier: "abc123xyz"
CreateRandom(ctx, "%s%d%s", 3)

// Uppercase codename: "ALPHA123"
CreateRandom(ctx, "%S%d", 5)

// Hex UUID-style: "a1b2c3d4e5f6"
CreateRandom(ctx, "%x", 12)

// Mixed case: "Test123ABC"
CreateRandom(ctx, "%S%d%S", 4)

// With separators: "callback-abc-123"
CreateRandom(ctx, "callback-%s-%d", 3)

// Email-style: "user@random.com"
CreateRandom(ctx, "%s@%s.com", 6)
```

**Security Considerations:**
- **CreateRandom**: Generated strings are pseudo-random, not cryptographically secure
- Use CreateRandom for identifiers and obfuscation, not cryptographic keys
- Random strings may not have sufficient entropy for security-critical uses
- Consider length requirements - longer strings are more unique
- **ConfigCheck**: May expose configuration details - use in secure environments
- Configuration errors may reveal internal infrastructure details
- Limit ConfigCheck calls to authenticated, trusted operators

**Notes:**
- CreateRandom length parameter determines characters per format specifier
- Zero length uses Mythic's default length per specifier
- Format string can mix specifiers and literal characters
- ConfigCheck requires authentication but no special permissions
- Configuration details returned in Config map vary by Mythic version
- Failed configuration checks don't prevent API usage - they're informational
- Some configuration issues may only appear under specific conditions

---

**Block Lists:**

- **DeleteBlockList(blockListID)** - Delete a block list and all its entries
  - File: `pkg/mythic/blocklist.go:24`
  - Tests: `tests/integration/blocklist_test.go:11`
  - GraphQL: `deleteBlockList` mutation
  - Input: Block list ID to delete
  - Returns: DeleteBlockListResponse with status and message
  - Validates block list ID is positive
  - Removes the entire block list including all associated entries
  - Used when a block list is no longer needed or was created in error

- **DeleteBlockListEntry(entryIDs)** - Delete specific entries from block lists
  - File: `pkg/mythic/blocklist.go:75`
  - Tests: `tests/integration/blocklist_test.go:36`
  - GraphQL: `deleteBlockListEntry` mutation
  - Input: List of entry IDs to delete (must be non-empty, all positive, no duplicates)
  - Returns: DeleteBlockListEntryResponse with status and count of deleted entries
  - Validates entry IDs list is non-empty, all IDs are positive, no duplicates
  - Allows removing specific IPs, domains, or user agents without deleting entire list
  - More efficient than deleting entries one at a time

**Block List System:**

Block lists in Mythic prevent specific IP addresses, domains, or user agents from accessing C2 infrastructure. This security feature helps:
- **Avoid Detection**: Block known security scanner IPs and user agents
- **Prevent Analysis**: Stop automated malware analysis systems from connecting
- **Operational Security**: Block IPs associated with security vendors and researchers
- **Incident Response**: Quickly block IPs that show signs of defensive analysis

**BlockList Structure:**
- ID: Unique block list identifier
- Name: Descriptive name for the block list (e.g., "Security Scanners", "Sandboxes")
- Description: Detailed description of what the list blocks
- OperationID: Associated operation
- Active: Boolean indicating if list is actively enforced
- Deleted: Soft delete flag

**BlockListEntry Structure:**
- ID: Unique entry identifier
- BlockListID: Parent block list ID
- Type: Entry type - "ip", "domain", or "user_agent"
- Value: The actual value to block (IP address, domain name, or user agent string)
- Description: Optional description of why this entry is blocked
- Active: Boolean indicating if entry is actively enforced
- Deleted: Soft delete flag

**Helper Methods:**
- **BlockList.String()**: Human-readable representation showing name and status
- **BlockList.IsActive()**: Returns true if list is active and not deleted
- **BlockList.IsDeleted()**: Returns true if list has been deleted
- **BlockListEntry.String()**: Human-readable representation showing type, value, and status
- **BlockListEntry.IsActive()**: Returns true if entry is active and not deleted
- **BlockListEntry.IsDeleted()**: Returns true if entry has been deleted
- **DeleteBlockListRequest.String()**: Display deletion request details
- **DeleteBlockListResponse.String()**: Display deletion result
- **DeleteBlockListResponse.IsSuccessful()**: Returns true if deletion succeeded
- **DeleteBlockListEntryRequest.String()**: Display entry deletion request
- **DeleteBlockListEntryResponse.String()**: Display entry deletion result with count
- **DeleteBlockListEntryResponse.IsSuccessful()**: Returns true if entry deletion succeeded

**Common Use Cases:**
```go
// Delete an outdated block list
result, err := client.DeleteBlockList(ctx, blockListID)
if err != nil {
    return err
}
if result.IsSuccessful() {
    fmt.Println(result.String())
}

// Remove specific entries from a block list
// For example, after confirming an IP is no longer a threat
entryIDs := []int{123, 456, 789}
result, err := client.DeleteBlockListEntry(ctx, entryIDs)
if err != nil {
    return err
}
fmt.Printf("Deleted %d entries\n", result.DeletedCount)

// Clean up multiple entries at once
suspiciousIPs := []int{101, 102, 103, 104, 105}
result, err = client.DeleteBlockListEntry(ctx, suspiciousIPs)
if result.IsSuccessful() {
    fmt.Printf("Removed %d IPs from block list\n", result.DeletedCount)
}
```

**Entry Types:**
- **ip**: Block specific IP addresses
  - Examples: "192.168.1.100", "10.0.0.50", "203.0.113.42"
  - Prevents connections from these IPs to C2 infrastructure
  - Useful for blocking known security vendor IPs

- **domain**: Block specific domain names
  - Examples: "scanner.security.com", "sandbox.analysis.net"
  - Blocks reverse DNS lookups or SNI requests matching these domains
  - Prevents security tools with known domains from analyzing payloads

- **user_agent**: Block specific user agent strings
  - Examples: "SecurityScanner/1.0", "wget/1.21", "python-requests/2.28"
  - Blocks HTTP/HTTPS requests with matching User-Agent headers
  - Prevents automated tools and scripts from downloading payloads

**Block List Management Workflow:**
1. **Create Block Lists**: Set up lists for different categories (scanners, sandboxes, etc.)
2. **Add Entries**: Populate with known security IPs, domains, and user agents
3. **Activate Lists**: Enable lists to start blocking matching traffic
4. **Monitor Traffic**: Watch for blocked connection attempts in logs
5. **Update Entries**: Add new threats, remove false positives
6. **Delete Obsolete Entries**: Remove entries that are no longer threats
7. **Delete Old Lists**: Clean up block lists no longer needed

**Deletion Strategies:**
- **Individual Entry Deletion**: Remove specific false positives while keeping list
- **Bulk Entry Deletion**: Clean up multiple obsolete entries efficiently
- **Full List Deletion**: Remove entire category when no longer relevant
- **Soft Deletes**: Mythic soft-deletes entries for audit trail (deleted flag)

**Security Considerations:**
- Block lists protect C2 infrastructure from analysis and detection
- Over-aggressive blocking may block legitimate target networks
- False positives in block lists can prevent valid callbacks
- Block list changes take effect immediately on active C2 profiles
- Deleted block lists may still be referenced in logs for audit purposes
- Consider using both IP and user agent blocking for defense in depth
- Regularly review and update block lists as threat landscape changes

**Common Block List Entries:**
- **Security Vendors**: IPs owned by Palo Alto, Cisco, CrowdStrike, etc.
- **Malware Analysis**: VirusTotal, Any.run, Joe Sandbox, Hybrid Analysis
- **Research Organizations**: University security labs, threat intelligence firms
- **Automated Scanners**: Shodan, Censys, ZoomEye, BinaryEdge
- **Cloud Providers**: AWS security scanning, Azure security center
- **User Agents**: wget, curl, python-requests, automated security tools

**Notes:**
- DeleteBlockList removes both the list and all its entries permanently
- DeleteBlockListEntry requires at least one entry ID
- All entry IDs must be positive integers
- Duplicate entry IDs in the same request are not allowed
- Deletion is immediate and affects active C2 profile filtering
- Block list IDs and entry IDs are unique within the Mythic instance
- Soft-deleted entries remain in database for audit but are not enforced
- Block lists are operation-specific - deleting one doesn't affect others

---

### ⏳ Pending (0/20)

**Dynamic Queries:**
- **DynamicQueryFunction()** - Dynamic parameter queries
  - GraphQL: `dynamic_query_function` mutation

- **DynamicQueryBuildParameter()** - Build parameter queries
  - GraphQL: `dynamicQueryBuildParameterFunction` mutation

- **TypedarrayParseFunction()** - Parse typed arrays
  - GraphQL: `typedarray_parse_function` mutation

**Staging:**
- **GetStagingInfo()** - Get payload staging info
  - Database: `staginginfo` table

---

## Implementation Priority Recommendations

### High Priority (Core Functionality)
1. ✅ **Operations Management** - Essential for multi-operation environments
2. ✅ **Payloads** - Critical for agent deployment
3. ✅ **Credentials** - Important for tracking compromised accounts
4. ✅ **C2 Profiles** - Needed for agent communication management
5. ✅ **Processes** - Important for situational awareness

### Medium Priority (Enhanced Features)
6. ✅ **Artifacts/IOCs** - Useful for tracking indicators
7. ✅ **Tags** - Organization and categorization
8. ✅ **Keylogs** - Credential harvesting operations
9. ✅ **MITRE ATT&CK** - Threat intelligence integration
10. ✅ **Reporting** - Operation documentation

### Low Priority (Advanced Features)
11. **Eventing/Workflows** - Automation for advanced users
12. ✅ **Browser Scripts** - Custom UI functionality
13. ✅ **Container Management** - Development/debugging
14. **Dynamic Queries** - Advanced parameter handling
15. ✅ **Proxy Operations** - Specialized networking

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

*Last updated: 2026-01-08*
*SDK Version: In Development*
*Mythic API Version: v3.4.x*
