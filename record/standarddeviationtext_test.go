package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestStandardDeviationText(t *testing.T) {
	testRecords := []*api.Record{
		{
			ID: "someid",
			Data: map[string]any{
				"text": "3 times",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "3 times",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "3 times",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "1 time Case sensitive 2 times",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "5 times",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "5 times",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "5 times",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "5 times",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "5 times",
			},
		},
		{
			ID: "someid",
			Data: map[string]any{
				"text": "1 time case sensitive 2 times",
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
			expected: pointer(1.247219128924647),
		},
		"list with different values case sensitive": {
			records:       testRecords,
			caseSensitive: true,
			expected:      pointer(1.6583123951777),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, _ := record.StandardDeviationText(c.records, "text", c.caseSensitive)
			if c.expected == nil {
				assert.Nil(t, actual)
			} else {
				assert.NotNil(t, actual)
				assert.Equal(t, *c.expected, *actual)
			}
		})
	}
}
