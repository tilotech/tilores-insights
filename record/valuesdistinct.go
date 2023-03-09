package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// ValuesDistinct returns all unique non-null values of the current records.
// By default, the case of the value is ignored.
func ValuesDistinct(records []*api.Record, path string, caseSensitive bool) ([]any, error) {
	result := make([]any, 0, len(records))
	unique := make(map[string]struct{}, len(records))
	for _, record := range records {
		val, err := ExtractString(record, path, caseSensitive)
		if err != nil {
			return nil, err
		}
		if val != nil {
			if _, ok := unique[*val]; !ok {
				unique[*val] = struct{}{}
				result = append(result, Extract(record, path))
			}
		}
	}
	return result, nil
}
