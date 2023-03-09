package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Flatten merges the array on the given path of all records into a single array.
//
// Using flatten on non-array fields will raise an error.
func Flatten(records []*api.Record, path string) ([]any, error) {
	result := []any{}
	for _, record := range records {
		array, err := ExtractArray(record, path)
		if err != nil {
			return nil, err
		}
		for _, e := range array {
			if e != nil {
				result = append(result, e)
			}
		}
	}
	return result, nil
}
