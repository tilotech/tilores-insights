package record

import api "github.com/tilotech/tilores-plugin-api"

// Last returns the last record from the provided list of records.
func Last(records []*api.Record) *api.Record {
	if len(records) == 0 {
		return nil
	}
	return records[len(records)-1]
}
