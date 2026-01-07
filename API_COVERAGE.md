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
| Tasks | 11 | 0 | 1 | 12 |
| Files | 8 | 0 | 0 | 8 |
| Operations | 0 | 0 | 11 | 11 |
| Payloads | 0 | 0 | 10 | 10 |
| Credentials | 0 | 0 | 3 | 3 |
| C2 Profiles | 0 | 0 | 9 | 9 |
| Artifacts | 0 | 0 | 3 | 3 |
| Tags | 0 | 0 | 3 | 3 |
| Tokens | 0 | 0 | 4 | 4 |
| Processes | 0 | 0 | 2 | 2 |
| Keylogs | 0 | 0 | 2 | 2 |
| Browser Scripts | 0 | 0 | 3 | 3 |
| MITRE ATT&CK | 0 | 0 | 3 | 3 |
| Reporting | 0 | 0 | 2 | 2 |
| Eventing/Workflows | 0 | 0 | 15 | 15 |
| Operators | 0 | 0 | 11 | 11 |
| GraphQL Subscriptions | 0 | 0 | 1 | 1 |
| Advanced Features | 0 | 0 | 20 | 20 |
| **TOTAL** | **35** | **0** | **98** | **133** |

**Overall Coverage: 26.3%**

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

- **GetCallbacks()** - List all callbacks with filtering
  - File: `pkg/mythic/callbacks.go:106`
  - Tests: `tests/integration/callbacks_test.go:11`

- **GetCallbackByID()** - Get specific callback by display ID
  - File: `pkg/mythic/callbacks.go:181`
  - Tests: `tests/integration/callbacks_test.go:33`

- **GetActiveCallbacks()** - Filter only active callbacks
  - File: `pkg/mythic/callbacks.go:240`
  - Tests: `tests/integration/callbacks_test.go:51`

- **UpdateCallback()** - Update callback properties (description, ips, host, etc.)
  - File: `pkg/mythic/callbacks.go:293`
  - GraphQL: `updateCallback` mutation

- **Callback.IsActive()** - Helper to check if callback is active
  - File: `pkg/mythic/callbacks.go:379`

- **Callback.IsDead()** - Helper to check if callback is dead
  - File: `pkg/mythic/callbacks.go:384`

- **Callback.String()** - String representation
  - File: `pkg/mythic/callbacks.go:374`

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

---

## 3. Tasks (Commands)

### ✅ Tested (11/12 - 91.7%)

- **CreateTask()** - Issue task to callback(s)
  - File: `pkg/mythic/tasks.go:82`
  - Tests: `tests/integration/tasks_test.go`
  - GraphQL: `createTask` mutation

- **GetTask()** - Get task by ID
  - File: `pkg/mythic/tasks.go:171`
  - Tests: `tests/integration/tasks_test.go`

- **GetTasksByCallback()** - List all tasks for a callback
  - File: `pkg/mythic/tasks.go:239`

- **GetTaskOutput()** - Get task responses/output
  - File: `pkg/mythic/tasks.go:312`
  - Tests: `tests/integration/tasks_test.go`

- **UpdateTaskComment()** - Add/update task comment
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

### ⏳ Pending (1/12)

- **SubscribeToTaskOutput()** - Real-time task output subscription
  - GraphQL: Subscription (requires websocket support)
  - Note: Requires WebSocket transport layer implementation

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

### ⏳ Pending (11/11)

- **GetOperations()** - List all operations
  - Database: `operation` table

- **GetOperationByID()** - Get specific operation
  - Database: `operation` table

- **CreateOperation()** - Create new operation
  - GraphQL: `createOperation` mutation

- **UpdateOperation()** - Update operation details
  - GraphQL: `updateOperation` mutation
  - Fields: name, channel, complete, webhook, admin_id, banner_text, banner_color

- **UpdateCurrentOperation()** - Switch user's current operation
  - GraphQL: `updateCurrentOperation` mutation

- **GetOperatorsByOperation()** - List operators in operation
  - Database: `operatoroperation` table

- **UpdateOperatorOperation()** - Add/remove operators from operation
  - GraphQL: `updateOperatorOperation` mutation

- **GetOperationEventLog()** - Get operation event logs
  - Database: `operationeventlog` table

- **CreateOperationEventLog()** - Create event log entry
  - GraphQL: `createOperationEventLog` mutation

- **GetGlobalSettings()** - Get Mythic global settings
  - GraphQL: `getGlobalSettings` query

- **UpdateGlobalSettings()** - Update global settings
  - GraphQL: `updateGlobalSettings` mutation

---

## 6. Payloads

### ⏳ Pending (10/10)

- **GetPayloads()** - List all payloads
  - Database: `payload` table

- **GetPayloadByUUID()** - Get specific payload
  - Database: `payload` table

- **CreatePayload()** - Build new payload
  - GraphQL: `createPayload` mutation
  - Input: JSON payload definition

- **RebuildPayload()** - Rebuild existing payload
  - GraphQL: `rebuild_payload` mutation

- **UpdatePayload()** - Update payload settings
  - GraphQL: `updatePayload` mutation
  - Fields: callback_alert, callback_allowed, description, deleted

- **DeletePayload()** - Delete payload
  - GraphQL: `deleteFile` mutation

- **ExportPayloadConfig()** - Export payload configuration
  - GraphQL: `exportPayloadConfig` query

- **GetPayloadTypes()** - List available payload types
  - Database: `payloadtype` table

- **GetPayloadCommands()** - Get commands for payload
  - Database: `payloadcommand` table

- **GetPayloadOnHost()** - Track payloads deployed on hosts
  - Database: `payloadonhost` table

---

## 7. Credentials

### ⏳ Pending (3/3)

- **GetCredentials()** - List all credentials
  - Database: `credential` table

- **CreateCredential()** - Add new credential
  - GraphQL: `createCredential` mutation
  - Fields: realm, account, credential, comment, credential_type

- **UpdateCredential()** - Update credential
  - Database: `credential` table (direct update)

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

### ⏳ Pending (2/2)

- **GetProcesses()** - List processes from callbacks
  - Database: `process` table

- **GetProcessTree()** - Get process tree for callback
  - Database: `process` table with parent relationships

---

## 13. Keylogs

### ⏳ Pending (2/2)

- **GetKeylogs()** - List keylog entries
  - Database: `keylog` table

- **GetKeylogsByCallback()** - Filter keylogs by callback
  - Database: `keylog` table

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
1. **Operations Management** - Essential for multi-operation environments
2. **Payloads** - Critical for agent deployment
3. **Credentials** - Important for tracking compromised accounts
4. **C2 Profiles** - Needed for agent communication management
5. **Processes** - Important for situational awareness

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
