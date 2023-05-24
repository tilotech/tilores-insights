package edge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/edge"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestFrequencyDistribution(t *testing.T) {
	cases := map[string]struct {
		edges    api.Edges
		top      int
		sortASC  bool
		expected []*edge.FrequencyDistributionEntry
	}{
		"no edges": {
			edges:    api.Edges{},
			top:      -1,
			expected: []*edge.FrequencyDistributionEntry{},
		},
		"single edge": {
			edges: api.Edges{
				"record-1:record-2:R1",
			},
			top: -1,
			expected: []*edge.FrequencyDistributionEntry{
				{
					RuleID:     "R1",
					Frequency:  1,
					Percentage: 1.0,
				},
			},
		},
		"multiple edges": {
			edges: api.Edges{
				"record-1:record-2:R1",
				"record-1:record-2:R2",
				"record-1:record-3:R1",
				"record-1:record-3:R2",
				"record-1:record-4:R1",
			},
			top: -1,
			expected: []*edge.FrequencyDistributionEntry{
				{
					RuleID:     "R1",
					Frequency:  3,
					Percentage: 0.6,
				},
				{
					RuleID:     "R2",
					Frequency:  2,
					Percentage: 0.4,
				},
			},
		},
		"multiple edges, sort asc": {
			edges: api.Edges{
				"record-1:record-2:R1",
				"record-1:record-2:R2",
				"record-1:record-3:R1",
				"record-1:record-3:R2",
				"record-1:record-4:R1",
			},
			top:     -1,
			sortASC: true,
			expected: []*edge.FrequencyDistributionEntry{
				{
					RuleID:     "R2",
					Frequency:  2,
					Percentage: 0.4,
				},
				{
					RuleID:     "R1",
					Frequency:  3,
					Percentage: 0.6,
				},
			},
		},
		"multiple edges, top 0": {
			edges: api.Edges{
				"record-1:record-2:R1",
				"record-1:record-2:R2",
				"record-1:record-3:R1",
				"record-1:record-3:R2",
				"record-1:record-4:R1",
			},
			top:      0,
			expected: []*edge.FrequencyDistributionEntry{},
		},
		"multiple edges, top 1": {
			edges: api.Edges{
				"record-1:record-2:R1",
				"record-1:record-2:R2",
				"record-1:record-3:R1",
				"record-1:record-3:R2",
				"record-1:record-4:R1",
			},
			top: 1,
			expected: []*edge.FrequencyDistributionEntry{
				{
					RuleID:     "R1",
					Frequency:  3,
					Percentage: 0.6,
				},
			},
		},
		"multiple edges, top 10": {
			edges: api.Edges{
				"record-1:record-2:R1",
				"record-1:record-2:R2",
				"record-1:record-3:R1",
				"record-1:record-3:R2",
				"record-1:record-4:R1",
			},
			top: 10,
			expected: []*edge.FrequencyDistributionEntry{
				{
					RuleID:     "R1",
					Frequency:  3,
					Percentage: 0.6,
				},
				{
					RuleID:     "R2",
					Frequency:  2,
					Percentage: 0.4,
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := edge.FrequencyDistribution(c.edges, c.top, c.sortASC)
			assert.Equal(t, c.expected, actual)
		})
	}
}
