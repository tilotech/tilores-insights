package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestCountDistinct(t *testing.T) {
	testRecords := []*api.Record{
		{
			ID: "someid",
			Data: map[string]any{
				"nested": map[string]any{
					"field1": "a",
					"field2": "b",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"nested": map[string]any{
					"field1": "c",
					"field2": "d",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"nested": map[string]any{
					"field2": "b",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"nested": map[string]any{
					"field1": "a",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"nested": map[string]any{
					"field1": "A",
					"field2": "b",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"nested": map[string]any{
					"field1": "a",
					"field2": "b",
					"field3": "e",
				},
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
			actual, err := record.CountDistinct(c.records, c.paths, c.caseSensitive)
			require.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}
