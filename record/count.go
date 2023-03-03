package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Count returns the amount of records in the provided list.
//
// Entries with nil values will be counted just as any non-nil entry.
func Count(records []*api.Record) int {
	return len(records)
}
