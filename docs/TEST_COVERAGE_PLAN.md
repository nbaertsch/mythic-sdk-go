# COMPREHENSIVE INTEGRATION TEST COVERAGE PLAN
## Mythic SDK Go - Full Code Path Testing

**Date**: 2026-01-20
**Status**: CRITICAL - Schema bugs found in production
**Total SDK Methods**: 232
**Currently Tested**: 11 (~4.7%)
**Target Coverage**: 100%

---

## EXECUTIVE SUMMARY

The Mythic Go SDK has **221 untested code paths** that led to production failures:
- GraphQL schema mismatches (`payloadtype_id` vs `payload_type_id`)
- Type mismatches (bool vs array fields)
- Missing required parameters (C2 `callback_host`)
- Field removals (`attack`, `attributes`, `supported_ui_features`)

These bugs were NOT caught by CI because:
1. Tests exist but don't validate response data structures
2. Tests may use mock data or old Mythic versions
3. Integration tests don't fully exercise all GraphQL queries
4. No schema validation layer

**This document provides a COMPLETE test coverage plan for all 232 SDK methods.**

---

## PART 1: CURRENT CI/CD ANALYSIS

### Current CI Setup (.github/workflows/integration.yml)

**What CI Does Right:**
- ✅ Installs latest Mythic from GitHub
- ✅ Installs Poseidon agent
- ✅ Installs HTTP C2 profile
- ✅ Runs 4 test phases in parallel
- ✅ Waits for containers to be ready

**Critical CI Gaps:**
- ❌ Uses `latest` Mythic (schema can change)
- ❌ Tests may pass on old schema but fail on new
- ❌ No GraphQL schema validation
- ❌ No field-level data validation
- ❌ Tests check for errors but not response correctness

### CI Test Phases

**Phase 1+2: Core APIs**
```regex
TestE2E_Auth|TestE2E_Operations|TestE2E_File|TestE2E_Credential|
TestE2E_Artifact|TestE2E_Tag|TestE2E_Operator|TestE2E_C2Profile|
TestE2E_MITREAttack|TestE2E_Payload
```

**Phase 3: Agent Tests**
```regex
TestE2E_Callback
```

**Phase 4: Advanced APIs**
```regex
TestE2E_Command|TestE2E_Task|TestE2E_Response|TestE2E_Process|
TestE2E_Screenshot|TestE2E_Report|TestE2E_Redirect|TestE2E_Token
```

**Phase 5: Edge Cases**
```regex
TestE2E_Context|TestE2E_Concurrent|TestE2E_Large|TestE2E_Rate|
TestE2E_Invalid|TestE2E_Network|TestE2E_Memory|TestE2E_Query
```

---

## PART 2: SCHEMA BUGS FOUND

### Bug Category 1: Field Name Mismatches

| SDK Field | Actual Schema | Files Affected | Impact |
|-----------|---------------|----------------|--------|
| `payloadtype_id` | `payload_type_id` | commands.go, buildparameters.go | **CRITICAL** - All command queries fail |
| `attack` | (doesn't exist) | commands.go | **HIGH** - Query fails |

### Bug Category 2: Type Mismatches

| SDK Field | SDK Type | Actual Type | Files Affected | Impact |
|-----------|----------|-------------|----------------|--------|
| `callback_alert` | `bool` | `array` | payloads.go | **CRITICAL** - Panic on GetPayloadByUUID |
| `auto_generated` | `bool` | `array` | payloads.go | **CRITICAL** - Panic on GetPayloadByUUID |
| `supported_ui_features` | `bool` | `array` | commands.go | **CRITICAL** - Panic on GetCommands |
| `attributes` | `string` | `jsonb` | commands.go | **HIGH** - Type conversion fails |

### Bug Category 3: Missing Required Parameters

| Operation | Required Param | Default Missing | Impact |
|-----------|----------------|-----------------|--------|
| CreatePayload | `callback_host` in C2 config | Yes | **CRITICAL** - "Missing C2 profile information" |

### Root Cause Analysis

**Why Tests Passed But Production Failed:**
1. Tests call `GetCommands()` but don't validate `PayloadTypeID` field is populated correctly
2. Tests check `err != nil` but not actual GraphQL field names
3. CI may use cached/older Mythic version
4. No test validates GraphQL introspection matches SDK expectations

---

## PART 3: COMPLETE CODE PATH INVENTORY

### Summary Statistics

| Category | Total Methods | Tested | Untested | Coverage % |
|----------|---------------|--------|----------|------------|
| **Authentication** | 8 | 1 | 7 | 12.5% |
| **Operations** | 11 | 2 | 9 | 18.2% |
| **Callbacks** | 10 | 1 | 9 | 10.0% |
| **Tasks** | 18 | 1 | 17 | 5.6% |
| **Commands** | 7 | 0 | 7 | 0.0% |
| **Payloads** | 12 | 3 | 9 | 25.0% |
| **Files** | 11 | 1 | 10 | 9.1% |
| **Credentials** | 6 | 1 | 5 | 16.7% |
| **C2 Profiles** | 9 | 1 | 8 | 11.1% |
| **Alerts** | 7 | 0 | 7 | 0.0% |
| **Artifacts** | 8 | 0 | 8 | 0.0% |
| **Attack (MITRE)** | 6 | 0 | 6 | 0.0% |
| **Build Parameters** | 4 | 0 | 4 | 0.0% |
| **Browser Scripts** | 3 | 0 | 3 | 0.0% |
| **Containers** | 4 | 0 | 4 | 0.0% |
| **Dynamic Query** | 3 | 0 | 3 | 0.0% |
| **Eventing** | 15 | 0 | 15 | 0.0% |
| **File Browser** | 3 | 0 | 3 | 0.0% |
| **Hosts** | 5 | 0 | 5 | 0.0% |
| **Keylogs** | 3 | 0 | 3 | 0.0% |
| **Operators** | 11 | 0 | 11 | 0.0% |
| **Processes** | 5 | 0 | 5 | 0.0% |
| **Proxy** | 2 | 0 | 2 | 0.0% |
| **Reporting** | 2 | 0 | 2 | 0.0% |
| **Responses** | 6 | 0 | 6 | 0.0% |
| **RPFWD** | 4 | 0 | 4 | 0.0% |
| **Screenshots** | 6 | 0 | 6 | 0.0% |
| **Subscriptions** | 2 | 0 | 2 | 0.0% |
| **Tags** | 11 | 0 | 11 | 0.0% |
| **Tokens** | 7 | 0 | 7 | 0.0% |
| **Utility** | 2 | 0 | 2 | 0.0% |
| **Blocklist** | 2 | 0 | 2 | 0.0% |
| **Staging** | 1 | 0 | 1 | 0.0% |
| **TOTAL** | **232** | **11** | **221** | **4.7%** |

---

## PART 4: PRIORITY-BASED TEST PLAN

### PRIORITY 0: CRITICAL SCHEMA VALIDATION (NEW)

**Must be added FIRST before any other tests:**

#### Test: `TestE2E_SchemaValidation`
**Purpose**: Validate GraphQL schema matches SDK expectations
**Rationale**: Prevents all schema mismatch bugs

```go
func TestE2E_SchemaValidation(t *testing.T) {
    client := AuthenticateTestClient(t)

    // Test 1: Validate command schema
    t.Run("ValidateCommandSchema", func(t *testing.T) {
        // Query GraphQL introspection for 'command' type
        schema := QuerySchema(t, client, "command")

        // Assert required fields exist with correct types
        AssertFieldExists(t, schema, "id", "Int")
        AssertFieldExists(t, schema, "cmd", "String")
        AssertFieldExists(t, schema, "payload_type_id", "Int") // NOT payloadtype_id
        AssertFieldExists(t, schema, "description", "String")
        AssertFieldExists(t, schema, "help_cmd", "String")
        AssertFieldExists(t, schema, "version", "Int")
        AssertFieldExists(t, schema, "author", "String")
        AssertFieldExists(t, schema, "script_only", "Boolean")

        // Assert removed fields DON'T exist
        AssertFieldNotExists(t, schema, "attack") // Moved to attackcommands relation

        // Assert type mismatches
        AssertFieldType(t, schema, "supported_ui_features", "array") // NOT bool
        AssertFieldType(t, schema, "attributes", "jsonb") // NOT string
    })

    // Test 2: Validate payload schema
    t.Run("ValidatePayloadSchema", func(t *testing.T) {
        schema := QuerySchema(t, client, "payload")

        AssertFieldExists(t, schema, "payload_type_id", "Int") // NOT payloadtype_id
        AssertFieldType(t, schema, "callback_alert", "array") // NOT bool
        AssertFieldType(t, schema, "auto_generated", "array") // NOT bool
        AssertFieldExists(t, schema, "deleted", "Boolean")
    })

    // Test 3: Validate buildparameter schema
    t.Run("ValidateBuildParameterSchema", func(t *testing.T) {
        schema := QuerySchema(t, client, "buildparameter")

        AssertFieldExists(t, schema, "payload_type_id", "Int") // NOT payloadtype_id
    })

    // Test 4: Validate C2 parameter requirements
    t.Run("ValidateC2ProfileParameters", func(t *testing.T) {
        profiles, err := client.GetC2Profiles(ctx)
        require.NoError(t, err)

        for _, profile := range profiles {
            params, err := client.GetC2ProfileParameters(ctx, profile.ID)
            require.NoError(t, err)

            // Document required parameters
            for _, param := range params {
                if param.Required {
                    t.Logf("C2 Profile %s requires parameter: %s (default: %s)",
                        profile.Name, param.Name, param.DefaultValue)
                }
            }
        }
    })
}
```

**Expected Result**: This test FAILS on current SDK, passes after fixes
**CI Integration**: Run this test FIRST in Phase 1

---

### PRIORITY 1: CRITICAL CRUD OPERATIONS

These methods are used in 90% of SDK interactions and MUST work:

#### 1.1 Commands & Parameters (commands.go) - **0% Tested**

**Test File**: `commands_test.go`

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetCommands` | `TestE2E_GetCommands_SchemaValidation` | All fields populated, PayloadTypeID > 0, no panics | Used by every payload build |
| `GetCommandParameters` | `TestE2E_GetCommandParameters_Complete` | Parameters exist, types correct, required flags accurate | Required for task building |
| `GetCommandWithParameters` | `TestE2E_GetCommandWithParameters_AllPayloadTypes` | Iterate all payload types, validate params | Core task builder |
| `GetLoadedCommands` | `TestE2E_GetLoadedCommands_WithCallback` | After loading command to callback, verify it appears | Callback management |
| `BuildTaskParams` (raw string) | `TestE2E_BuildTaskParams_RawString` | Shell command, run, execute commands | Most common task type |
| `BuildTaskParams` (JSON) | `TestE2E_BuildTaskParams_JSONParams` | Commands with parameters, validation | Structured tasks |
| `BuildTaskParams` (required missing) | `TestE2E_BuildTaskParams_MissingRequired` | Must error when required param missing | Error handling |

**Expected Failures on Current SDK:**
- ✓ GetCommands panics on `supported_ui_features` (bool vs array)
- ✓ GetCommands panics on `attributes` (string vs jsonb)
- ✓ GetCommands fails on `payloadtype_id` field name
- ✓ GetCommandWithParameters fails on same fields

**Fix Verification:**
- All fields must be nullable or correct type
- PayloadTypeID must always be > 0 for valid commands
- BuildTaskParams must handle all command types

---

#### 1.2 Task Lifecycle (tasks.go) - **5.6% Tested**

**Test File**: `tasks_test.go`

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `IssueTask` | `TestE2E_IssueTask_RawString` | Issue `whoami`, get task ID back | Core functionality |
| `IssueTask` | `TestE2E_IssueTask_WithParams` | Issue parameterized task (ls -la), validate params sent | Structured tasks |
| `IssueTask` | `TestE2E_IssueTask_InvalidCommand` | Issue non-existent command, expect error | Error handling |
| `GetTask` | `TestE2E_GetTask_Complete` | After IssueTask, GetTask returns same task | Task retrieval |
| `GetTasksForCallback` | `TestE2E_GetTasksForCallback_Pagination` | Get tasks, validate all fields populated | Already tested ✓ |
| `GetTaskOutput` | `TestE2E_GetTaskOutput_MultipleResponses` | Issue task with output, retrieve all responses | Output handling |
| `WaitForTaskComplete` | `TestE2E_WaitForTaskComplete_Success` | Issue task, wait for completion, validate status | Synchronous operations |
| `WaitForTaskComplete` | `TestE2E_WaitForTaskComplete_Timeout` | Wait with short timeout, expect timeout error | Timeout handling |
| `WaitForTaskComplete` | `TestE2E_WaitForTaskComplete_Error` | Issue invalid task, wait should return error | Error propagation |
| `UpdateTask` | `TestE2E_UpdateTask_Comment` | Update task comment, verify persisted | Task management |
| `ReissueTask` | `TestE2E_ReissueTask_FailedTask` | Issue task, mark failed, reissue, get new task | Failure recovery |
| `RequestOpsecBypass` | `TestE2E_RequestOpsecBypass_Approval` | Request bypass, verify status change | OPSEC workflow |
| `AddMITREAttackToTask` | `TestE2E_AddMITREAttackToTask_MultipleT Numbers` | Tag task with T1234, T5678, verify persisted | ATT&CK mapping |
| `GetTasksByStatus` | `TestE2E_GetTasksByStatus_Completed` | Get completed tasks only, validate filter | Status filtering |
| `GetTaskArtifacts` | `TestE2E_GetTaskArtifacts_AfterExecution` | Run task creating artifact, verify captured | Artifact tracking |

**Expected Failures on Current SDK:**
- IssueTask may fail if using commands with schema errors
- WaitForTaskComplete may not handle all status codes

---

#### 1.3 Payload Lifecycle (payloads.go) - **25% Tested**

**Test File**: `payloads_test.go`

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `CreatePayload` | `TestE2E_CreatePayload_MinimalConfig` | Poseidon + HTTP + Linux, verify UUID returned | Core functionality |
| `CreatePayload` | `TestE2E_CreatePayload_AllCommands` | Include all 47+ commands, verify build succeeds | Full feature payload |
| `CreatePayload` | `TestE2E_CreatePayload_MissingC2Host` | No callback_host, expect error | Validates required params |
| `CreatePayload` | `TestE2E_CreatePayload_InvalidOS` | Invalid OS value, expect error | Input validation |
| `GetPayloadByUUID` | `TestE2E_GetPayloadByUUID_AfterCreate` | Create, then get, validate all fields | Retrieval validation |
| `GetPayloadByUUID` | `TestE2E_GetPayloadByUUID_BuildPhases` | Check build_phase field during build | Build progress |
| `WaitForPayloadComplete` | `TestE2E_WaitForPayloadComplete_Success` | Wait for build, expect success | Already tested ✓ |
| `WaitForPayloadComplete` | `TestE2E_WaitForPayloadComplete_BuildError` | Trigger build error, wait should return error | Error handling |
| `DownloadPayload` | `TestE2E_DownloadPayload_Complete` | Build payload, download, verify file size > 0 | File download |
| `DeletePayload` | `TestE2E_DeletePayload_MarkDeleted` | Delete payload, verify deleted=true | Cleanup |
| `RebuildPayload` | `TestE2E_RebuildPayload_SameConfig` | Rebuild existing payload, get new UUID | Rebuild workflow |
| `GetPayloadCommands` | `TestE2E_GetPayloadCommands_VerifyIncluded` | Get commands for payload, match build request | Command inclusion |

**Expected Failures on Current SDK:**
- ✓ GetPayloadByUUID panics on `callback_alert` (bool vs array)
- ✓ GetPayloadByUUID panics on `auto_generated` (bool vs array)
- ✓ CreatePayload fails without callback_host parameter

---

#### 1.4 Callback Management (callbacks.go) - **10% Tested**

**Test File**: `callbacks_test.go`

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetAllCallbacks` | `TestE2E_GetAllCallbacks_Complete` | Get all, validate structure | Callback inventory |
| `GetAllActiveCallbacks` | `TestE2E_GetAllActiveCallbacks_OnlyActive` | Already tested ✓, add validation for Active=true | Tested |
| `GetCallbackByID` | `TestE2E_GetCallbackByID_ValidID` | Get callback by display ID, validate fields | Single retrieval |
| `UpdateCallback` | `TestE2E_UpdateCallback_Description` | Update description, verify persisted | Callback management |
| `UpdateCallback` | `TestE2E_UpdateCallback_Sleep` | Update sleep interval, verify sent to agent | Configuration |
| `DeleteCallback` | `TestE2E_DeleteCallback_MarkDeleted` | Delete callback, verify deleted=true | Cleanup |
| `AddCallbackGraphEdge` | `TestE2E_AddCallbackGraphEdge_ParentChild` | Create parent-child relationship, verify graph | Process tree |
| `RemoveCallbackGraphEdge` | `TestE2E_RemoveCallbackGraphEdge_BreakLink` | Remove edge, verify disconnected | Graph management |
| `ExportCallbackConfig` | `TestE2E_ExportCallbackConfig_JSONValid` | Export config, validate JSON structure | Config export |
| `ImportCallbackConfig` | `TestE2E_ImportCallbackConfig_Restore` | Export, delete, import, verify restored | Config import |

---

#### 1.5 File Operations (files.go) - **9.1% Tested**

**Test File**: `files_test.go`

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `UploadFile` | `TestE2E_UploadFile_SmallFile` | Upload 1KB file, verify success | File upload |
| `UploadFile` | `TestE2E_UploadFile_LargeFile` | Upload 10MB file, verify chunks | Large file handling |
| `GetFiles` | `TestE2E_GetFiles_AfterUpload` | Already tested ✓, add pagination test | Tested |
| `GetFileByID` | `TestE2E_GetFileByID_ValidID` | Get specific file, validate metadata | Single retrieval |
| `DownloadFile` | `TestE2E_DownloadFile_CompleteFile` | Upload, download, verify SHA256 match | File integrity |
| `DownloadFile` | `TestE2E_DownloadFile_IncompleteChunks` | Partial upload, download should error | Error handling |
| `DeleteFile` | `TestE2E_DeleteFile_MarkDeleted` | Delete file, verify deleted=true | Cleanup |
| `BulkDownloadFiles` | `TestE2E_BulkDownloadFiles_Multiple` | Upload 5 files, bulk download, verify count | Bulk operations |
| `PreviewFile` | `TestE2E_PreviewFile_TextFile` | Upload text, preview first 1KB | Preview feature |
| `GetDownloadedFiles` | `TestE2E_GetDownloadedFiles_Filter` | Get only downloaded files, validate filter | File filtering |

---

### PRIORITY 2: AUTHENTICATION & SESSION MANAGEMENT

#### 2.1 Authentication Flow (auth.go) - **12.5% Tested**

**Test File**: `auth_test.go`

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `Login` | `TestE2E_Login_ValidCredentials` | Login, get tokens, validate format | Core auth |
| `Login` | `TestE2E_Login_InvalidCredentials` | Bad password, expect auth error | Error handling |
| `GetMe` | `TestE2E_GetMe_AfterLogin` | Already tested ✓, validate all fields | Tested |
| `CreateAPIToken` | `TestE2E_CreateAPIToken_ValidToken` | Create token, use for auth, validate works | Token creation |
| `CreateAPIToken` | `TestE2E_CreateAPIToken_WithExpiry` | Create with 1 hour expiry, validate expires_at | Expiry handling |
| `RefreshAccessToken` | `TestE2E_RefreshAccessToken_BeforeExpiry` | Wait for near-expiry, refresh, get new token | Token refresh |
| `Logout` | `TestE2E_Logout_InvalidatesSession` | Logout, subsequent call should fail | Session cleanup |
| `ClearRefreshToken` | `TestE2E_ClearRefreshToken_CantRefresh` | Clear token, refresh should fail | Token cleanup |

---

### PRIORITY 3: ADVANCED FEATURES

#### 3.1 Build Parameters (buildparameters.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetBuildParameters` | `TestE2E_GetBuildParameters_AllTypes` | Get all, validate types (String, Boolean, etc.) | Parameter discovery |
| `GetBuildParametersByPayloadType` | `TestE2E_GetBuildParametersByPayloadType_Poseidon` | Get for Poseidon, validate required params | Type-specific params |
| `GetBuildParameterInstances` | `TestE2E_GetBuildParameterInstances_AfterBuild` | Build payload, get instances, verify values | Instance tracking |
| `GetBuildParameterInstancesByPayload` | `TestE2E_GetBuildParameterInstancesByPayload_SpecificPayload` | Get instances for one payload, validate | Payload config |

**Expected Failures:**
- ✓ All methods fail on `payloadtype_id` field name

---

#### 3.2 MITRE ATT&CK (attack.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetAttackTechniques` | `TestE2E_GetAttackTechniques_AllT Numbers` | Get all techniques, validate T-numbers | Technique inventory |
| `GetAttackTechniqueByTNum` | `TestE2E_GetAttackTechniqueByTNum_T1059` | Get T1059, validate name/description | Single lookup |
| `GetAttackByTask` | `TestE2E_GetAttackByTask_AfterTagging` | Tag task with technique, get techniques, verify | Task mapping |
| `GetAttackByCommand` | `TestE2E_GetAttackByCommand_InherentMapping` | Commands have default techniques, verify | Command mapping |
| `GetAttacksByOperation` | `TestE2E_GetAttacksByOperation_Coverage` | Get all techniques used in operation | Operation coverage |

---

#### 3.3 Eventing System (eventing.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `EventingTriggerManual` | `TestE2E_EventingTriggerManual_SimpleWorkflow` | Trigger webhook event, verify executes | Manual trigger |
| `EventingTriggerKeyword` | `TestE2E_EventingTriggerKeyword_OnTaskOutput` | Set keyword trigger, issue task, verify fires | Keyword automation |
| `EventingTriggerCancel` | `TestE2E_EventingTriggerCancel_Running` | Cancel running event, verify stopped | Cancellation |
| `EventingTriggerRetry` | `TestE2E_EventingTriggerRetry_Failed` | Retry failed event, verify re-executes | Retry logic |
| `EventingExportWorkflow` | `TestE2E_EventingExportWorkflow_ValidJSON` | Export workflow, validate JSON structure | Export feature |
| `EventingImportContainerWorkflow` | `TestE2E_EventingImportContainerWorkflow_FromFile` | Import workflow, verify appears | Import feature |
| `UpdateEventGroupApproval` | `TestE2E_UpdateEventGroupApproval_Approve` | Approve pending event, verify proceeds | Approval workflow |
| `SendExternalWebhook` | `TestE2E_SendExternalWebhook_HTTPPost` | Send webhook, verify received (mock server) | External integration |

---

#### 3.4 Subscriptions (subscriptions.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `Subscribe` | `TestE2E_Subscribe_CallbackUpdates` | Subscribe to callbacks, create callback, verify notification | Real-time updates |
| `Unsubscribe` | `TestE2E_Unsubscribe_StopReceiving` | Subscribe, unsubscribe, verify no more notifications | Cleanup |

---

### PRIORITY 4: OPERATIONS & OPERATORS

#### 4.1 Operations Management (operations.go) - **18.2% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `CreateOperation` | `TestE2E_CreateOperation_NewOp` | Create operation, verify created | Operation creation |
| `UpdateOperation` | `TestE2E_UpdateOperation_Name` | Update name, verify persisted | Operation management |
| `GetOperationByID` | `TestE2E_GetOperationByID_ValidID` | Get specific operation, validate fields | Single retrieval |
| `GetOperatorsByOperation` | `TestE2E_GetOperatorsByOperation_Members` | Get operators in operation, validate list | Member management |
| `UpdateOperatorOperation` | `TestE2E_UpdateOperatorOperation_ChangeRole` | Change operator role, verify updated | Permission management |
| `GetOperationEventLog` | `TestE2E_GetOperationEventLog_AuditTrail` | Get event log, validate entries | Audit trail |
| `CreateOperationEventLog` | `TestE2E_CreateOperationEventLog_CustomEntry` | Create log entry, verify appears | Manual logging |

---

#### 4.2 Operator Management (operators.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetOperators` | `TestE2E_GetOperators_AllUsers` | Get all operators, validate structure | User inventory |
| `CreateOperator` | `TestE2E_CreateOperator_NewUser` | Create operator, verify created | User creation |
| `UpdateOperatorStatus` | `TestE2E_UpdateOperatorStatus_Active` | Change status, verify updated | User management |
| `UpdatePasswordAndEmail` | `TestE2E_UpdatePasswordAndEmail_ChangeCredentials` | Update creds, verify can login with new | Credential management |
| `GetOperatorPreferences` | `TestE2E_GetOperatorPreferences_Theme` | Get prefs, validate structure | User preferences |
| `UpdateOperatorPreferences` | `TestE2E_UpdateOperatorPreferences_Change` | Update prefs, verify persisted | Preference management |
| `CreateInviteLink` | `TestE2E_CreateInviteLink_ValidLink` | Create link, validate URL | User invitation |

---

### PRIORITY 5: ARTIFACT TRACKING & REPORTING

#### 5.1 Artifacts (artifacts.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `CreateArtifact` | `TestE2E_CreateArtifact_NewEntry` | Create artifact, verify saved | Artifact tracking |
| `GetArtifacts` | `TestE2E_GetArtifacts_AllTypes` | Get all, validate types (File, Process, etc.) | Artifact inventory |
| `GetArtifactsByHost` | `TestE2E_GetArtifactsByHost_Filter` | Get for specific host, validate filter | Host-specific |
| `UpdateArtifact` | `TestE2E_UpdateArtifact_AddNote` | Update notes, verify persisted | Artifact annotation |
| `DeleteArtifact` | `TestE2E_DeleteArtifact_Remove` | Delete artifact, verify removed | Cleanup |

---

#### 5.2 Credentials (credentials.go) - **16.7% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetCredentials` | `TestE2E_GetCredentials_AllTypes` | Already tested ✓, add type validation | Tested |
| `CreateCredential` | `TestE2E_CreateCredential_Password` | Create plaintext password, verify saved | Cred storage |
| `CreateCredential` | `TestE2E_CreateCredential_Hash` | Create NTLM hash, verify saved | Hash storage |
| `UpdateCredential` | `TestE2E_UpdateCredential_AddComment` | Update comment, verify persisted | Cred annotation |
| `DeleteCredential` | `TestE2E_DeleteCredential_Remove` | Delete credential, verify removed | Cleanup |

---

### PRIORITY 6: OBSERVABILITY & MONITORING

#### 6.1 Screenshots (screenshots.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetScreenshots` | `TestE2E_GetScreenshots_ForCallback` | Get screenshots, validate structure | Screenshot retrieval |
| `DownloadScreenshot` | `TestE2E_DownloadScreenshot_ValidImage` | Download, verify image format | Image download |
| `GetScreenshotThumbnail` | `TestE2E_GetScreenshotThumbnail_Smaller` | Get thumbnail, verify smaller than full | Thumbnail feature |
| `DeleteScreenshot` | `TestE2E_DeleteScreenshot_Remove` | Delete, verify removed | Cleanup |

---

#### 6.2 Processes (processes.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetProcesses` | `TestE2E_GetProcesses_AllProcesses` | Get processes, validate structure | Process tracking |
| `GetProcessesByCallback` | `TestE2E_GetProcessesByCallback_Specific` | Get for callback, validate filter | Callback processes |
| `GetProcessTree` | `TestE2E_GetProcessTree_Hierarchy` | Get tree, validate parent-child relationships | Process hierarchy |

---

#### 6.3 Keylogs (keylogs.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetKeylogs` | `TestE2E_GetKeylogs_AllLogs` | Get keylogs, validate structure | Keylog retrieval |
| `GetKeylogsByCallback` | `TestE2E_GetKeylogsByCallback_Specific` | Get for callback, validate filter | Callback keylogs |

---

#### 6.4 Responses (responses.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetResponsesByTask` | `TestE2E_GetResponsesByTask_AllResponses` | Get task responses, validate order | Response retrieval |
| `GetResponsesByCallback` | `TestE2E_GetResponsesByCallback_Paginated` | Get callback responses, test pagination | Bulk retrieval |
| `SearchResponses` | `TestE2E_SearchResponses_Filter` | Search with filter, validate results | Search feature |
| `GetLatestResponses` | `TestE2E_GetLatestResponses_Recent` | Get recent N responses, validate order | Latest data |

---

### PRIORITY 7: NETWORKING & INFRASTRUCTURE

#### 7.1 C2 Profiles (c2profiles.go) - **11.1% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetC2Profiles` | `TestE2E_GetC2Profiles_AllProfiles` | Already tested ✓, add parameter validation | Tested |
| `StartStopProfile` | `TestE2E_StartStopProfile_Toggle` | Stop HTTP, verify stopped, start, verify running | Profile control |
| `CreateC2Instance` | `TestE2E_CreateC2Instance_NewInstance` | Create new C2 instance, verify config | Instance creation |
| `GetProfileOutput` | `TestE2E_GetProfileOutput_Logs` | Get C2 logs, validate format | Log retrieval |
| `C2HostFile` | `TestE2E_C2HostFile_ServeFile` | Host file, verify accessible via HTTP | File hosting |
| `C2GetIOC` | `TestE2E_C2GetIOC_Indicators` | Get IOCs, validate format | IOC extraction |

---

#### 7.2 RPFWD (rpfwd.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `CreateRPFWD` | `TestE2E_CreateRPFWD_NewForward` | Create reverse port forward, verify created | RPFWD creation |
| `GetRPFWDs` | `TestE2E_GetRPFWDs_AllForwards` | Get all forwards, validate structure | RPFWD inventory |
| `DeleteRPFWD` | `TestE2E_DeleteRPFWD_Remove` | Delete forward, verify removed | Cleanup |

---

#### 7.3 Proxy (proxy.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `ToggleProxy` | `TestE2E_ToggleProxy_Enable` | Enable proxy for task, verify routing | Proxy control |
| `TestProxy` | `TestE2E_TestProxy_Connectivity` | Test proxy connection, verify success | Proxy testing |

---

### PRIORITY 8: ADVANCED FEATURES

#### 8.1 Dynamic Query (dynamicquery.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `DynamicQueryFunction` | `TestE2E_DynamicQueryFunction_Execute` | Execute dynamic query, validate result | Dynamic queries |
| `DynamicBuildParameter` | `TestE2E_DynamicBuildParameter_Options` | Get dynamic options, validate list | Dynamic config |

---

#### 8.2 Hosts & File Browser (hosts.go, filebrowser.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetHosts` | `TestE2E_GetHosts_AllHosts` | Get hosts, validate structure | Host inventory |
| `GetHostNetworkMap` | `TestE2E_GetHostNetworkMap_Topology` | Get network map, validate connections | Network topology |
| `GetFileBrowserObjects` | `TestE2E_GetFileBrowserObjects_AllFiles` | Get file browser objects, validate | File browser |

---

#### 8.3 Tags (tags.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `CreateTagType` | `TestE2E_CreateTagType_NewType` | Create tag type, verify created | Tag type creation |
| `CreateTag` | `TestE2E_CreateTag_ApplyToCallback` | Tag callback, verify applied | Tagging |
| `GetTags` | `TestE2E_GetTags_ForSource` | Get tags for object, validate | Tag retrieval |
| `DeleteTag` | `TestE2E_DeleteTag_Remove` | Delete tag, verify removed | Cleanup |

---

#### 8.4 Alerts (alerts.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetAlerts` | `TestE2E_GetAlerts_AllTypes` | Get alerts, validate types | Alert inventory |
| `CreateCustomAlert` | `TestE2E_CreateCustomAlert_NewAlert` | Create alert, verify appears | Alert creation |
| `ResolveAlert` | `TestE2E_ResolveAlert_MarkResolved` | Resolve alert, verify status | Alert management |

---

#### 8.5 Reporting (reporting.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GenerateReport` | `TestE2E_GenerateReport_HTML` | Generate HTML report, validate format | Report generation |
| `GetRedirectRules` | `TestE2E_GetRedirectRules_ForPayload` | Get redirect rules, validate | Redirect rules |

---

#### 8.6 Tokens (tokens.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetTokens` | `TestE2E_GetTokens_AllTokens` | Get tokens, validate structure | Token inventory |
| `GetCallbackTokens` | `TestE2E_GetCallbackTokens_ForCallback` | Get tokens for callback, validate | Callback tokens |
| `DeleteAPIToken` | `TestE2E_DeleteAPIToken_Revoke` | Delete token, verify revoked | Token cleanup |

---

#### 8.7 Containers (containers.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `ContainerListFiles` | `TestE2E_ContainerListFiles_PayloadType` | List files in container, validate | Container inspection |
| `ContainerDownloadFile` | `TestE2E_ContainerDownloadFile_Config` | Download file from container, verify | File retrieval |
| `ContainerWriteFile` | `TestE2E_ContainerWriteFile_Upload` | Write file to container, verify | File upload |

---

#### 8.8 Browser Scripts (browserscripts.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetBrowserScripts` | `TestE2E_GetBrowserScripts_AllScripts` | Get scripts, validate structure | Script inventory |
| `CustomBrowserExport` | `TestE2E_CustomBrowserExport_Config` | Export browser config, validate JSON | Config export |

---

#### 8.9 Staging (staging.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `GetStagingInfo` | `TestE2E_GetStagingInfo_Details` | Get staging info, validate | Staging details |

---

#### 8.10 Blocklist (blocklist.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `DeleteBlockList` | `TestE2E_DeleteBlockList_Remove` | Delete blocklist, verify removed | Blocklist management |

---

#### 8.11 Utility (utility.go) - **0% Tested**

| Method | Test Name | What To Validate | Why Critical |
|--------|-----------|------------------|--------------|
| `CreateRandom` | `TestE2E_CreateRandom_Generate` | Generate random data, validate format | Random generation |
| `ConfigCheck` | `TestE2E_ConfigCheck_Validate` | Check server config, validate response | Config validation |

---

## PART 5: TEST IMPLEMENTATION STRATEGY

### Phase 1: Schema Validation Foundation (Week 1)

**Goal**: Prevent ALL schema bugs

1. Implement `TestE2E_SchemaValidation` covering:
   - All GraphQL types (command, payload, buildparameter, etc.)
   - Field existence validation
   - Type validation (bool vs array, string vs jsonb)
   - Required parameter validation

2. Add GraphQL introspection helpers:
   ```go
   func QuerySchema(t *testing.T, client *Client, typeName string) SchemaType
   func AssertFieldExists(t *testing.T, schema SchemaType, fieldName string, expectedType string)
   func AssertFieldNotExists(t *testing.T, schema SchemaType, fieldName string)
   func AssertFieldType(t *testing.T, schema SchemaType, fieldName string, expectedType string)
   ```

3. Run this test FIRST in CI Phase 1
4. **Expected**: Test FAILS on current SDK, PASSES after fixes

### Phase 2: Critical CRUD Operations (Week 2-3)

**Goal**: 100% coverage of Priority 1 methods

1. Commands & Parameters (7 tests)
2. Task Lifecycle (15 tests)
3. Payload Lifecycle (12 tests)
4. Callback Management (10 tests)
5. File Operations (10 tests)

**Total**: 54 new tests

### Phase 3: Authentication & Operations (Week 4)

**Goal**: Complete session management coverage

1. Authentication Flow (8 tests)
2. Operations Management (7 tests)
3. Build Parameters (4 tests)
4. MITRE ATT&CK (5 tests)

**Total**: 24 new tests

### Phase 4: Advanced Features (Week 5-6)

**Goal**: Complete feature parity

1. Eventing System (8 tests)
2. Subscriptions (2 tests)
3. Operator Management (7 tests)
4. C2 Profiles (6 tests)
5. Dynamic Query (2 tests)

**Total**: 25 new tests

### Phase 5: Observability & Reporting (Week 7)

**Goal**: Complete monitoring coverage

1. Artifacts (5 tests)
2. Credentials (4 tests)
3. Screenshots (4 tests)
4. Processes (3 tests)
5. Keylogs (2 tests)
6. Responses (4 tests)
7. Reporting (2 tests)

**Total**: 24 new tests

### Phase 6: Networking & Infrastructure (Week 8)

**Goal**: Complete infrastructure coverage

1. RPFWD (3 tests)
2. Proxy (2 tests)
3. Hosts & File Browser (3 tests)
4. Tokens (3 tests)
5. Containers (3 tests)

**Total**: 14 new tests

### Phase 7: Remaining Features (Week 9)

**Goal**: 100% method coverage

1. Tags (4 tests)
2. Alerts (3 tests)
3. Browser Scripts (2 tests)
4. Staging (1 test)
5. Blocklist (1 test)
6. Utility (2 tests)

**Total**: 13 new tests

---

## PART 6: CI/CD IMPROVEMENTS

### Current CI Issues

1. **Uses `latest` Mythic** - schema can change
2. **No schema validation** - tests pass despite field mismatches
3. **No field-level assertions** - only checks `err != nil`
4. **Integration tests may use stale data** - not creating fresh objects

### Recommended CI Changes

#### 1. Pin Mythic Version

```yaml
- name: Clone Mythic Framework
  run: |
    # Pin to specific version or tag
    git clone --branch v3.2.0 --depth 1 https://github.com/its-a-feature/Mythic.git /tmp/mythic
    # OR use commit SHA
    # git clone https://github.com/its-a-feature/Mythic.git /tmp/mythic
    # cd /tmp/mythic && git checkout abc123def456
```

**Rationale**: Prevents schema drift from breaking tests

#### 2. Add Schema Validation Phase

```yaml
- name: Run Schema Validation Tests
  run: |
    go test -v -tags=integration ./tests/integration/... \
      -run "TestE2E_SchemaValidation" \
      -timeout 5m

  # If schema validation fails, STOP immediately
  # Don't run other tests
```

**Rationale**: Catch schema bugs FIRST before other tests

#### 3. Add Mythic Version Documentation

```yaml
- name: Document Mythic Version
  run: |
    cd /tmp/mythic
    echo "Mythic Version: $(cat VERSION)" >> $GITHUB_STEP_SUMMARY
    echo "Git Commit: $(git rev-parse HEAD)" >> $GITHUB_STEP_SUMMARY
    echo "Schema Version: $(curl -k -s https://127.0.0.1:7443/api/version)" >> $GITHUB_STEP_SUMMARY
```

**Rationale**: Know which Mythic version tests ran against

#### 4. Enhanced Test Assertions

Add to all tests:

```go
// Don't just check error
commands, err := client.GetCommands(ctx)
require.NoError(t, err)

// ALSO validate response structure
require.NotEmpty(t, commands, "GetCommands returned empty list")

for _, cmd := range commands {
    // Validate required fields populated
    assert.NotZero(t, cmd.ID, "Command ID is zero")
    assert.NotEmpty(t, cmd.Cmd, "Command name is empty")
    assert.NotZero(t, cmd.PayloadTypeID, "PayloadTypeID is zero")
    assert.NotEmpty(t, cmd.Description, "Description is empty")

    // Validate field types
    assert.IsType(t, "", cmd.Cmd, "Cmd should be string")
    assert.IsType(t, 0, cmd.PayloadTypeID, "PayloadTypeID should be int")
    assert.IsType(t, false, cmd.ScriptOnly, "ScriptOnly should be bool")
}
```

#### 5. Test Data Cleanup

Ensure tests clean up after themselves:

```go
func TestE2E_Example(t *testing.T) {
    client := AuthenticateTestClient(t)

    // Create test data
    payload, err := client.CreatePayload(ctx, req)
    require.NoError(t, err)

    // Cleanup on test end
    t.Cleanup(func() {
        client.DeletePayload(context.Background(), payload.UUID)
    })

    // Run test...
}
```

#### 6. Parallel Test Execution

```yaml
strategy:
  matrix:
    test-group:
      - "Schema Validation"
      - "Commands & Tasks"
      - "Payloads & Files"
      - "Callbacks & Operations"
      - "Advanced Features"
```

**Rationale**: Faster CI runs

---

## PART 7: TEST COVERAGE METRICS

### Coverage Goals

| Metric | Current | Target | Timeline |
|--------|---------|--------|----------|
| Method Coverage | 4.7% (11/232) | 100% (232/232) | 9 weeks |
| Line Coverage | Unknown | 80%+ | 9 weeks |
| Branch Coverage | Unknown | 70%+ | 9 weeks |
| Integration Tests | 39 files | 154 new tests | 9 weeks |
| Schema Validation | 0% | 100% | Week 1 |

### Weekly Milestones

| Week | Focus Area | Tests Added | Cumulative Coverage |
|------|-----------|-------------|---------------------|
| 1 | Schema Validation | 1 (critical) | 5% |
| 2-3 | Critical CRUD | 54 | 28% |
| 4 | Auth & Operations | 24 | 38% |
| 5-6 | Advanced Features | 25 | 49% |
| 7 | Observability | 24 | 59% |
| 8 | Infrastructure | 14 | 65% |
| 9 | Remaining | 13 | **100%** |

---

## PART 8: ACCEPTANCE CRITERIA

### Definition of "Complete Coverage"

A method has complete test coverage when:

1. **Happy Path Test**: Validates successful execution with valid inputs
2. **Error Handling Test**: Validates error cases (invalid input, not found, etc.)
3. **Schema Validation**: All response fields are asserted (type, non-zero, format)
4. **Integration Test**: Uses real Mythic instance, not mocks
5. **Cleanup**: Test cleans up created resources
6. **Documentation**: Test name clearly describes what's validated

### Test Quality Checklist

Each test must:
- ✅ Use `require.NoError()` for critical assertions
- ✅ Use `assert.*()` for field validations
- ✅ Validate response structure, not just error=nil
- ✅ Create fresh test data (don't rely on existing data)
- ✅ Clean up created resources in `t.Cleanup()`
- ✅ Have clear test name: `TestE2E_<Method>_<Scenario>`
- ✅ Run in < 30 seconds (unless testing timeouts)
- ✅ Be idempotent (can run multiple times)

---

## PART 9: IMMEDIATE ACTION ITEMS

### Week 1 Sprint Plan

**Day 1-2: Schema Validation Foundation**
1. Implement GraphQL introspection helpers
2. Write `TestE2E_SchemaValidation` test
3. Run test, document all schema mismatches
4. Fix all schema mismatches in SDK
5. Verify test passes

**Day 3-5: Commands & Tasks Critical Path**
1. Implement 7 command tests
2. Implement 5 task tests (IssueTask, GetTask, WaitForTaskComplete, etc.)
3. Run tests, fix any bugs found
4. Document test results

**Expected Outcome**:
- Schema validation catches future bugs automatically
- Commands & Tasks fully tested
- 3 new bugs discovered and fixed (estimate)

---

## CONCLUSION

**Current State**: SDK has 4.7% integration test coverage (11/232 methods)

**Root Cause**: Tests exist but don't validate GraphQL schema, field types, or response structures

**Impact**: Production bugs (schema mismatches, type errors, missing parameters)

**Solution**: 154 new integration tests covering all 232 methods with field-level validation

**Timeline**: 9 weeks to 100% coverage

**Priority 0**: Schema validation test (Week 1) - prevents ALL schema bugs

**Next Steps**:
1. Review and approve this plan
2. Start Week 1: Schema validation
3. Fix all discovered schema bugs
4. Continue with Priority 1 critical CRUD tests

**Success Criteria**: No SDK bugs reach production ever again
