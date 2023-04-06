package record_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestExtract(t *testing.T) {
	dataJSON := `
	{
		"value": "string",
		"nested": {
			"value": "nested string value",
			"super": {
				"value": "Super Nested String Value"
			}
		},
		"int": 123,
		"list": [
			"abc",
			"DEF",
			"geh"
		],
		"nullValue": null,
		"emptyString": ""
	}
	`
	data := map[string]any{}
	err := json.Unmarshal([]byte(dataJSON), &data)
	require.NoError(t, err)

	r := &api.Record{
		ID:   "some-id",
		Data: data,
	}

	cases := map[string]struct {
		useCustomData bool
		customData    *api.Record
		expected      any
	}{
		"value": {
			expected: "string",
		},
		"nested.value": {
			expected: "nested string value",
		},
		"nested.super.value": {
			expected: "Super Nested String Value",
		},
		"nested": {
			expected: map[string]any{
				"value": "nested string value",
				"super": map[string]any{
					"value": "Super Nested String Value",
				},
			},
		},
		"int": {
			expected: 123.0, // json parses numbers as float64!
		},
		"list": {
			expected: []any{
				"abc",
				"DEF",
				"geh",
			},
		},
		"list.0": {
			expected: "abc",
		},
		"list.1": {
			expected: "DEF",
		},
		"list.2": {
			expected: "geh",
		},
		"nil as input": {
			useCustomData: true,
			customData:    &api.Record{},
			expected:      nil,
		},
		"nil record as input": {
			useCustomData: true,
			customData:    nil,
			expected:      nil,
		},
		"nonexistent": {
			expected: nil,
		},
		"emptyString": {
			expected: "",
		},
		"nullValue": {
			expected: nil,
		},
		"non.existent": {
			expected: nil,
		},
		"nested.nonexistent": {
			expected: nil,
		},
		"nested.value.nonexistent": {
			expected: nil,
		},
		"int.nonexistent": {
			expected: nil,
		},
		"list.a": {
			expected: nil,
		},
		"list.4": {
			expected: nil,
		},
		"list.-1": {
			expected: nil,
		},
	}

	for path, c := range cases {
		t.Run(path, func(t *testing.T) {
			input := r
			if c.useCustomData {
				input = c.customData
			}
			actual := record.Extract(input, path)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestExtractNumber(t *testing.T) {
	dataJSON := `
	{
		"nonnumeric": "string",
		"int": 123,
		"float": 123.4,
		"nullValue": null,
		"numericText": "123",
		"exponent": "1e3",
		"nested": {
			"value": "123"
		}
	}
	`
	data := map[string]any{}
	err := json.Unmarshal([]byte(dataJSON), &data)
	require.NoError(t, err)

	r := &api.Record{
		ID:   "some-id",
		Data: data,
	}

	cases := map[string]struct {
		expected    any
		expectError bool
	}{
		"nonnumeric": {
			expectError: true,
		},
		"int": {
			expected: pointer(123.0),
		},
		"float": {
			expected: pointer(123.4),
		},
		"nullValue": {
			expected: nil,
		},
		"numericText": {
			expected: pointer(123.0),
		},
		"exponent": {
			expected: pointer(1000.0),
		},
		"nested": {
			expectError: true,
		},
	}

	for path, c := range cases {
		t.Run(path, func(t *testing.T) {
			actual, err := record.ExtractNumber(r, path)
			if c.expectError {
				assert.Error(t, err)
			} else if c.expected == nil {
				require.NoError(t, err)
				assert.Nil(t, actual)
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.expected, actual)
			}
		})
	}
}

func TestExtractString(t *testing.T) {
	dataJSON := `
	{
		"nested": {
			"value": "nested string value",
			"super": {
				"value": "Super Nested String Value"
			}
		},
		"list": [
			"abc",
			"DEF",
			"geh"
		],
		"keepUpper": "Has Upper Case",
		"caseInsensitive": "Has Upper Case",
    "bool": true,
		"emptyString": "",
		"int": 123,
		"float": 123.4,
		"nullValue": null,
		"numericText": "123",
		"exponent": "1e3",
		"nested": {
			"propB": "valB",
			"propA": "valA"
		}
	}
	`
	data := map[string]any{}
	err := json.Unmarshal([]byte(dataJSON), &data)
	require.NoError(t, err)

	r := &api.Record{
		ID:   "some-id",
		Data: data,
	}

	cases := map[string]struct {
		expected      *string
		caseSensitive bool
	}{
		"bool": {
			expected: pointer("true"),
		},
		"int": {
			expected: pointer("123"),
		},
		"float": {
			expected: pointer("123.4"),
		},
		"nullValue": {
			expected: nil,
		},
		"numericText": {
			expected: pointer("123"),
		},
		"exponent": {
			expected: pointer("1e3"),
		},
		"nested": {
			caseSensitive: true,
			expected:      pointer(`{"propA":"valA","propB":"valB"}`),
		},
		"keepUpper": {
			caseSensitive: true,
			expected:      pointer("Has Upper Case"),
		},
		"caseInsensitive": {
			expected: pointer("has upper case"),
		},
		"list.0": {
			expected: pointer("abc"),
		},
		"list.1": {
			expected: pointer("def"),
		},
		"list.2": {
			expected: pointer("geh"),
		},
		"list": {
			expected: pointer(`["abc","def","geh"]`),
		},
		"emptyString": {
			expected: pointer(""),
		},
	}

	for path, c := range cases {
		t.Run(path, func(t *testing.T) {
			actual, err := record.ExtractString(r, path, c.caseSensitive)
			require.NoError(t, err)
			if c.expected == nil {
				assert.Nil(t, actual)
			} else {
				require.NotNil(t, actual)
				assert.Equal(t, *c.expected, *actual)
			}
		})
	}
}

func TestExtractTime(t *testing.T) {
	dataJSON := `
	{
		"time": "2023-03-07T16:06:05Z",
		"timeISO": "2022-02-28T06:56:47.778565",
		"not-a-time": "something else"
	}
	`
	data := map[string]any{}
	err := json.Unmarshal([]byte(dataJSON), &data)
	require.NoError(t, err)

	r := &api.Record{
		ID:   "some-id",
		Data: data,
	}

	cases := map[string]struct {
		expected    *time.Time
		expectError bool
	}{
		"time": {
			expected: pointer(time.Date(2023, 3, 7, 16, 6, 5, 0, time.UTC)),
		},
		"timeISO": {
			expected: pointer(time.Date(2022, 2, 28, 6, 56, 47, 778565000, time.UTC)),
		},
		"not-a-time": {
			expectError: true,
		},
		"nullValue": {
			expected: nil,
		},
	}

	for path, c := range cases {
		t.Run(path, func(t *testing.T) {
			actual, err := record.ExtractTime(r, path)
			if c.expectError {
				assert.Error(t, err)
			} else if c.expected == nil {
				assert.Nil(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, actual)
				assert.True(t, c.expected.Equal(*actual), "expected %v and %v to be equal", c.expected, actual)
			}
		})
	}
}
