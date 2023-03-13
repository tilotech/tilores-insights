package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Limit returns a new record list that contains up to 'count' records.
//
// By default it takes the first records from the list. If offset was provided
// it will skip 'offset' records.
//
// If the list does not provide enough records, then an empty list is returned.
func Limit(records []*api.Record, count int, offset int) []*api.Record {
	l := len(records)
	if count <= 0 || offset < 0 || offset >= l {
		return []*api.Record{}
	}
	n := min(count, l-offset)
	return records[offset : offset+n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
