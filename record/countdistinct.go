package record

import (
	api "github.com/tilotech/tilores-plugin-api"
	"gitlab.com/tilotech/go-helpers/stringset"
)

// CountDistinct returns the number of unique non-null values for the provided
// path.
//
// If multiple paths were provided, then each unique combination of the path
// values will be considered.
// If all paths are null, then this does not count as a new value. However, if
// at least one path has a value, then this does count as a new value.
func CountDistinct(records []*api.Record, paths []string, caseSensitive bool) (int, error) {
	recordSet := stringset.New()
	for _, record := range records {
		recordVal := ""
		foundAVal := false
		for _, path := range paths {
			val, err := ExtractString(record, path, caseSensitive)
			if err != nil {
				return 0, err
			}
			if val != nil {
				foundAVal = true
				recordVal += *val
			}
		}
		if foundAVal {
			recordSet.Add(recordVal)
		}
	}
	return len(recordSet), nil
}
