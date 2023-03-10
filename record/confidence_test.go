package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestConfidence(t *testing.T) {
	testRecords := []*api.Record{
		{
			ID: "someid",
			Data: map[string]any{
				"text": "a",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "a",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "a",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "b",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "c",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "c",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "c",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "c",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "c",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "B",
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
		caseSensitive bool
		expected      *float64
		expectError   bool
	}{
		"empty list": {
			records:  []*api.Record{},
			expected: nil,
		},
		"nil list": {
			records:  nil,
			expected: nil,
		},
		"list with all nil values": {
			records: []*api.Record{
				nil,
				nil,
			},
			expected: nil,
		},
		"list with different values": {
			records:  testRecords,
			expected: pointer(0.38),
		},
		"list with different values case sensitive": {
			records:       testRecords,
			caseSensitive: true,
			expected:      pointer(0.36),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := record.Confidence(c.records, "text", c.caseSensitive)
			require.NoError(t, err)
			if c.expected == nil {
				assert.Nil(t, actual)
			} else {
				require.NotNil(t, actual)
				assert.Equal(t, *c.expected, *actual)
			}
		})
	}
}
