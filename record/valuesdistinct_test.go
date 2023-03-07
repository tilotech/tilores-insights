package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestValuesDistinct(t *testing.T) {
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
		expected      []any
	}{
		"empty list": {
			records:  []*api.Record{},
			path:     "nested.field1",
			expected: []any{},
		},
		"nil list": {
			records:  nil,
			path:     "nested.field1",
			expected: []any{},
		},
		"list with different values on field": {
			records:  testRecords,
			path:     "nested.field1",
			expected: []any{"a", "c"},
		},
		"list with different values on object": {
			records: testRecords,
			path:    "nested",
			expected: []any{
				map[string]interface{}{
					"field1": "a",
					"field2": "b",
				},
				map[string]interface{}{
					"field1": "c",
					"field2": "d",
				},
				map[string]interface{}{
					"field2": "b",
				},
				map[string]interface{}{
					"field1": "a",
				},
				map[string]interface{}{
					"field1": "a",
					"field2": "b",
					"field3": "e",
				},
			},
		},
		"list with different values case sensitive": {
			records:       testRecords,
			path:          "nested.field1",
			caseSensitive: true,
			expected:      []any{"a", "c", "A"},
		},
		"list with different values on object case sensitive": {
			records:       testRecords,
			path:          "nested",
			caseSensitive: true,
			expected: []any{
				map[string]interface{}{
					"field1": "a",
					"field2": "b",
				},
				map[string]interface{}{
					"field1": "c",
					"field2": "d",
				},
				map[string]interface{}{
					"field2": "b",
				},
				map[string]interface{}{
					"field1": "a",
				},
				map[string]interface{}{
					"field1": "A",
					"field2": "b",
				},
				map[string]interface{}{
					"field1": "a",
					"field2": "b",
					"field3": "e",
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, _ := record.ValuesDistinct(c.records, c.path, c.caseSensitive)
			assert.ElementsMatch(t, c.expected, actual)
		})
	}
}
