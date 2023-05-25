package edge

import (
	"sort"
	"strings"

	api "github.com/tilotech/tilores-plugin-api"
)

// FrequencyDistributionEntry represents a single row of a frequency distribution table.
type FrequencyDistributionEntry struct {
	RuleID     string  `json:"value"`
	Frequency  int     `json:"frequency"`
	Percentage float64 `json:"percentage"`
}

// FrequencyDistribution returns how often a rule is present.
//
// By default, the results are ordered with the highest percentage first, but
// it can be changed using the 'sortASC' option.
//
// Using the 'top' option it is possible to limit the results to only the n
// highest or lowest results.
//
// This function does not consider implicit rule usages based on duplicates.
func FrequencyDistribution(edges api.Edges, top int, sortASC bool) []*FrequencyDistributionEntry {
	if top == 0 {
		return []*FrequencyDistributionEntry{}
	}

	freq := map[string]int{}
	counted := 0
	for _, edge := range edges {
		rule := strings.SplitN(edge, ":", 3)[2]
		freq[rule]++
		counted++
	}

	result := make([]*FrequencyDistributionEntry, 0, len(freq))
	if counted == 0 {
		return result
	}

	for rule, f := range freq {
		result = append(result, &FrequencyDistributionEntry{
			RuleID:     rule,
			Frequency:  f,
			Percentage: float64(f) / float64(counted),
		})
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
		return result[0:top]
	}
	return result
}
