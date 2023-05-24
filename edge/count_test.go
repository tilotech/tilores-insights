package edge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/edge"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestCount(t *testing.T) {
	cases := map[string]struct {
		edges    api.Edges
		expected int
	}{
		"empty list": {
			edges:    api.Edges{},
			expected: 0,
		},
		"filled list": {
			edges: api.Edges{
				"record-1:record-2:R1",
				"record-1:record-2:R2",
				"record-1:record-3:R1",
			},
			expected: 3,
		},
		"nil list": {
			edges:    nil,
			expected: 0,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.expected, edge.Count(c.edges))
		})
	}
}
