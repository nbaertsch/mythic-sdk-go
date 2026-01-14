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
