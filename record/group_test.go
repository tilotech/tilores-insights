package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	insights "github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestGroup(t *testing.T) {
	r1 := &api.Record{
		ID: "r1",
		Data: map[string]any{
			"value":  "a",
			"value2": true,
		},
	}
	r2 := &api.Record{
		ID: "r2",
		Data: map[string]any{
			"value":  nil,
			"value2": true,
		},
	}
	r3 := &api.Record{
		ID: "r3",
		Data: map[string]any{
			"value":  "b",
			"value2": true,
		},
	}
	r4 := &api.Record{
		ID: "r4",
		Data: map[string]any{
			"value":  "A",
			"value2": true,
		},
	}
	r5 := &api.Record{
		ID: "r5",
		Data: map[string]any{
			"value":  1,
			"value2": true,
		},
	}
	r6 := &api.Record{
		ID: "r6",
		Data: map[string]any{
			"value":  []any{"a", "b"},
			"value2": true,
		},
	}
	r7 := &api.Record{
		ID: "r7",
		Data: map[string]any{
			"value":  []any{"A", "b"},
			"value2": true,
		},
	}
	r8 := &api.Record{
		ID: "r8",
		Data: map[string]any{
			"value":  "a",
			"value2": false,
		},
	}
	defaultRecords := []*api.Record{
		r1, r2, r3, r4, r5, r6, r7, r8,
	}

	cases := map[string]struct {
		records       []*api.Record
		paths         []string
		caseSensitive bool
		expected      [][]*api.Record
	}{
		"empty records": {
			records:  []*api.Record{},
			paths:    []string{"value"},
			expected: [][]*api.Record{},
		},
		"nil records": {
			records:  nil,
			paths:    []string{"value"},
			expected: [][]*api.Record{},
		},
		"empty paths": {
			records:  defaultRecords,
			paths:    []string{},
			expected: [][]*api.Record{defaultRecords},
		},
		"group on value": {
			records: defaultRecords,
			paths:   []string{"value"},
			expected: [][]*api.Record{
				{r1, r4, r8}, // value: a
				{r2},         // value: nil
				{r3},         // value: b
				{r5},         // value: 1
				{r6, r7},     // value [a, b]
			},
		},
		"group on value, case sensitive": {
			records:       defaultRecords,
			paths:         []string{"value"},
			caseSensitive: true,
			expected: [][]*api.Record{
				{r1, r8}, // value: a
				{r2},     // value: nil
				{r3},     // value: b
				{r4},     // value: A
				{r5},     // value: 1
				{r6},     // value: [A, b]
				{r7},     // value: [a, b]
			},
		},
		"group on value and value 2": {
			records: defaultRecords,
			paths:   []string{"value", "value2"},
			expected: [][]*api.Record{
				{r1, r4}, // value: a/true
				{r2},     // value: nil/true
				{r3},     // value: b/true
				{r5},     // value: 1/true
				{r6, r7}, // value [a, b]/true
				{r8},     // value: a/false
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := insights.Group(c.records, c.paths, c.caseSensitive)
			require.NoError(t, err)
			assert.ElementsMatch(t, c.expected, actual)
		})
	}
}
