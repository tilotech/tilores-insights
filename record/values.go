package record

import api "github.com/tilotech/tilores-plugin-api"

// Values returns all non-null values of the current records.
func Values(records []*api.Record, path string) []any {
	result := make([]any, 0, len(records))
	_ = Visit(records, path, func(val any, _ *api.Record) error {
		if val != nil {
			result = append(result, val)
		}
		return nil
	})
	return result
}
