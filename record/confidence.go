package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Confidence describes the probability of having the one truly correct value
// for the provided path.
//
// The resulting value is a float ranging from 0 to 1 representing a percentage.
//
// Example: For the values ["a","a","a","b"]
//
// a: 3 * 0.75
//
// b: 1 * 0.25
//
// confidence: 62.5%
//
// Null values are ignored in the calculation.
// Returns null if all values are null.
func Confidence(records []*api.Record, path string, caseSensitive bool) (*float64, error) {
	if len(records) == 0 {
		return nil, nil
	}
	frequencies := make(map[string]int, len(records))
	valueCount := 0

	for _, record := range records {
		val, err := ExtractString(record, path, caseSensitive)
		if err != nil {
			return nil, err
		}
		if val != nil {
			valueCount++
			if _, ok := frequencies[*val]; !ok {
				frequencies[*val] = 1
			} else {
				frequencies[*val]++
			}
		}
	}
	if valueCount == 0 {
		return nil, nil
	}
	probSum := 0.0
	for _, freq := range frequencies {
		prob := float64(freq) / float64(valueCount)
		probSum += prob * float64(freq)
	}
	return pointer(probSum / float64(valueCount)), nil
}
