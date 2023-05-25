package edge_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/edge"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestMatrix(t *testing.T) {
	cases := map[string]struct {
		edges      api.Edges
		duplicates api.Duplicates
		links      []string
		expected   []edge.MatrixEntry
	}{
		"no edges, no duplicates": {
			edges:      api.Edges{},
			duplicates: api.Duplicates{},
			expected:   []edge.MatrixEntry{},
		},
		"single edge, link list provided": {
			edges: api.Edges{
				"record-1:record-2:R1",
			},
			duplicates: api.Duplicates{},
			links:      []string{"R1", "R2"},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"R1": true,
						"R2": false,
					},
				},
			},
		},
		"single edge, link list not provided": {
			edges: api.Edges{
				"record-1:record-2:R1",
			},
			duplicates: api.Duplicates{},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"R1": true,
					},
				},
			},
		},
		"multiple edges, link list provided": {
			edges: api.Edges{
				"record-1:record-2:R1",
				"record-2:record-1:R2",
				"record-1:record-2:R3",
				"record-2:record-3:R1",
			},
			duplicates: api.Duplicates{},
			links:      []string{"R1", "R2"},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"R1": true,
						"R2": true,
					},
				},
				{
					A: "record-2",
					B: "record-3",
					Links: map[string]bool{
						"R1": true,
						"R2": false,
					},
				},
			},
		},
		"multiple edges, link list not provided": {
			edges: api.Edges{
				"record-1:record-2:R1",
				"record-2:record-1:R2",
				"record-1:record-2:R3",
				"record-2:record-3:R1",
			},
			duplicates: api.Duplicates{},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"R1": true,
						"R2": true,
						"R3": true,
					},
				},
				{
					A: "record-2",
					B: "record-3",
					Links: map[string]bool{
						"R1": true,
						"R2": false,
						"R3": false,
					},
				},
			},
		},
		"single duplicate, link list provided": {
			edges: api.Edges{},
			duplicates: api.Duplicates{
				"record-1": []string{
					"record-2",
				},
			},
			links: []string{
				"R1",
				"duplicate",
			},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"R1":        false,
						"duplicate": true,
					},
				},
			},
		},
		"single duplicate, link list not provided": {
			edges: api.Edges{},
			duplicates: api.Duplicates{
				"record-1": []string{
					"record-2",
				},
			},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"duplicate": true,
					},
				},
			},
		},
		"multiple duplicates, link list provided": {
			edges: api.Edges{},
			duplicates: api.Duplicates{
				"record-1": []string{
					"record-2",
					"record-3",
				},
				"record-5": []string{
					"record-4",
				},
			},
			links: []string{
				"R1",
				"duplicate",
			},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"R1":        false,
						"duplicate": true,
					},
				},
				{
					A: "record-1",
					B: "record-3",
					Links: map[string]bool{
						"R1":        false,
						"duplicate": true,
					},
				},
				{
					A: "record-4",
					B: "record-5",
					Links: map[string]bool{
						"R1":        false,
						"duplicate": true,
					},
				},
			},
		},
		"duplicates with rule groups, link list provided": {
			edges: api.Edges{},
			duplicates: api.Duplicates{
				"G1:record-1": []string{
					"record-2",
					"record-3",
				},
				"G2:record-1": []string{
					"record-3",
				},
				"G2:record-2": []string{
					"record-4",
				},
				"G3:record-1": []string{
					"record-3",
					"record-4",
				},
			},
			links: []string{
				"R1",
				"G1:duplicate",
				"G2:duplicate",
			},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"R1":           false,
						"G1:duplicate": true,
						"G2:duplicate": false,
					},
				},
				{
					A: "record-1",
					B: "record-3",
					Links: map[string]bool{
						"R1":           false,
						"G1:duplicate": true,
						"G2:duplicate": true,
					},
				},
				{
					A: "record-2",
					B: "record-4",
					Links: map[string]bool{
						"R1":           false,
						"G1:duplicate": false,
						"G2:duplicate": true,
					},
				},
			},
		},
		"duplicates with rule groups, link list not provided": {
			edges: api.Edges{},
			duplicates: api.Duplicates{
				"G1:record-1": []string{
					"record-2",
					"record-3",
				},
				"G2:record-1": []string{
					"record-3",
				},
				"G2:record-2": []string{
					"record-4",
				},
				"G3:record-1": []string{
					"record-3",
					"record-4",
				},
			},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"G1:duplicate": true,
						"G2:duplicate": false,
						"G3:duplicate": false,
					},
				},
				{
					A: "record-1",
					B: "record-3",
					Links: map[string]bool{
						"G1:duplicate": true,
						"G2:duplicate": true,
						"G3:duplicate": true,
					},
				},
				{
					A: "record-1",
					B: "record-4",
					Links: map[string]bool{
						"G1:duplicate": false,
						"G2:duplicate": false,
						"G3:duplicate": true,
					},
				},
				{
					A: "record-2",
					B: "record-4",
					Links: map[string]bool{
						"G1:duplicate": false,
						"G2:duplicate": true,
						"G3:duplicate": false,
					},
				},
			},
		},
		"mixed edges and duplicates, link list not provided": {
			edges: api.Edges{
				"record-1:record-2:R1",
				"record-2:record-3:R2",
				"record-2:record-3:R3",
			},
			duplicates: api.Duplicates{
				"record-1": []string{
					"record-4",
				},
			},
			expected: []edge.MatrixEntry{
				{
					A: "record-1",
					B: "record-2",
					Links: map[string]bool{
						"R1":        true,
						"R2":        false,
						"R3":        false,
						"duplicate": false,
					},
				},
				{
					A: "record-1",
					B: "record-4",
					Links: map[string]bool{
						"R1":        false,
						"R2":        false,
						"R3":        false,
						"duplicate": true,
					},
				},
				{
					A: "record-2",
					B: "record-3",
					Links: map[string]bool{
						"R1":        false,
						"R2":        true,
						"R3":        true,
						"duplicate": false,
					},
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := edge.Matrix(c.edges, c.duplicates, c.links)
			assert.ElementsMatch(t, c.expected, actual)
		})
	}
}
