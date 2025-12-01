package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Max returns the highest value of the provided numeric path.
//
// Returns null if all values are null.
func Max(records []*api.Record, path string) (*float64, error) {
	var maxVal *float64
	err := VisitNumber(records, path, func(number *float64, _ *api.Record) error {
		if number != nil {
			if maxVal == nil || *number > *maxVal {
				maxVal = number
			}
		}
		return nil
	})
	return maxVal, err
}
