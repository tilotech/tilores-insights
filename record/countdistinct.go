package record

import (
	"strings"

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
	recordKeys, err := extractStringKeys(records, paths, caseSensitive)
	if err != nil {
		return 0, err
	}
	keys := map[string]struct{}{}
	for _, key := range recordKeys {
		if !strings.Contains(key, ":s:") {
			continue
		}
		keys[key] = struct{}{}
	}
	return len(keys), nil
}
