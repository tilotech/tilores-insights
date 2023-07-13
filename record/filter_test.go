package record_test

import (
	"testing"
	"time"

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
			"numeric": 12.3,
			"time":    "2022-01-01T00:00:00Z",
		},
	}
	r2 := &api.Record{
		ID: "r2",
		Data: map[string]any{
			"value":   "string B",
			"numeric": "123",
			"time":    "2023-01-01T00:00:00Z",
		},
	}
	r3 := &api.Record{
		ID: "r3",
		Data: map[string]any{
			"value":   "other string A",
			"numeric": 1234.0,
			"time":    "2024-01-01T00:00:00Z",
		},
	}
	r4 := &api.Record{
		ID: "r4",
		Data: map[string]any{
			"numeric": 1.234567e+06,
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
					Path:   "value",
					Equals: "string a",
				},
			},
			expected: []*api.Record{r1},
		},
		"equal too big numeric": {
			records: []*api.Record{r4},
			conditions: []*insights.FilterCondition{
				{
					Path:   "numeric",
					Equals: 1234567,
				},
			},
			expected: []*api.Record{r4},
		},
		"equal numeric as string": {
			records: []*api.Record{r4},
			conditions: []*insights.FilterCondition{
				{
					Path:   "numeric",
					Equals: "1.234567e+06",
				},
			},
			expected: []*api.Record{r4},
		},
		"equal strings, case sensitive": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					Equals:        "string A",
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
					Equals:        "string a",
					CaseSensitive: pointer(true),
				},
			},
			expected: []*api.Record{},
		},
		"equal strings, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "does not exist",
					Equals: "string a",
				},
			},
			expected: []*api.Record{},
		},
		"equal map": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path: "map",
					Equals: map[string]any{
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
					Equals: pointer(map[string]any{
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
		"less than": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:     "numeric",
					LessThan: pointer(150.0),
				},
			},
			expected: []*api.Record{r1, r2},
		},
		"less than, equal second record": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:     "numeric",
					LessThan: pointer(123.0),
				},
			},
			expected: []*api.Record{r1},
		},
		"less than, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:     "does not exist",
					LessThan: pointer(150.0),
				},
			},
			expected: []*api.Record{},
		},
		"less than, not a number": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:     "value",
					LessThan: pointer(150.0),
				},
			},
			expectError: true,
		},
		"less equal": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:       "numeric",
					LessEquals: pointer(150.0),
				},
			},
			expected: []*api.Record{r1, r2},
		},
		"less equal, equal second record": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:       "numeric",
					LessEquals: pointer(123.0),
				},
			},
			expected: []*api.Record{r1, r2},
		},
		"less equal, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:       "does not exist",
					LessEquals: pointer(150.0),
				},
			},
			expected: []*api.Record{},
		},
		"less equal, not a number": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:       "value",
					LessEquals: pointer(150.0),
				},
			},
			expectError: true,
		},
		"greater than": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:        "numeric",
					GreaterThan: pointer(100.0),
				},
			},
			expected: []*api.Record{r2, r3},
		},
		"greater than, equal second record": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:        "numeric",
					GreaterThan: pointer(123.0),
				},
			},
			expected: []*api.Record{r3},
		},
		"greater than, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:        "does not exist",
					GreaterThan: pointer(100.0),
				},
			},
			expected: []*api.Record{},
		},
		"greater than, not a number": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:        "value",
					GreaterThan: pointer(100.0),
				},
			},
			expectError: true,
		},
		"greater equal": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "numeric",
					GreaterEquals: pointer(100.0),
				},
			},
			expected: []*api.Record{r2, r3},
		},
		"greater equal, equal second record": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "numeric",
					GreaterEquals: pointer(123.0),
				},
			},
			expected: []*api.Record{r2, r3},
		},
		"greater equal, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "does not exist",
					GreaterEquals: pointer(100.0),
				},
			},
			expected: []*api.Record{},
		},
		"greater equal, not a number": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:          "value",
					GreaterEquals: pointer(100.0),
				},
			},
			expectError: true,
		},
		"after": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "time",
					After: pointer(time.Date(2022, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{r2, r3},
		},
		"after, equal second record": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "time",
					After: pointer(time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{r3},
		},
		"after, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "does not exist",
					After: pointer(time.Date(2022, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{},
		},
		"after, not a time": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "value",
					After: pointer(time.Date(2022, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expectError: true,
		},
		"since": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "time",
					Since: pointer(time.Date(2022, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{r2, r3},
		},
		"since, equal second record": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "time",
					Since: pointer(time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{r2, r3},
		},
		"since, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "does not exist",
					Since: pointer(time.Date(2022, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{},
		},
		"since, not a time": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "value",
					Since: pointer(time.Date(2022, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expectError: true,
		},
		"before": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "time",
					Before: pointer(time.Date(2023, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{r1, r2},
		},
		"before, equal second record": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "time",
					Before: pointer(time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{r1},
		},
		"before, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "does not exist",
					Before: pointer(time.Date(2023, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{},
		},
		"before, not a time": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "value",
					Before: pointer(time.Date(2023, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expectError: true,
		},
		"until": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "time",
					Until: pointer(time.Date(2023, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{r1, r2},
		},
		"until, equal second record": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "time",
					Until: pointer(time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{r1, r2},
		},
		"until, nil value": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "does not exist",
					Until: pointer(time.Date(2023, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expected: []*api.Record{},
		},
		"until, not a time": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:  "value",
					Until: pointer(time.Date(2023, 06, 01, 0, 0, 0, 0, time.UTC)),
				},
			},
			expectError: true,
		},
		"equal strings, inverted": {
			records: defaultRecords,
			conditions: []*insights.FilterCondition{
				{
					Path:   "value",
					Equals: "string a",
					Invert: pointer(true),
				},
			},
			expected: []*api.Record{r2, r3},
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
