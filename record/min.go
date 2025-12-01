package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Min returns the lowest value of the provided numeric path.
//
// Using min on non-numeric paths will raise an error.
// Returns null if all values are null.
func Min(records []*api.Record, path string) (*float64, error) {
	var minVal *float64
	err := VisitNumber(records, path, func(number *float64, _ *api.Record) error {
		if number != nil {
			if minVal == nil || *number < *minVal {
				minVal = number
			}
		}
		return nil
	})
	return minVal, err
}
