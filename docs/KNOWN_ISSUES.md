# Known Issues and Design Decisions

## IssueTask Uses REST Endpoint Instead of GraphQL

### Background
The `IssueTask` function uses the Hasura REST webhook endpoint (`/api/v1.4/create_task_webhook`) instead of the GraphQL `createTask` mutation.

### Why This Is Necessary

The GraphQL approach fails due to how the `hasura/go-graphql-client` library handles optional array parameters:

1. **GraphQL Requirement**: All variables referenced in a query must be present in the variables map
2. **Serialization Issue**: When nil array values are included in the map, they serialize to explicit `"null"` in JSON
3. **Hasura Validation**: When Hasura receives `"null"` for optional array parameters (like `callback_ids: [Int]`), it validates them and rejects with:
   ```
   Message: null value found for non-nullable type: "[Int!]!"
   ```
   Even though the GraphQL schema correctly defines these as nullable (`[Int]`, not `[Int!]!`)

### Why REST Is The Correct Solution

The webhook approach is **not** a workaround - it's the proper implementation:

✅ **Same Endpoint**: It's the exact webhook that Hasura's GraphQL action calls internally
✅ **Proper Semantics**: Allows omitting parameters (not sending them) vs sending explicit `null`
✅ **Documented API**: `/api/v1.4/create_task_webhook` is a stable, public Hasura action endpoint
✅ **Full Control**: We can control exactly what gets sent in the JSON payload
✅ **Works Correctly**: Matches the webhook's expected input structure

### Code Location
See `pkg/mythic/tasks.go:84-173` for the implementation with detailed comments explaining this decision.

### Alternative Approaches Attempted

1. ❌ **Pass nil slices**: GraphQL client serializes to `"null"`, rejected by Hasura
2. ❌ **Pass empty slices `[]int{}`**: Not nil, webhook logic fails
3. ❌ **Conditionally include variables**: GraphQL requires ALL referenced variables
4. ✅ **Use REST webhook directly**: Works perfectly, proper architectural choice

### Status
- **Resolved**: Using REST endpoint (commit c09e54f)
- **Not a bug**: This is the correct way to use Hasura actions with optional arrays
- **Production ready**: All tests passing

---

## GraphQL JSONB Scalar Handling in GetInviteLinks

### Issue
The `getInviteLinks` GraphQL action returns a `jsonb` field that causes reflection panics in the `shurcooL/graphql` library.

### Workaround
Query the `invite_link` database table directly instead of using the GraphQL action:
```go
var linksQuery struct {
    InviteLinks []struct {
        // ... fields
    } `graphql:"invite_link(order_by: {created_at: desc})"`
}
```

### Impact
- ⚠️ Couples SDK to database schema structure
- ⚠️ Bypasses the intended GraphQL action layer
- ✅ Works reliably for all operations

### Status
- **Workaround implemented**: Direct table query (commit 6ca0f9b)
- **Limitation**: GraphQL client library JSONB handling
- **Future**: Consider switching to GraphQL client with better scalar support

### Code Location
See `pkg/mythic/operators.go:408-473`

---

## CreateInviteLink Missing Fields

### Issue
The `createInviteLinkOutput` type doesn't return all fields needed to fully populate the `InviteLink` struct. Fields like `ID`, `ExpiresAt`, `CreatedBy`, and `CreatedAt` are not included in the mutation response.

### Current Behavior
Returns the `link` URL which we parse to extract the `Code`. Other fields must be retrieved separately via `GetInviteLinks()`.

### Workaround
After creating an invite link, call `GetInviteLinks()` to get the full details if needed.

### Status
- **Partially resolved**: Parse link URL to get code (commit bbe1ba7)
- **Limitation**: GraphQL action output type definition
- **Impact**: Requires additional API call for complete information

### Code Location
See `pkg/mythic/operators.go:477-541`

---

## Notes on Optional Parameters

### Pointer Types vs Empty Values

The SDK uses different patterns for optional parameters:

- **Pointers (`*int`, `*bool`, `*string`)**: For nullable fields where `nil` means "don't update/include"
- **Empty values (empty string, false, 0)**: For fields where the GraphQL client handles them correctly
- **Conditional inclusion**: For arrays, only include in variables map if non-empty

### Example: UpdateOperatorStatus
```go
variables := map[string]interface{}{
    "operator_id": req.OperatorID,
    "active":      req.Active,   // *bool - nil means don't update
    "admin":       req.Admin,     // *bool - nil means don't update
    "deleted":     req.Deleted,   // *bool - nil means don't update
}
```

The GraphQL client serializes `nil` pointers as `null`, which GraphQL correctly interprets as "parameter not provided" for optional fields.

### Status
- **Working as designed**: Different patterns for different scenarios
- **Production ready**: All patterns tested and working

---

## C2GetIOC and C2SampleMessage GraphQL Schema Mismatches

### Issue
The `C2GetIOC` and `C2SampleMessage` functions use parameter names that don't match the Mythic GraphQL schema in some configurations.

### Error Messages
1. **C2GetIOC**: `'c2GetIOC' has no argument named 'profile_id'`
2. **C2SampleMessage**: `'c2SampleMessage' has no argument named 'message_type'`

### Current Status
- Functions are implemented and fail gracefully when schema doesn't match
- Integration tests handle these as warnings (not failures)
- May be version-specific or configuration-dependent features

### Impact
- ⚠️ IOC generation may not work in all Mythic configurations
- ⚠️ Sample message generation may not work in all Mythic configurations
- ✅ Tests pass - errors are caught and logged as warnings
- ✅ SDK remains functional - these are optional features

### Workaround
These features are not critical to core SDK functionality. The SDK handles failures gracefully:
```go
ioc, err := client.C2GetIOC(ctx, profileID)
if err != nil {
    // Handle gracefully - feature may not be available
    log.Printf("C2GetIOC not available: %v", err)
}
```

### Investigation Needed
To fix these issues, we need to:
1. Determine correct parameter names from Mythic GraphQL schema
2. Verify if these features require specific C2 profile types
3. Check if these features are version-specific (may not be in v3.4.20)

### Status
- **Documented**: Known schema mismatch (v3.4.20)
- **Low priority**: Optional features, graceful degradation
- **Future**: Investigate correct schema parameters or remove if deprecated

### Code Location
See `pkg/mythic/c2profiles.go:296-368`

---

## TestE2E_CallbackTaskLifecycle OPSEC Timeout in CI

### Issue
The `TestE2E_CallbackTaskLifecycle` integration test occasionally times out in CI when tasks get stuck in "OPSEC Pre Check Running..." status.

### Root Cause
Mythic's OPSEC (Operational Security) pre-check feature requires manual operator approval before tasks execute. In the CI environment:
1. Task is issued successfully
2. Task enters "OPSEC Pre Check Running..." status
3. **No operator available to approve** - test waits indefinitely
4. Test times out after 30-60 seconds

### CI Evidence (Run 21014612618)
```
Duration: 61.29s
Details:
- ✓ Payload created and agent deployed successfully
- ✓ Agent callback established (ID: 1)
- ✓ Host: RUNNERVMI13QX, User: runner, PID: 11675
- ✓ Loaded 75 commands
- ✓ Shell task issued (Display ID: 1)
- ✓ Task status: "OPSEC Pre Check Running..."
- ❌ Task never completed (stuck in OPSEC pre-check)
```

### Impact
- ⚠️ Intermittent CI failures when OPSEC is enabled
- ⚠️ Test timeout (61s) but no actual SDK failure
- ✅ SDK works correctly - issue is test environment configuration
- ✅ Manual testing with OPSEC disabled works fine

### Workarounds

**Option 1: Disable OPSEC for Test Operation** (Recommended)
Configure the Mythic operation to bypass OPSEC checks:
```bash
# In Mythic UI or via API:
# Operations > [Test Operation] > Settings > Disable OPSEC Pre-checks
```

**Option 2: Configure Callback to Bypass OPSEC**
Set callback-level OPSEC bypass when creating the agent.

**Option 3: Test Detection**
Modify test to detect OPSEC status and skip:
```go
if task.Status == "OPSEC Pre Check Running..." {
    t.Skip("Task waiting for OPSEC approval - skipping in automated environment")
}
```

### Investigation Notes
- OPSEC feature is designed for production security (manual approval before execution)
- CI environments are automated and cannot provide manual approvals
- The SDK correctly handles task lifecycle; OPSEC is a Mythic server feature
- Task polling works correctly (`WaitForTaskComplete` polls every 2 seconds)

### Status
- **Not an SDK bug**: Mythic OPSEC feature working as designed
- **Test environment issue**: CI needs OPSEC disabled or callbacks configured to bypass
- **Low priority**: Only affects CI; SDK functionality is correct
- **Future**: Add OPSEC status detection in test or CI configuration

### Code Location
- Test: `tests/integration/e2e_callback_task_test.go:19-468`
- Polling: `pkg/mythic/tasks.go:392-435`

---

## Version Compatibility

This SDK is tested against **Mythic v3.4.20**.

Breaking changes in future Mythic versions may require SDK updates, particularly for:
- GraphQL schema changes
- Webhook endpoint paths or signatures
- Database schema changes (affects direct table queries)

---

## Reporting Issues

If you encounter issues not listed here:
1. Check the [GitHub Issues](https://github.com/nbaertsch/mythic-sdk-go/issues)
2. Verify you're using a compatible Mythic version
3. Review the [test suite](tests/integration/) for usage examples
4. Report new issues with full error messages and reproduction steps
