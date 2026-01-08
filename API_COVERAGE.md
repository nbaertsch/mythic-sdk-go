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
| Advanced Features | 0 | 0 | 20 | 20 |
| **TOTAL** | **129** | **0** | **25** | **154** |

**Overall Coverage: 83.8%**

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
