package record

import (
	"math"

	api "github.com/tilotech/tilores-plugin-api"
)

// StandardDeviation calculates the standard deviation for the provided numeric
// path.
//
// Using standardDeviation on non-numeric paths will raise an error.
// Null values are ignored in the calculation.
// Returns null if all values are null.
func StandardDeviation(records []*api.Record, path string) (*float64, error) {
	if len(records) == 0 {
		return nil, nil
	}
	avg, err := Average(records, path)
	if err != nil {
		return nil, err
	}
	difSquareSum := 0.0
	counted := 0.0
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			dif := *number - *avg
			difSquareSum += dif * dif
			counted++
		}
	}
	if counted == 0 {
		return nil, nil
	}
	return pointer(math.Sqrt(difSquareSum / counted)), nil
}

// TODO: Add text standard deviation
