package record

import api "github.com/tilotech/tilores-plugin-api"

// Values returns all non-null values of the current records.
func Values(records []*api.Record, path string) []any {
	result := make([]any, 0, len(records))
	for _, record := range records {
		val := Extract(record, path)
		if val != nil {
			result = append(result, val)
		}
	}
	return result
}
