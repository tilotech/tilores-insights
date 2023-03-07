package record

import api "github.com/tilotech/tilores-plugin-api"
import "gitlab.com/tilotech/go-helpers/stringset"

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
