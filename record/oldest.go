package record

import (
	"time"

	api "github.com/tilotech/tilores-plugin-api"
)

// Oldest returns the Record for where the provided time path has the lowest
// (least recent) value.
//
// Returns null if the list is empty or does not contain records with the
// provided path.
//
// Using oldest on non-time paths will raise an error.
func Oldest(records []*api.Record, path string) (*api.Record, error) {
	var record *api.Record
	var oldestTime *time.Time

	for _, r := range records {
		t, err := ExtractTime(r, path)
		if err != nil {
			return nil, err
		}
		if t == nil {
			continue
		}
		if oldestTime == nil || oldestTime.After(*t) {
			oldestTime = t
			record = r
		}
	}
	return record, nil
}
