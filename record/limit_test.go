package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	insights "github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestLimit(t *testing.T) {
	r1 := &api.Record{
		ID:   "r1",
		Data: map[string]any{},
	}
	r2 := &api.Record{
		ID:   "r2",
		Data: map[string]any{},
	}
	r3 := &api.Record{
		ID:   "r3",
		Data: map[string]any{},
	}

	defaultRecords := []*api.Record{
		r1, r2, r3,
	}

	cases := map[string]struct {
		records  []*api.Record
		count    int
		offset   int
		expected []*api.Record
	}{
		"empty list": {
			records:  []*api.Record{},
			count:    1,
			offset:   0,
			expected: []*api.Record{},
		},
		"nil list": {
			records:  nil,
			count:    1,
			offset:   0,
			expected: []*api.Record{},
		},
		"limit 0": {
			records:  defaultRecords,
			count:    0,
			offset:   0,
			expected: []*api.Record{},
		},
		"limit 2": {
			records:  defaultRecords,
			count:    2,
			offset:   0,
			expected: []*api.Record{r1, r2},
		},
		"limit 3": {
			records:  defaultRecords,
			count:    3,
			offset:   0,
			expected: []*api.Record{r1, r2, r3},
		},
		"limit 4": {
			records:  defaultRecords,
			count:    4,
			offset:   0,
			expected: []*api.Record{r1, r2, r3},
		},
		"limit 1,1": {
			records:  defaultRecords,
			count:    1,
			offset:   1,
			expected: []*api.Record{r2},
		},
		"limit 4,1": {
			records:  defaultRecords,
			count:    4,
			offset:   1,
			expected: []*api.Record{r2, r3},
		},
		"limit 1,3": {
			records:  defaultRecords,
			count:    1,
			offset:   3,
			expected: []*api.Record{},
		},
		"limit -1": {
			records:  defaultRecords,
			count:    -1,
			offset:   0,
			expected: []*api.Record{},
		},
		"limit 1,-1": {
			records:  defaultRecords,
			count:    1,
			offset:   -1,
			expected: []*api.Record{},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := insights.Limit(c.records, c.count, c.offset)
			assert.Equal(t, c.expected, actual)
		})
	}
}
