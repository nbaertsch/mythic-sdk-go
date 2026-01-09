package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// GetStagingInfo retrieves all payload staging information for the current operation.
// Staging is used when payloads are delivered in multiple parts or hosted on C2 infrastructure.
//
// Staging information includes:
//   - Staging UUIDs for payload identification
//   - Encryption/decryption keys for staged payloads
//   - C2 profile association for delivery
//   - Expiration times for temporary staging
//   - Active status for enabled/disabled staging
//
// This is useful for:
//   - Monitoring active staged payloads
//   - Managing staged payload delivery
//   - Rotating encryption keys for staged content
//   - Tracking staging endpoint usage
//   - Cleaning up expired staging entries
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//
// Returns:
//   - []*types.StagingInfo: List of staging information entries
//   - error: Error if the operation fails
//
// Example:
//
//	// Get all staging info for current operation
//	stagingList, err := client.GetStagingInfo(ctx)
//	if err != nil {
//	    return err
//	}
//
//	fmt.Printf("Found %d staging entries\n", len(stagingList))
//	for _, staging := range stagingList {
//	    if staging.IsActive() {
//	        fmt.Printf("  %s\n", staging.String())
//	        if staging.HasEncryption() {
//	            fmt.Println("    Uses encryption")
//	        }
//	        if staging.HasExpiration() {
//	            fmt.Printf("    Expires: %s\n", staging.ExpirationTime)
//	        }
//	    }
//	}
func (c *Client) GetStagingInfo(ctx context.Context) ([]*types.StagingInfo, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	var query struct {
		StagingInfo []struct {
			ID             int    `graphql:"id"`
			PayloadID      int    `graphql:"payload_id"`
			StagingUUID    string `graphql:"staging_uuid"`
			EncryptionKey  string `graphql:"encryption_key"`
			DecryptionKey  string `graphql:"decryption_key"`
			C2ProfileID    int    `graphql:"c2profile_id"`
			OperationID    int    `graphql:"operation_id"`
			CreationTime   string `graphql:"creation_time"`
			ExpirationTime string `graphql:"expiration_time"`
			Active         bool   `graphql:"active"`
			Deleted        bool   `graphql:"deleted"`
			OperatorID     int    `graphql:"operator_id"`
		} `graphql:"staginginfo(where: {operation_id: {_eq: $operation_id}, deleted: {_eq: false}}, order_by: {creation_time: desc})"`
	}

	operationID := c.GetCurrentOperation()
	if operationID == nil {
		return nil, WrapError("GetStagingInfo", ErrNotAuthenticated, "no current operation set")
	}

	variables := map[string]interface{}{
		"operation_id": *operationID,
	}

	if err := c.executeQuery(ctx, &query, variables); err != nil {
		return nil, WrapError("GetStagingInfo", err, "failed to get staging info")
	}

	// Convert to types.StagingInfo
	result := make([]*types.StagingInfo, len(query.StagingInfo))
	for i, s := range query.StagingInfo {
		result[i] = &types.StagingInfo{
			ID:             s.ID,
			PayloadID:      s.PayloadID,
			StagingUUID:    s.StagingUUID,
			EncryptionKey:  s.EncryptionKey,
			DecryptionKey:  s.DecryptionKey,
			C2ProfileID:    s.C2ProfileID,
			OperationID:    s.OperationID,
			CreationTime:   s.CreationTime,
			ExpirationTime: s.ExpirationTime,
			Active:         s.Active,
			Deleted:        s.Deleted,
			OperatorID:     s.OperatorID,
		}
	}

	return result, nil
}
