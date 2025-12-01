package record

import (
	"cmp"

	api "github.com/tilotech/tilores-plugin-api"
	"golang.org/x/exp/slices"
)

// FrequencyDistributionEntry represents a single row of a frequency distribution table.
type FrequencyDistributionEntry struct {
	Value       any     `json:"value"`
	Frequency   int     `json:"frequency"`
	Percentage  float64 `json:"percentage"`
	originalPos int
}

// FrequencyDistribution returns how often a non-null value for the provided
// field is present.
//
// By default, the results are ordered with the highest percentage first, but
// it can be changed using the 'sortASC' option.
//
// Using the 'top' option it is possible to limit the results to only the n
// highest or lowest results.
//
// Values with with equal frequency will always be returned in the order of the
// first occurrence for that value.
func FrequencyDistribution(records []*api.Record, path string, caseSensitive bool, top int, sortASC bool) ([]*FrequencyDistributionEntry, error) { //nolint:gocognit
	if top == 0 {
		return []*FrequencyDistributionEntry{}, nil
	}
	positionMap := map[*api.Record]int{}
	for i := range records {
		positionMap[records[i]] = i
	}
	entriesMap := make(map[string]*FrequencyDistributionEntry, len(records))
	counted := 0
	err := Visit(records, path, func(v any, record *api.Record) error {
		val, err := validateString(v, caseSensitive)
		if err != nil {
			return err
		}
		if val != nil {
			counted++
			if _, ok := entriesMap[*val]; !ok {
				entriesMap[*val] = &FrequencyDistributionEntry{
					Value:       v,
					Frequency:   1,
					Percentage:  0.0,
					originalPos: positionMap[record],
				}
			} else {
				entriesMap[*val].Frequency++
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := make([]*FrequencyDistributionEntry, 0, len(entriesMap))
	if counted == 0 {
		return result, nil
	}
	for _, entry := range entriesMap {
		entry.Percentage = float64(entry.Frequency) / float64(counted)
		result = append(result, entry)
	}
	sortFunc := func(a, b *FrequencyDistributionEntry) int {
		if a.Frequency == b.Frequency {
			return cmp.Compare(a.originalPos, b.originalPos)
		}
		return cmp.Compare(b.Frequency, a.Frequency)
	}
	if sortASC {
		sortFunc = func(a, b *FrequencyDistributionEntry) int {
			if a.Frequency == b.Frequency {
				return cmp.Compare(a.originalPos, b.originalPos)
			}
			return cmp.Compare(a.Frequency, b.Frequency)
		}
	}
	slices.SortStableFunc(result, sortFunc)
	if top > 0 && top < len(result) {
		return result[0:top], nil
	}
	return result, nil
}
