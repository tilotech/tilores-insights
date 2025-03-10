package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Max returns the highest value of the provided numeric path.
//
// Returns null if all values are null.
func Max(records []*api.Record, path string) (*float64, error) {
	var maxVal *float64
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			if maxVal == nil || *number > *maxVal {
				maxVal = number
			}
		}
	}
	return maxVal, nil
}
