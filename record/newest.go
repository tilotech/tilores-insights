package record

import (
	"time"

	api "github.com/tilotech/tilores-plugin-api"
)

// Newest returns the Record for where the provided time path has the highest
// (most recent) value.
//
// Returns null if the list is empty or does not contain records with the
// provided path.
//
// Using newest on non-time paths will raise an error.
func Newest(records []*api.Record, path string) (*api.Record, error) {
	var record *api.Record
	var newestTime *time.Time

	for _, r := range records {
		t, err := ExtractTime(r, path)
		if err != nil {
			return nil, err
		}
		if t == nil {
			continue
		}
		if newestTime == nil || newestTime.Before(*t) {
			newestTime = t
			record = r
		}
	}
	return record, nil
}
