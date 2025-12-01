package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Sum returns the sum of the provided numeric path.
//
// Using sum on non-numeric paths will raise an error.
// Null values are ignored in the calculation.
// Returns null if all values are null.
func Sum(records []*api.Record, path string) (*float64, error) {
	sum := 0.0
	counted := 0.0
	err := VisitNumber(records, path, func(number *float64, _ *api.Record) error {
		if number != nil {
			sum += *number
			counted++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if counted == 0 {
		return nil, nil
	}
	return pointer(sum), nil
}
