package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestCount(t *testing.T) {
	cases := map[string]struct {
		records  []*api.Record
		expected int
	}{
		"empty list": {
			records:  []*api.Record{},
			expected: 0,
		},
		"filled list": {
			records: []*api.Record{
				{},
				{},
				{},
			},
			expected: 3,
		},
		"nil list": {
			records:  nil,
			expected: 0,
		},
		"list with nil values": {
			records: []*api.Record{
				{},
				nil,
			},
			expected: 2,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.expected, record.Count(c.records))
		})
	}
}
