package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	insights "github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestFilter(t *testing.T) {
	r1 := &api.Record{
		ID: "r1",
		Data: map[string]any{
			"value": "string A",
			"map": map[string]any{
				"foo": "bar",
				"faz": "BAZ",
			},
		},
	}
	r2 := &api.Record{
		ID: "r2",
		Data: map[string]any{
			"value": "string B",
		},
	}
	r3 := &api.Record{
		ID: "r3",
		Data: map[string]any{
			"value": "other string A",
		},
	}
	defaultRecords := []*api.Record{
		r1, r2, r3,
	}

	cases := map[string]struct {
		records     []*api.Record
		conditions  []*insights.FilterCondition
		expected    []*api.Record
		expectError bool
	}{
		"no conditions": {
			records:    defaultRecords,
			conditions: []*insights.FilterCondition{},
			expected:   defaultRecords,
		},
		"no criteria": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path: "value",
				},
			},
			expected: defaultRecords,
		},
		"equal strings": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "value",
					Equal: "string a",
				},
			},
			expected: []*api.Record{r1},
		},
		"equal strings, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					Equal:         "string A",
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{r1},
		},
		"not equal strings, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					Equal:         "string a",
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{},
		},
		"equal strings, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "does not exist",
					Equal: "string a",
				},
			},
			expected: []*api.Record{},
		},
		"equal map": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path: "map",
					Equal: map[string]any{
						"faz": "baz",
						"foo": "bar",
					},
				},
			},
			expected: []*api.Record{r1},
		},
		"not equal map, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path: "map",
					Equal: pointer(map[string]any{
						"faz": "baz",
						"foo": "bar",
					}),
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{},
		},
		"all is null": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "does not exist",
					IsNull: pointer(true),
				},
			},
			expected: defaultRecords,
		},
		"some is null": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "map",
					IsNull: pointer(true),
				},
			},
			expected: []*api.Record{r2, r3},
		},
		"none is not null": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "does not exist",
					IsNull: pointer(false),
				},
			},
			expected: []*api.Record{},
		},
		"some is not null": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "map",
					IsNull: pointer(false),
				},
			},
			expected: []*api.Record{r1},
		},
		"string starts with": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:       "value",
					StartsWith: pointer("string"),
				},
			},
			expected: []*api.Record{r1, r2},
		},
		"string starts with, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:       "does not exist",
					StartsWith: pointer("string"),
				},
			},
			expected: []*api.Record{},
		},
		"string starts with, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					StartsWith:    pointer("string"),
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{r1, r2},
		},
		"not string starts with, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					StartsWith:    pointer("STRING"),
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{},
		},
		"string ends with": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:     "value",
					EndsWith: pointer("string a"),
				},
			},
			expected: []*api.Record{r1, r3},
		},
		"string ends with, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:     "does not exist",
					EndsWith: pointer("string a"),
				},
			},
			expected: []*api.Record{},
		},
		"string ends with, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					EndsWith:      pointer("string A"),
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{r1, r3},
		},
		"not string ends with, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					EndsWith:      pointer("string a"),
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{},
		},
		"string like regex": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:      "value",
					LikeRegex: pointer(`String\s+a$`),
				},
			},
			expected: []*api.Record{r1, r3},
		},
		"string like regex, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:      "does not exist",
					LikeRegex: pointer(`string\s+a$`),
				},
			},
			expected: []*api.Record{},
		},
		"string like regex, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					LikeRegex:     pointer(`string\s+A$`),
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{r1, r3},
		},
		"not string like regex, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					LikeRegex:     pointer(`string\s+a$`),
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{},
		},
		"string like regex, invalid regex": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:      "value",
					LikeRegex: pointer(`foo(`),
				},
			},
			expectError: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := insights.Filter(c.records, c.conditions)

			if c.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}
