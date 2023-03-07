package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestFrequencyDistribution(t *testing.T) {
	testRecords := []*api.Record{
		{
			ID: "someid",
			Data: map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "a",
					"field2": "b",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "c",
					"field2": "d",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]interface{}{
				"nested": map[string]interface{}{
					"field2": "b",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "a",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "A",
					"field2": "b",
				},
			},
		}, {
			ID: "someid",
			Data: map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "a",
					"field2": "b",
					"field3": "e",
				},
			},
		},
		{
			ID:   "someid",
			Data: map[string]interface{}{},
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
		expectedFirst *record.FrequencyDistributionEntry
		expectedLast  *record.FrequencyDistributionEntry
	}{
		"empty list": {
			records:  []*api.Record{},
			path:     "nested.field1",
			expected: []*record.FrequencyDistributionEntry{},
		},
		"nil list": {
			records:  nil,
			path:     "nested.field1",
			expected: []*record.FrequencyDistributionEntry{},
		},
		"on field": {
			records: testRecords,
			path:    "nested.field1",
			expectedFirst: &record.FrequencyDistributionEntry{
				Value:      "a",
				Frequency:  4,
				Percentage: 0.8,
			},
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "a",
					Frequency:  4,
					Percentage: 0.8,
				},
				{
					Value:      "c",
					Frequency:  1,
					Percentage: 0.2,
				},
			},
		},
		"on field top 1": {
			records: testRecords,
			path:    "nested.field1",
			top:     1,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "a",
					Frequency:  4,
					Percentage: 0.8,
				},
			},
		},
		"on field asc": {
			records: testRecords,
			path:    "nested.field1",
			sortASC: true,
			expectedLast: &record.FrequencyDistributionEntry{
				Value:      "a",
				Frequency:  4,
				Percentage: 0.8,
			},
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "a",
					Frequency:  4,
					Percentage: 0.8,
				},
				{
					Value:      "c",
					Frequency:  1,
					Percentage: 0.2,
				},
			},
		},
		"on object": {
			records: testRecords,
			path:    "nested",
			expectedFirst: &record.FrequencyDistributionEntry{
				Value: map[string]interface{}{
					"field1": "a",
					"field2": "b",
				},
				Frequency:  2,
				Percentage: 0.3333333333333333,
			},
			expected: []*record.FrequencyDistributionEntry{
				{
					Value: map[string]interface{}{
						"field1": "a",
						"field2": "b",
					},
					Frequency:  2,
					Percentage: 0.3333333333333333,
				},
				{
					Value: map[string]interface{}{
						"field1": "c",
						"field2": "d",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field2": "b",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field1": "a",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field1": "a",
						"field2": "b",
						"field3": "e",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
			},
		},
		"on object top 3": {
			records: testRecords,
			path:    "nested",
			top:     3,
			expectedFirst: &record.FrequencyDistributionEntry{
				Value: map[string]interface{}{
					"field1": "a",
					"field2": "b",
				},
				Frequency:  2,
				Percentage: 0.3333333333333333,
			},
			expected: []*record.FrequencyDistributionEntry{
				{
					Value: map[string]interface{}{
						"field1": "a",
						"field2": "b",
					},
					Frequency:  2,
					Percentage: 0.3333333333333333,
				},
				{
					Value: map[string]interface{}{
						"field1": "c",
						"field2": "d",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field2": "b",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
			},
		},
		"on object asc": {
			records: testRecords,
			path:    "nested",
			sortASC: true,
			expectedLast: &record.FrequencyDistributionEntry{
				Value: map[string]interface{}{
					"field1": "a",
					"field2": "b",
				},
				Frequency:  2,
				Percentage: 0.3333333333333333,
			},
			expected: []*record.FrequencyDistributionEntry{
				{
					Value: map[string]interface{}{
						"field1": "a",
						"field2": "b",
					},
					Frequency:  2,
					Percentage: 0.3333333333333333,
				},
				{
					Value: map[string]interface{}{
						"field1": "c",
						"field2": "d",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field2": "b",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field1": "a",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field1": "a",
						"field2": "b",
						"field3": "e",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
			},
		},
		"case sensitive": {
			records:       testRecords,
			path:          "nested.field1",
			caseSensitive: true,
			expectedFirst: &record.FrequencyDistributionEntry{
				Value:      "a",
				Frequency:  3,
				Percentage: 0.6,
			},
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "a",
					Frequency:  3,
					Percentage: 0.6,
				},
				{
					Value:      "c",
					Frequency:  1,
					Percentage: 0.2,
				},
				{
					Value:      "A",
					Frequency:  1,
					Percentage: 0.2,
				},
			},
		},
		"case sensitive asc": {
			records:       testRecords,
			path:          "nested.field1",
			caseSensitive: true,
			sortASC:       true,
			expectedLast: &record.FrequencyDistributionEntry{
				Value:      "a",
				Frequency:  3,
				Percentage: 0.6,
			},
			expected: []*record.FrequencyDistributionEntry{
				{
					Value:      "a",
					Frequency:  3,
					Percentage: 0.6,
				},
				{
					Value:      "c",
					Frequency:  1,
					Percentage: 0.2,
				},
				{
					Value:      "A",
					Frequency:  1,
					Percentage: 0.2,
				},
			},
		},
		"on object case sensitive": {
			records:       testRecords,
			path:          "nested",
			caseSensitive: true,
			expected: []*record.FrequencyDistributionEntry{
				{
					Value: map[string]interface{}{
						"field1": "a",
						"field2": "b",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field1": "c",
						"field2": "d",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field2": "b",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field1": "a",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field1": "A",
						"field2": "b",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
				{
					Value: map[string]interface{}{
						"field1": "a",
						"field2": "b",
						"field3": "e",
					},
					Frequency:  1,
					Percentage: 0.16666666666666666,
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, _ := record.FrequencyDistribution(c.records, c.path, c.caseSensitive, c.top, c.sortASC)
			if c.top < 1 {
				assert.ElementsMatch(t, c.expected, actual)
			}
			if c.top > 1 {
				assert.Len(t, actual, c.top)
			}
			if c.top == 1 {
				assert.Equal(t, c.expected, actual)
			}
			if c.expectedFirst != nil {
				assert.NotEmpty(t, actual)
				assert.Equal(t, c.expectedFirst, actual[0])
			}
			if c.expectedLast != nil {
				assert.NotEmpty(t, actual)
				assert.Equal(t, c.expectedLast, actual[len(actual)-1])
			}
		})
	}
}
