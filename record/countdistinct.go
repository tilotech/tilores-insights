package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// CountDistinct returns the number of unique non-null values for the provided
// path.
//
// If multiple paths were provided, then each unique combination of the path
// values will be considered.
// If all paths are null, then this does not count as a new value. However, if
// at least one path has a value, then this does count as a new value.
func CountDistinct(records []*api.Record, paths []string, caseSensitive bool) (int, error) {
	nilKey, _ := extractStringKeys(&api.Record{Data: map[string]any{}}, paths, caseSensitive)

	keys := map[string]struct{}{}
	for _, record := range records {
		key, err := extractStringKeys(record, paths, caseSensitive)
		if err != nil {
			return 0, err
		}
		if key != nilKey {
			keys[key] = struct{}{}
		}
	}
	return len(keys), nil
}
