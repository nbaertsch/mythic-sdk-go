package mythic

import (
	"context"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic/types"
)

// DeleteBlockList deletes a block list and all its entries.
// Block lists are used to prevent specific IP addresses, domains, or user agents
// from accessing C2 infrastructure, helping avoid detection by security tools.
//
// Deleting a block list removes the list and all associated entries.
// This is typically done when a block list is no longer needed or was created
// in error.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - blockListID: The ID of the block list to delete
//
// Returns:
//   - *types.DeleteBlockListResponse: Result of the deletion operation
//   - error: Error if the operation fails
//
// Example:
//
//	// Delete a block list that's no longer needed
//	result, err := client.DeleteBlockList(ctx, blockListID)
//	if err != nil {
//	    return err
//	}
//	if result.IsSuccessful() {
//	    fmt.Println("Block list deleted successfully")
//	}
func (c *Client) DeleteBlockList(ctx context.Context, blockListID int) (*types.DeleteBlockListResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if blockListID <= 0 {
		return nil, WrapError("DeleteBlockList", ErrInvalidInput, "block list ID must be positive")
	}

	var mutation struct {
		Response struct {
			Status  string `graphql:"status"`
			Message string `graphql:"message"`
			Error   string `graphql:"error"`
		} `graphql:"deleteBlockList(blocklist_id: $blocklist_id)"`
	}

	variables := map[string]interface{}{
		"blocklist_id": blockListID,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("DeleteBlockList", err, "failed to delete block list")
	}

	// Map to types.DeleteBlockListResponse
	response := &types.DeleteBlockListResponse{
		Status:  mutation.Response.Status,
		Message: mutation.Response.Message,
		Error:   mutation.Response.Error,
	}

	// Check for error in response
	if !response.IsSuccessful() {
		return response, WrapError("DeleteBlockList", ErrOperationFailed, response.Error)
	}

	return response, nil
}

// DeleteBlockListEntry deletes one or more entries from block lists.
// This allows removing specific IPs, domains, or user agents from block lists
// without deleting the entire list.
//
// Multiple entries can be deleted in a single call by providing multiple entry IDs.
// This is more efficient than deleting entries one at a time.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - entryIDs: List of block list entry IDs to delete
//
// Returns:
//   - *types.DeleteBlockListEntryResponse: Result with count of deleted entries
//   - error: Error if the operation fails
//
// Example:
//
//	// Delete specific entries from a block list
//	entryIDs := []int{123, 456, 789}
//	result, err := client.DeleteBlockListEntry(ctx, entryIDs)
//	if err != nil {
//	    return err
//	}
//	if result.IsSuccessful() {
//	    fmt.Printf("Deleted %d entries\n", result.DeletedCount)
//	}
func (c *Client) DeleteBlockListEntry(ctx context.Context, entryIDs []int) (*types.DeleteBlockListEntryResponse, error) {
	if err := c.EnsureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Validate input
	if len(entryIDs) == 0 {
		return nil, WrapError("DeleteBlockListEntry", ErrInvalidInput, "entry IDs list cannot be empty")
	}

	// Validate all IDs are positive
	for i, id := range entryIDs {
		if id <= 0 {
			return nil, WrapError("DeleteBlockListEntry", ErrInvalidInput, "all entry IDs must be positive")
		}
		// Check for duplicates
		for j := i + 1; j < len(entryIDs); j++ {
			if entryIDs[i] == entryIDs[j] {
				return nil, WrapError("DeleteBlockListEntry", ErrInvalidInput, "duplicate entry IDs not allowed")
			}
		}
	}

	var mutation struct {
		Response struct {
			Status       string `graphql:"status"`
			DeletedCount int    `graphql:"deleted_count"`
			Message      string `graphql:"message"`
			Error        string `graphql:"error"`
		} `graphql:"deleteBlockListEntry(entry_ids: $entry_ids)"`
	}

	variables := map[string]interface{}{
		"entry_ids": entryIDs,
	}

	if err := c.executeMutation(ctx, &mutation, variables); err != nil {
		return nil, WrapError("DeleteBlockListEntry", err, "failed to delete block list entries")
	}

	// Map to types.DeleteBlockListEntryResponse
	response := &types.DeleteBlockListEntryResponse{
		Status:       mutation.Response.Status,
		DeletedCount: mutation.Response.DeletedCount,
		Message:      mutation.Response.Message,
		Error:        mutation.Response.Error,
	}

	// Check for error in response
	if !response.IsSuccessful() {
		return response, WrapError("DeleteBlockListEntry", ErrOperationFailed, response.Error)
	}

	return response, nil
}
