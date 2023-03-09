package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	insights "github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestSort(t *testing.T) {
	r1 := &api.Record{
		ID: "r1",
		Data: map[string]any{
			"string":        "a",
			"number":        10.0,
			"similarstring": "b",
			"mixed":         2.0,
		},
	}
	r2 := &api.Record{
		ID: "r2",
		Data: map[string]any{
			"string":        "c",
			"number":        5.0,
			"similarstring": "a",
			"mixed":         "01",
		},
	}
	r3 := &api.Record{
		ID: "r3",
		Data: map[string]any{
			"string":        nil,
			"number":        5.0,
			"similarstring": nil,
			"mixed":         10.0,
		},
	}
	r4 := &api.Record{
		ID: "r4",
		Data: map[string]any{
			"string":        "b",
			"number":        "1",
			"similarstring": "a",
			"mixed":         "not a number",
		},
	}
	defaultRecords := []*api.Record{
		r1, r2, r3, r4,
	}

	cases := map[string]struct {
		records  []*api.Record
		criteria []*insights.SortCriteria
		expected []*api.Record
	}{
		"no criteria": {
			records:  defaultRecords,
			criteria: []*insights.SortCriteria{},
			expected: defaultRecords,
		},
		"sort by string ASC": {
			records: defaultRecords,
			criteria: []*insights.SortCriteria{
				{
					Path: "string",
					ASC:  true,
				},
			},
			expected: []*api.Record{r1, r4, r2, r3},
		},
		"sort by string DESC": {
			records: defaultRecords,
			criteria: []*insights.SortCriteria{
				{
					Path: "string",
					ASC:  false,
				},
			},
			expected: []*api.Record{r3, r2, r4, r1},
		},
		"sort by number": {
			records: defaultRecords,
			criteria: []*insights.SortCriteria{
				{
					Path: "number",
					ASC:  true,
				},
			},
			expected: []*api.Record{r4, r2, r3, r1},
		},
		"sort by number DESC": {
			records: defaultRecords,
			criteria: []*insights.SortCriteria{
				{
					Path: "number",
					ASC:  false,
				},
			},
			expected: []*api.Record{r1, r2, r3, r4},
		},
		"sort by similarstring ASC, number DESC": {
			records: defaultRecords,
			criteria: []*insights.SortCriteria{
				{
					Path: "similarstring",
					ASC:  false,
				},
				{
					Path: "number",
					ASC:  true,
				},
			},
			expected: []*api.Record{r3, r1, r4, r2},
		},
		"sort by mixed": {
			records: defaultRecords,
			criteria: []*insights.SortCriteria{
				{
					Path: "mixed",
					ASC:  true,
				},
			},
			expected: []*api.Record{r2, r3, r1, r4},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := insights.Sort(c.records, c.criteria)
			require.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}
