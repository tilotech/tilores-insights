package record

import (
	"math"

	api "github.com/tilotech/tilores-plugin-api"
)

// StandardDeviationText calculates the standard deviation for the value
// frequency of the provided path.
//
// Null values are ignored in the calculation.
// Returns null if all values are null.
func StandardDeviationText(records []*api.Record, path string, caseSensitive bool) (*float64, error) {
	if len(records) == 0 {
		return nil, nil
	}
	frequencies := make(map[string]float64, len(records))
	counted := 0.0
	for _, record := range records {
		val, err := ExtractString(record, path, caseSensitive)
		if err != nil {
			return nil, err
		}
		if val != nil {
			counted++
			if _, ok := frequencies[*val]; !ok {
				frequencies[*val] = 1.0
			} else {
				frequencies[*val]++
			}
		}
	}
	if counted == 0.0 {
		return nil, nil
	}

	sum := 0.0
	for _, v := range frequencies {
		sum += v
	}
	count := float64(len(frequencies))
	avg := sum / count

	difSquareSum := 0.0
	for _, v := range frequencies {
		dif := v - avg
		difSquareSum += dif * dif
	}

	return pointer(math.Sqrt(difSquareSum / count)), nil
}
