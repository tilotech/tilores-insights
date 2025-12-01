package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Flatten merges the array on the given path of all records into a single array.
//
// Using flatten on non-array fields will raise an error.
func Flatten(records []*api.Record, path string) ([]any, error) {
	result := []any{}
	err := VisitArray(records, path, func(array []any, _ *api.Record) error {
		for _, e := range array {
			if e != nil {
				result = append(result, e)
			}
		}
		return nil
	})
	return result, err
}
