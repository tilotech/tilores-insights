package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestCountDistinct(t *testing.T) {
	testRecords := []*api.Record{
		{
			"someid",
			map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "a",
					"field2": "b",
				},
			},
		},
		{
			"someid",
			map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "c",
					"field2": "d",
				},
			},
		},
		{
			"someid",
			map[string]interface{}{
				"nested": map[string]interface{}{
					"field2": "b",
				},
			},
		},
		{
			"someid",
			map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "a",
				},
			},
		},
		{
			"someid",
			map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "A",
					"field2": "b",
				},
			},
		}, {
			"someid",
			map[string]interface{}{
				"nested": map[string]interface{}{
					"field1": "a",
					"field2": "b",
					"field3": "e",
				},
			},
		},
		{
			"someid",
			map[string]interface{}{},
		},
		nil,
	}
	cases := map[string]struct {
		records       []*api.Record
		paths         []string
		caseSensitive bool
		expected      int
	}{
		"empty list": {
			records:  []*api.Record{},
			paths:    []string{"nested.field1"},
			expected: 0,
		},
		"nil list": {
			records:  nil,
			paths:    []string{"nested.field1"},
			expected: 0,
		},
		"list with different values on single field": {
			records:  testRecords,
			paths:    []string{"nested.field1"},
			expected: 2,
		},
		"list with different values on multiple fields": {
			records:  testRecords,
			paths:    []string{"nested.field1", "nested.field2"},
			expected: 4,
		},
		"list with different values on object": {
			records:  testRecords,
			paths:    []string{"nested"},
			expected: 5,
		},
		"list with different values case sensitive": {
			records:       testRecords,
			paths:         []string{"nested.field1", "nested.field2"},
			caseSensitive: true,
			expected:      5,
		},
		"list with different values on object case sensitive": {
			records:       testRecords,
			paths:         []string{"nested"},
			caseSensitive: true,
			expected:      6,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, _ := record.CountDistinct(c.records, c.paths, c.caseSensitive)
			assert.Equal(t, c.expected, actual)
		})
	}
}
