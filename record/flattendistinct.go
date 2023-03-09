package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// FlattenDistinct merges the array on the given path of all records into a
// single array where each value is unique.
//
// Using flatten on non-array fields will raise an error.
//
// By default, the case of the value is ignored.
func FlattenDistinct(records []*api.Record, path string, caseSensitive bool) ([]any, error) {
	result := []any{}
	unique := map[string]struct{}{}
	for _, record := range records {
		array, err := ExtractArray(record, path)
		if err != nil {
			return nil, err
		}
		for _, e := range array {
			if e == nil {
				continue
			}
			val, err := valueToString(e, caseSensitive)
			if err != nil {
				return nil, err
			}
			if _, ok := unique[*val]; !ok {
				unique[*val] = struct{}{}
				result = append(result, e)
			}
		}
	}
	return result, nil
}
