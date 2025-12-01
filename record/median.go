package record

import (
	"sort"

	api "github.com/tilotech/tilores-plugin-api"
)

// Median returns the median value of the provided numeric path.
//
// Using median on non-numeric paths will raise an error.
// Null values are ignored in the calculation.
// Returns null if all values are null.
func Median(records []*api.Record, path string) (*float64, error) {
	numbers := []float64{}
	err := VisitNumber(records, path, func(number *float64, _ *api.Record) error {
		if number != nil {
			numbers = append(numbers, *number)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	counted := len(numbers)
	if counted == 0 {
		return nil, nil
	}
	sort.Float64s(numbers)
	if counted%2 == 1 {
		return pointer(numbers[(counted / 2)]), nil
	}
	return pointer((numbers[counted/2] + numbers[(counted/2)-1]) / 2), nil
}
