package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	insights "github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestOldest(t *testing.T) {
	r1 := &api.Record{
		ID: "r1",
		Data: map[string]any{
			"time":         "2022-01-01T00:00:00Z",
			"invalid-time": "2022-01-01T00:00:00Z",
			"nil-records":  nil,
		},
	}
	r2 := &api.Record{
		ID: "r2",
		Data: map[string]any{
			"time":         "2023-01-01T00:00:00Z",
			"invalid-time": "not a valid time",
		},
	}
	r3 := &api.Record{
		ID: "r3",
		Data: map[string]any{
			"time": nil,
		},
	}

	defaultRecords := []*api.Record{
		r1, r2, r3,
	}

	cases := map[string]struct {
		records     []*api.Record
		expected    *api.Record
		expectError bool
	}{
		"empty records": {
			records:  []*api.Record{},
			expected: nil,
		},
		"nil-records": {
			records:  nil,
			expected: nil,
		},
		"time": {
			records:  defaultRecords,
			expected: r1,
		},
		"invalid-time": {
			records:     defaultRecords,
			expectError: true,
		},
	}

	for path, c := range cases {
		t.Run(path, func(t *testing.T) {
			actual, err := insights.Oldest(c.records, path)
			if c.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Same(t, c.expected, actual)
			}
		})
	}
}
