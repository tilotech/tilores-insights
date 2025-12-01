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
	recordKeys := make(map[string][]string)
	for _, record := range records {
		if record == nil {
			continue
		}
		recordKeys[record.ID] = []string{""}
	}
	for _, path := range paths {
		keys := map[string][]*string{}
		err := VisitString(records, path, caseSensitive, func(val *string, record *api.Record) error {
			if record == nil {
				return nil
			}
			if _, ok := keys[record.ID]; !ok {
				keys[record.ID] = make([]*string, 0, 1)
			}
			keys[record.ID] = append(keys[record.ID], val)
			return nil
		})
		if err != nil {
			return 0, err
		}
		for rid, values := range keys {
			recordKeys[rid] = mergePathKeys(recordKeys[rid], values)
		}
	}
	return collectDistinctKeyCount(recordKeys), nil
}

func mergePathKeys(existingKeys []string, values []*string) []string {
	newKeys := make([]string, 0, len(existingKeys)*len(values))
	for _, existingKey := range existingKeys {
		for _, val := range values {
			var newKey string
			if val == nil {
				newKey = existingKey + ":n:"
			} else {
				newKey = existingKey + ":s:" + *val
			}
			newKeys = append(newKeys, newKey)
		}
	}
	return newKeys
}

func collectDistinctKeyCount(recordKeys map[string][]string) int {
	allKeys := make(map[string]struct{})
	for _, keys := range recordKeys {
		for _, key := range keys {
			if _, ok := allKeys[key]; ok || !strings.Contains(key, ":s:") {
				continue
			}
			allKeys[key] = struct{}{}
		}
	}
	return len(allKeys)
}
