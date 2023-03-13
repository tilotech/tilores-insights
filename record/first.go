package record

import api "github.com/tilotech/tilores-plugin-api"

// First returns the first record from the provided list of records.
func First(records []*api.Record) *api.Record {
	if len(records) == 0 {
		return nil
	}
	return records[0]
}
