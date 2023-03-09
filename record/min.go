package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Min returns the lowest value of the provided numeric path.
//
// Using min on non-numeric paths will raise an error.
// Returns null if all values are null.
func Min(records []*api.Record, path string) (*float64, error) {
	var min *float64
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			if min == nil || *number < *min {
				min = number
			}
		}
	}
	return min, nil
}
