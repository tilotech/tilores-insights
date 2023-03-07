package record

import (
	"sort"

	api "github.com/tilotech/tilores-plugin-api"
)

type FrequencyDistributionEntry struct {
	Value      any     `json:"value"`
	Frequency  int     `json:"frequency"`
	Percentage float64 `json:"percentage"`
}

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
	if top > 0 {
		return result[0:top], nil
	}
	return result, nil
}
