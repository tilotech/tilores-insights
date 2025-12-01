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

	err := VisitTime(records, path, func(t *time.Time, r *api.Record) error {
		if t == nil {
			return nil
		}
		if oldestTime == nil || oldestTime.After(*t) {
			oldestTime = t
			record = r
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}
