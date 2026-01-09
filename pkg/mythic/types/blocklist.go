package types

import "fmt"

// BlockList represents a collection of blocked IPs, domains, or user agents.
type BlockList struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OperationID int    `json:"operation_id"`
	Active      bool   `json:"active"`
	Deleted     bool   `json:"deleted"`
}

// String returns a human-readable representation of the block list.
func (b *BlockList) String() string {
	status := "active"
	if b.Deleted {
		status = "deleted"
	} else if !b.Active {
		status = "inactive"
	}
	return fmt.Sprintf("BlockList '%s' (%s)", b.Name, status)
}

// IsActive returns true if the block list is active and not deleted.
func (b *BlockList) IsActive() bool {
	return b.Active && !b.Deleted
}

// IsDeleted returns true if the block list has been deleted.
func (b *BlockList) IsDeleted() bool {
	return b.Deleted
}

// BlockListEntry represents an individual entry in a block list.
type BlockListEntry struct {
	ID          int    `json:"id"`
	BlockListID int    `json:"blocklist_id"`
	Type        string `json:"type"` // "ip", "domain", "user_agent"
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Active      bool   `json:"active"`
	Deleted     bool   `json:"deleted"`
}

// String returns a human-readable representation of the block list entry.
func (b *BlockListEntry) String() string {
	status := "active"
	if b.Deleted {
		status = "deleted"
	} else if !b.Active {
		status = "inactive"
	}
	return fmt.Sprintf("%s: %s (%s)", b.Type, b.Value, status)
}

// IsActive returns true if the entry is active and not deleted.
func (b *BlockListEntry) IsActive() bool {
	return b.Active && !b.Deleted
}

// IsDeleted returns true if the entry has been deleted.
func (b *BlockListEntry) IsDeleted() bool {
	return b.Deleted
}

// DeleteBlockListRequest represents a request to delete a block list.
type DeleteBlockListRequest struct {
	BlockListID int `json:"blocklist_id"`
}

// String returns a human-readable representation of the request.
func (d *DeleteBlockListRequest) String() string {
	return fmt.Sprintf("Delete block list ID %d", d.BlockListID)
}

// DeleteBlockListResponse represents the response from deleting a block list.
type DeleteBlockListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (d *DeleteBlockListResponse) String() string {
	if d.Status == "success" {
		if d.Message != "" {
			return d.Message
		}
		return "Block list deleted successfully"
	}
	return fmt.Sprintf("Failed to delete block list: %s", d.Error)
}

// IsSuccessful returns true if the deletion succeeded.
func (d *DeleteBlockListResponse) IsSuccessful() bool {
	return d.Status == "success"
}

// DeleteBlockListEntryRequest represents a request to delete block list entries.
type DeleteBlockListEntryRequest struct {
	EntryIDs []int `json:"entry_ids"`
}

// String returns a human-readable representation of the request.
func (d *DeleteBlockListEntryRequest) String() string {
	return fmt.Sprintf("Delete %d block list entries", len(d.EntryIDs))
}

// DeleteBlockListEntryResponse represents the response from deleting entries.
type DeleteBlockListEntryResponse struct {
	Status       string `json:"status"`
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
}

// String returns a human-readable representation of the response.
func (d *DeleteBlockListEntryResponse) String() string {
	if d.Status == "success" {
		return fmt.Sprintf("Deleted %d block list entries", d.DeletedCount)
	}
	return fmt.Sprintf("Failed to delete entries: %s", d.Error)
}

// IsSuccessful returns true if the deletion succeeded.
func (d *DeleteBlockListEntryResponse) IsSuccessful() bool {
	return d.Status == "success"
}
