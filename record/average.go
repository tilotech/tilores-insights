package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Average returns the average value of the provided numeric path.
//
// Using average on non-numeric paths will raise an error.
// Null values are ignored in the calculation.
// Returns null if all values are null.
func Average(records []*api.Record, path string) (*float64, error) {
	if len(records) == 0 {
		return nil, nil
	}
	sum := 0.0
	counted := 0.0
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			sum += *number
			counted++
		}
	}
	if counted == 0 {
		return nil, nil
	}
	return pointer(sum / counted), nil
}
