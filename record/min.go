package record

import (
	"math"

	api "github.com/tilotech/tilores-plugin-api"
)

// Min returns the lowest value of the provided numeric path.
//
// Using min on non-numeric paths will raise an error.
// Returns null if all values are null.
func Min(records []*api.Record, path string) (*float64, error) {
	if len(records) == 0 {
		return nil, nil
	}
	min := 0.0
	counted := 0.0
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			if counted == 0 {
				min = *number
			} else {
				min = math.Min(min, *number)
			}
			counted++
		}
	}
	if counted == 0 {
		return nil, nil
	}
	return pointer(min), nil
}
