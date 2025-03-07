package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestFrequencyDistribution(t *testing.T) {
	testRecords := []*api.Record{
		{
			ID: "someid",
			Data: map[string]any{
				"field1": "c",
				"case":   "a",
				"nested": map[string]any{
					"field1": "a",
					"field2": "b",
					"field3": "c",
				},
				"nestedCase": map[string]any{
					"field1": "a",
					"field2": "a",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"field1": "a",
				"case":   "a",
				"nested": map[string]any{
					"field1": "a",
					"field2": "a",
				},
				"nestedCase": map[string]any{
					"field1": "a",
					"field2": "a",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"field1": "A",
				"case":   "A",
				"nested": map[string]any{
					"field1": "A",
					"field2": "A",
				},
				"nestedCase": map[string]any{
					"field1": "A",
					"field2": "A",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"field1": "b",
				"case":   "b",
				"nested": map[string]any{
					"field2": "b",
				},
				"nestedCase": map[string]any{
					"field2": "b",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"field1": "b",
				"case":   "b",
				"nested": map[string]any{
					"field2": "b",
				},
				"nestedCase": map[string]any{
					"field2": "b",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"field1": "a",
				"case":   "a",
				"nested": map[string]any{
					"field1": "a",
					"field2": "a",
				},
				"nestedCase": map[string]any{
					"field1": "a",
					"field2": "a",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"field1": "d",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"field1": "d",
			},
		},
		{
			ID:   "someid",
			Data: map[string]any{},
		},
		nil,
	}
	cases := map[string]struct {
		records       []*api.Record
		path          string
		caseSensitive bool
		top           int
		sortASC       bool
		expected      []*record.FrequencyDistributionEntry
	}{
		"empty list": {
			records:  []*api.Record{},
			path:     "field1",
			top:      -1,
			expected: []*record.FrequencyDistributionEntry{},
		},
		"nil list": {
			records:  nil,
			path:     "field1",
			top:      -1,
			expected: []*record.FrequencyDistributionEntry{},
		},
		"on field": {
			records: testRecords,
			path:    "field1",
			top:     -1,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "a",
					Frequency:  3,
					Percentage: 0.375,
				},
				{
					Value:      "b",
					Frequency:  2,
					Percentage: 0.25,
				},
				{
					Value:      "d",
					Frequency:  2,
					Percentage: 0.25,
				},
				{
					Value:      "c",
					Frequency:  1,
					Percentage: 0.125,
				},
			},
		},
		"on field top more than actual number of entries": {
			records: testRecords,
			path:    "field1",
			top:     5,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "a",
					Frequency:  3,
					Percentage: 0.375,
				},
				{
					Value:      "b",
					Frequency:  2,
					Percentage: 0.25,
				},
				{
					Value:      "d",
					Frequency:  2,
					Percentage: 0.25,
				},
				{
					Value:      "c",
					Frequency:  1,
					Percentage: 0.125,
				},
			},
		},
		"on field top 1": {
			records: testRecords,
			path:    "field1",
			top:     1,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "a",
					Frequency:  3,
					Percentage: 0.375,
				},
			},
		},
		"on field asc": {
			records: testRecords,
			path:    "field1",
			sortASC: true,
			top:     -1,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "c",
					Frequency:  1,
					Percentage: 0.125,
				},
				{
					Value:      "b",
					Frequency:  2,
					Percentage: 0.25,
				},
				{
					Value:      "d",
					Frequency:  2,
					Percentage: 0.25,
				},
				{
					Value:      "a",
					Frequency:  3,
					Percentage: 0.375,
				},
			},
		},
		"on object": {
			records: testRecords,
			path:    "nested",
			top:     -1,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value: map[string]any{
						"field1": "a",
						"field2": "a",
					},
					Frequency:  3,
					Percentage: 0.5,
				},
				{
					Value: map[string]any{
						"field2": "b",
					},
					Frequency:  2,
					Percentage: 0.3333333333333333,
				},
				{
					Value: map[string]any{
						"field1": "a",
						"field2": "b",
						"field3": "c",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
			},
		},
		"on object top 2": {
			records: testRecords,
			path:    "nested",
			top:     2,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value: map[string]any{
						"field1": "a",
						"field2": "a",
					},
					Frequency:  3,
					Percentage: 0.5,
				},
				{
					Value: map[string]any{
						"field2": "b",
					},
					Frequency:  2,
					Percentage: 0.3333333333333333,
				},
			},
		},
		"on object asc": {
			records: testRecords,
			path:    "nested",
			sortASC: true,
			top:     -1,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value: map[string]any{
						"field1": "a",
						"field2": "b",
						"field3": "c",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]any{
						"field2": "b",
					},
					Frequency:  2,
					Percentage: 0.3333333333333333,
				},
				{
					Value: map[string]any{
						"field1": "a",
						"field2": "a",
					},
					Frequency:  3,
					Percentage: 0.5,
				},
			},
		},
		"case sensitive": {
			records:       testRecords,
			path:          "case",
			caseSensitive: true,
			top:           -1,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "a",
					Frequency:  3,
					Percentage: 0.5,
				},
				{
					Value:      "b",
					Frequency:  2,
					Percentage: 0.3333333333333333,
				},
				{
					Value:      "A",
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
			},
		},
		"case sensitive asc": {
			records:       testRecords,
			path:          "case",
			caseSensitive: true,
			sortASC:       true,
			top:           -1,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "A",
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value:      "b",
					Frequency:  2,
					Percentage: 0.3333333333333333,
				},
				{
					Value:      "a",
					Frequency:  3,
					Percentage: 0.5,
				},
			},
		},
		"on object case sensitive": {
			records:       testRecords,
			path:          "nestedCase",
			caseSensitive: true,
			top:           -1,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value: map[string]any{
						"field1": "a",
						"field2": "a",
					},
					Frequency:  3,
					Percentage: 0.5,
				},
				{
					Value: map[string]any{
						"field2": "b",
					},
					Frequency:  2,
					Percentage: 0.3333333333333333,
				},
				{
					Value: map[string]any{
						"field1": "A",
						"field2": "A",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
			},
		},
		"top 0 returns an empty table": {
			records:  testRecords,
			path:     "nestedCase",
			top:      0,
			expected: []*record.FrequencyDistributionEntry{},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := record.FrequencyDistribution(c.records, c.path, c.caseSensitive, c.top, c.sortASC)
			require.NoError(t, err)
			assert.EqualExportedValues(t, c.expected, actual)
		})
	}
}
