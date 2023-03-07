package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// ValuesDistinct returns all unique non-null values of the current records.
// By default, the case of the value is ignored.
func ValuesDistinct(records []*api.Record, path string, caseSensitive bool) ([]any, error) {
	unique := make(map[string]any, len(records))
	for _, record := range records {
		val, err := ExtractString(record, path, caseSensitive)
		if err != nil {
			return nil, err
		}
		if val != nil {
			if _, ok := unique[*val]; !ok {
				unique[*val] = Extract(record, path)
			}
		}
	}
	result := make([]any, 0, len(unique))
	for _, val := range unique {
		result = append(result, val)
	}
	return result, nil
}
