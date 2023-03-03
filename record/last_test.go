package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestLast(t *testing.T) {
	lastRecord := &api.Record{}
	cases := map[string]struct {
		records  []*api.Record
		expected *api.Record
	}{
		"empty list": {
			records:  []*api.Record{},
			expected: nil,
		},
		"filled list": {
			records: []*api.Record{
				{},
				{},
				lastRecord,
			},
			expected: lastRecord,
		},
		"nil list": {
			records:  nil,
			expected: nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := record.Last(c.records)
			if c.expected == nil {
				assert.Nil(t, actual)
			} else {
				assert.Same(t, actual, c.expected)
			}
		})
	}
}
