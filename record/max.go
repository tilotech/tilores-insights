package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Max returns the highest value of the provided numeric path.
//
// Returns null if all values are null.
func Max(records []*api.Record, path string) (*float64, error) {
	var max *float64
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			if max == nil || *number > *max {
				max = number
			}
		}
	}
	return max, nil
}
