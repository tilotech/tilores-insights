package record

import (
	"sort"

	api "github.com/tilotech/tilores-plugin-api"
)

// FrequencyDistributionEntry represents a single row of a frequency distribution table.
type FrequencyDistributionEntry struct {
	Value      any     `json:"value"`
	Frequency  int     `json:"frequency"`
	Percentage float64 `json:"percentage"`
}

// FrequencyDistribution returns how often a non-null value for the provided
// field is present.
//
// By default, the results are ordered with the highest percentage first, but
// it can be changed using the 'sortASC' option.
//
// Using the 'top' option it is possible to limit the results to only the n
// highest or lowest results.
func FrequencyDistribution(records []*api.Record, path string, caseSensitive bool, top int, sortASC bool) ([]*FrequencyDistributionEntry, error) {
	entriesMap := make(map[string]*FrequencyDistributionEntry, len(records))
	counted := 0
	for _, record := range records {
		val, err := ExtractString(record, path, caseSensitive)
		if err != nil {
			return nil, err
		}
		if val != nil {
			counted++
			if _, ok := entriesMap[*val]; !ok {
				entriesMap[*val] = &FrequencyDistributionEntry{
					Value:      Extract(record, path),
					Frequency:  1,
					Percentage: 0.0,
				}
			} else {
				entriesMap[*val].Frequency++
			}
		}
	}
	result := make([]*FrequencyDistributionEntry, 0, len(entriesMap))
	if counted == 0 {
		return result, nil
	}
	for _, entry := range entriesMap {
		entry.Percentage = float64(entry.Frequency) / float64(counted)
		result = append(result, entry)
	}
	sortFunc := func(i, j int) bool {
		return result[i].Percentage > result[j].Percentage
	}
	if sortASC {
		sortFunc = func(i, j int) bool {
			return result[i].Percentage < result[j].Percentage
		}
	}
	sort.Slice(result, sortFunc)
	if top > 0 && top < len(result) {
		return result[0:top], nil
	}
	return result, nil
}
