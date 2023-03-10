package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestFlattenDistinct(t *testing.T) {
	testRecords := []*api.Record{
		{
			ID: "someid",
			Data: map[string]any{
				"arr": []any{
					"v1",
					"v2",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"arr": []any{
					"v2",
					"v3",
				},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"arr": []any{},
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"arr": nil,
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"arr": []any{
					4,
					nil,
					"",
					"",
					map[string]any{
						"nested": true,
					},
					"v5",
					"V5",
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
		path          string
		caseSensitive bool
		expected      []any
		expectError   bool
	}{
		"empty list": {
			records:  []*api.Record{},
			path:     "arr",
			expected: []any{},
		},
		"nil list": {
			records:  nil,
			path:     "arr",
			expected: []any{},
		},
		"list with different values": {
			records:  testRecords,
			path:     "arr",
			expected: []any{"v1", "v2", "v3", 4, "", map[string]any{"nested": true}, "v5"},
		},
		"list with different values case sensitive": {
			records:       testRecords,
			path:          "arr",
			caseSensitive: true,
			expected:      []any{"v1", "v2", "v3", 4, "", map[string]any{"nested": true}, "v5", "V5"},
		},
		"non existent field": {
			records:  testRecords,
			path:     "none",
			expected: []any{},
		},
		"fails on none array / nil": {
			records: []*api.Record{
				{
					ID: "someid",
					Data: map[string]any{
						"arr": []string{
							"v1",
							"v2",
						},
					},
				},
				{
					ID: "someid",
					Data: map[string]any{
						"arr": "v2",
					},
				},
			},
			path:        "arr",
			expectError: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := record.FlattenDistinct(c.records, c.path, c.caseSensitive)
			if c.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.expected, actual)
			}
		})
	}
}
