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

func TestVisit(t *testing.T) {
	data1JSON := `
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
		"nestedList": [
			{"a": 1, "b": "b1"},
			{"a": 2, "b": "b2"},
			{"a": 3, "b": "b3"}
		],
		"nestedNestedList": [
			{
				"a": [
					{"aSub": 101},
					{"aSub": 102}
				]
			},
			{
				"a": [
					{"aSub": 201},
					{"aSub": 202}
				]
			},
			{
				"a": [
					{"aSub": 301},
					{"aSub": 302}
				]
			}
		],
		"nullValue": null,
		"emptyString": "",
		"time": "2023-03-07T16:06:05Z",
		"timeISO": "2022-02-28T06:56:47.778565"
	}
	`
	data1 := map[string]any{}
	err := json.Unmarshal([]byte(data1JSON), &data1)
	require.NoError(t, err)
	data2JSON := `
	{
		"value": "string2",
		"nested": {
			"value": "nested string value2",
			"super": {
				"value": "Super Nested String Value2"
			}
		},
		"int": 456,
		"list": [
			"abc2",
			"DEF2",
			"geh2",
			"ijk2"
		],
		"nestedList": [
			{"a": 4, "b": "b4"},
			{"a": 5, "b": "b5"},
			{"a": 6, "b": "b6"}
		],
		"nestedNestedList": [
			{
				"a": [
					{"aSub": 401},
					{"aSub": 402}
				]
			},
			{
				"a": [
					{"aSub": 501},
					{"aSub": 502}
				]
			},
			{
				"a": [
					{"aSub": 601},
					{"aSub": 602}
				]
			}
		],
		"nullValue": null,
		"emptyString": "",
		"time": "2025-03-07T16:06:05Z",
		"timeISO": "2024-02-28T06:56:47.778565"
	}
	`
	data2 := map[string]any{}
	err = json.Unmarshal([]byte(data2JSON), &data2)
	require.NoError(t, err)

	records := []*api.Record{
		{
			ID:   "some-id",
			Data: data1,
		},
		{
			ID:   "some-id",
			Data: data2,
		},
	}

	cases := map[string]struct {
		expected            []any
		expectNumberErr     bool
		expectedNumber      []*float64
		expectStringErr     bool
		expectedString      []*string
		stringCaseSensitive bool
		expectTimeErr       bool
		expectedTime        []*time.Time
		expectArrayErr      bool
		expectedArray       [][]any
		useCustomData       bool
		customData          []*api.Record
	}{
		"value": {
			expected:        []any{"string", "string2"},
			expectNumberErr: true,
			expectedString:  []*string{pointer("string"), pointer("string2")},
			expectTimeErr:   true,
			expectArrayErr:  true,
		},
		"nested.value": {
			expected:        []any{"nested string value", "nested string value2"},
			expectNumberErr: true,
			expectedString:  []*string{pointer("nested string value"), pointer("nested string value2")},
			expectTimeErr:   true,
			expectArrayErr:  true,
		},
		"nested.super.value": {
			expected:        []any{"Super Nested String Value", "Super Nested String Value2"},
			expectNumberErr: true,
			expectedString:  []*string{pointer("super nested string value"), pointer("super nested string value2")},
			expectTimeErr:   true,
			expectArrayErr:  true,
		},
		"int": {
			expected:       []any{123.0, 456.0}, // json parses numbers as float64
			expectedNumber: []*float64{pointer(123.0), pointer(456.0)},
			expectedString: []*string{pointer("123"), pointer("456")},
			expectTimeErr:  true,
			expectArrayErr: true,
		},
		"list": {
			expected: []any{
				[]any{
					"abc",
					"DEF",
					"geh",
				},
				[]any{
					"abc2",
					"DEF2",
					"geh2",
					"ijk2",
				},
			},
			expectNumberErr: true,
			expectedString: []*string{
				pointer(`["abc","def","geh"]`),
				pointer(`["abc2","def2","geh2","ijk2"]`),
			},
			expectTimeErr: true,
			expectedArray: [][]any{
				{
					"abc",
					"DEF",
					"geh",
				},
				{
					"abc2",
					"DEF2",
					"geh2",
					"ijk2",
				},
			},
		},
		"list.0": {
			expected:        []any{"abc", "abc2"},
			expectNumberErr: true,
			expectedString:  []*string{pointer("abc"), pointer("abc2")},
			expectTimeErr:   true,
			expectArrayErr:  true,
		},
		"list.*": {
			expected:            []any{"abc", "DEF", "geh", "abc2", "DEF2", "geh2", "ijk2"},
			expectNumberErr:     true,
			stringCaseSensitive: true,
			expectedString: []*string{
				pointer("abc"),
				pointer("DEF"),
				pointer("geh"),
				pointer("abc2"),
				pointer("DEF2"),
				pointer("geh2"),
				pointer("ijk2"),
			},
			expectTimeErr:  true,
			expectArrayErr: true,
		},
		"nil as input": {
			useCustomData:  true,
			customData:     []*api.Record{{}, nil, {}},
			expected:       []any{nil, nil, nil},
			expectedNumber: []*float64{nil, nil, nil},
			expectedString: []*string{nil, nil, nil},
			expectedTime:   []*time.Time{nil, nil, nil},
			expectedArray:  [][]any{nil, nil, nil},
		},
		"nonexistent": {
			expected:       []any{nil, nil},
			expectedNumber: []*float64{nil, nil},
			expectedString: []*string{nil, nil},
			expectedTime:   []*time.Time{nil, nil},
			expectedArray:  [][]any{nil, nil},
		},
		"emptyString": {
			expected:        []any{"", ""},
			expectNumberErr: true,
			expectedString:  []*string{pointer(""), pointer("")},
			expectTimeErr:   true,
			expectArrayErr:  true,
		},
		"nullValue": {
			expected:       []any{nil, nil},
			expectedNumber: []*float64{nil, nil},
			expectedString: []*string{nil, nil},
			expectedTime:   []*time.Time{nil, nil},
			expectedArray:  [][]any{nil, nil},
		},
		"non.existent": {
			expected:       []any{nil, nil},
			expectedNumber: []*float64{nil, nil},
			expectedString: []*string{nil, nil},
			expectedTime:   []*time.Time{nil, nil},
			expectedArray:  [][]any{nil, nil},
		},
		"nested.nonexistent": {
			expected:       []any{nil, nil},
			expectedNumber: []*float64{nil, nil},
			expectedString: []*string{nil, nil},
			expectedTime:   []*time.Time{nil, nil},
			expectedArray:  [][]any{nil, nil},
		},
		"nested.value.nonexistent": {
			expected:       []any{nil, nil},
			expectedNumber: []*float64{nil, nil},
			expectedString: []*string{nil, nil},
			expectedTime:   []*time.Time{nil, nil},
			expectedArray:  [][]any{nil, nil},
		},
		"int.nonexistent": {
			expected:       []any{nil, nil},
			expectedNumber: []*float64{nil, nil},
			expectedString: []*string{nil, nil},
			expectedTime:   []*time.Time{nil, nil},
			expectedArray:  [][]any{nil, nil},
		},
		"list.a": {
			expected:       []any{nil, nil},
			expectedNumber: []*float64{nil, nil},
			expectedString: []*string{nil, nil},
			expectedTime:   []*time.Time{nil, nil},
			expectedArray:  [][]any{nil, nil},
		},
		"list.3": {
			expected:        []any{nil, "ijk2"},
			expectNumberErr: true,
			expectedString:  []*string{nil, pointer("ijk2")},
			expectTimeErr:   true,
			expectArrayErr:  true,
		},
		"list.-1": {
			expected:       []any{nil, nil},
			expectedNumber: []*float64{nil, nil},
			expectedString: []*string{nil, nil},
			expectedTime:   []*time.Time{nil, nil},
			expectedArray:  [][]any{nil, nil},
		},
		"nestedList.0.a": {
			expected:       []any{1.0, 4.0},
			expectedNumber: []*float64{pointer(1.0), pointer(4.0)},
			expectedString: []*string{pointer("1"), pointer("4")},
			expectTimeErr:  true,
			expectArrayErr: true,
		},
		"nestedList.*.a": {
			expected: []any{1.0, 2.0, 3.0, 4.0, 5.0, 6.0},
			expectedNumber: []*float64{
				pointer(1.0),
				pointer(2.0),
				pointer(3.0),
				pointer(4.0),
				pointer(5.0),
				pointer(6.0),
			},
			expectedString: []*string{
				pointer("1"),
				pointer("2"),
				pointer("3"),
				pointer("4"),
				pointer("5"),
				pointer("6"),
			},
			expectTimeErr:  true,
			expectArrayErr: true,
		},
		"nestedNestedList.0.a.0.aSub": {
			expected:       []any{101.0, 401.0},
			expectedNumber: []*float64{pointer(101.0), pointer(401.0)},
			expectedString: []*string{pointer("101"), pointer("401")},
			expectTimeErr:  true,
			expectArrayErr: true,
		},
		"nestedNestedList.*.a.*.aSub": {
			expected: []any{101.0, 102.0, 201.0, 202.0, 301.0, 302.0, 401.0, 402.0, 501.0, 502.0, 601.0, 602.0},
			expectedNumber: []*float64{
				pointer(101.0),
				pointer(102.0),
				pointer(201.0),
				pointer(202.0),
				pointer(301.0),
				pointer(302.0),
				pointer(401.0),
				pointer(402.0),
				pointer(501.0),
				pointer(502.0),
				pointer(601.0),
				pointer(602.0),
			},
			expectedString: []*string{
				pointer("101"),
				pointer("102"),
				pointer("201"),
				pointer("202"),
				pointer("301"),
				pointer("302"),
				pointer("401"),
				pointer("402"),
				pointer("501"),
				pointer("502"),
				pointer("601"),
				pointer("602"),
			},
			expectTimeErr:  true,
			expectArrayErr: true,
		},
		"nestedNestedList.*.a": {
			expected: []any{
				[]any{
					map[string]any{"aSub": 101.0},
					map[string]any{"aSub": 102.0},
				},
				[]any{
					map[string]any{"aSub": 201.0},
					map[string]any{"aSub": 202.0},
				},
				[]any{
					map[string]any{"aSub": 301.0},
					map[string]any{"aSub": 302.0},
				},
				[]any{
					map[string]any{"aSub": 401.0},
					map[string]any{"aSub": 402.0},
				},
				[]any{
					map[string]any{"aSub": 501.0},
					map[string]any{"aSub": 502.0},
				},
				[]any{
					map[string]any{"aSub": 601.0},
					map[string]any{"aSub": 602.0},
				},
			},
			expectNumberErr: true,
			expectedString: []*string{
				pointer(`[{"aSub":101},{"aSub":102}]`),
				pointer(`[{"aSub":201},{"aSub":202}]`),
				pointer(`[{"aSub":301},{"aSub":302}]`),
				pointer(`[{"aSub":401},{"aSub":402}]`),
				pointer(`[{"aSub":501},{"aSub":502}]`),
				pointer(`[{"aSub":601},{"aSub":602}]`),
			},
			stringCaseSensitive: true,
			expectTimeErr:       true,
			expectedArray: [][]any{
				{
					map[string]any{"aSub": 101.0},
					map[string]any{"aSub": 102.0},
				},
				{
					map[string]any{"aSub": 201.0},
					map[string]any{"aSub": 202.0},
				},
				{
					map[string]any{"aSub": 301.0},
					map[string]any{"aSub": 302.0},
				},
				{
					map[string]any{"aSub": 401.0},
					map[string]any{"aSub": 402.0},
				},
				{
					map[string]any{"aSub": 501.0},
					map[string]any{"aSub": 502.0},
				},
				{
					map[string]any{"aSub": 601.0},
					map[string]any{"aSub": 602.0},
				},
			},
		},
		"time": {
			expected:        []any{"2023-03-07T16:06:05Z", "2025-03-07T16:06:05Z"},
			expectNumberErr: true,
			expectedString:  []*string{pointer("2023-03-07t16:06:05z"), pointer("2025-03-07t16:06:05z")},
			expectedTime: []*time.Time{
				pointer(time.Date(2023, 3, 7, 16, 6, 5, 0, time.UTC)),
				pointer(time.Date(2025, 3, 7, 16, 6, 5, 0, time.UTC)),
			},
			expectArrayErr: true,
		},
		"timeISO": {
			expected:        []any{"2022-02-28T06:56:47.778565", "2024-02-28T06:56:47.778565"},
			expectNumberErr: true,
			expectedString:  []*string{pointer("2022-02-28t06:56:47.778565"), pointer("2024-02-28t06:56:47.778565")},
			expectedTime: []*time.Time{
				pointer(time.Date(2022, 2, 28, 6, 56, 47, 778565000, time.UTC)),
				pointer(time.Date(2024, 2, 28, 6, 56, 47, 778565000, time.UTC)),
			},
			expectArrayErr: true,
		},
	}

	for path, c := range cases {
		t.Run(path, func(t *testing.T) {
			input := records
			if c.useCustomData {
				input = c.customData
			}
			actual := make([]any, 0)
			err := record.Visit(input, path, func(val any, _ *api.Record) error {
				actual = append(actual, val)
				return nil
			})
			assert.NoError(t, err)
			assert.Equal(t, c.expected, actual)

			actualNumber := make([]*float64, 0)
			err = record.VisitNumber(input, path, func(val *float64, _ *api.Record) error {
				actualNumber = append(actualNumber, val)
				return nil
			})
			if c.expectNumberErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.expectedNumber, actualNumber)
			}

			actualString := make([]*string, 0)
			err = record.VisitString(input, path, c.stringCaseSensitive, func(val *string, _ *api.Record) error {
				actualString = append(actualString, val)
				return nil
			})
			if c.expectStringErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.expectedString, actualString)
			}

			actualTime := make([]*time.Time, 0)
			err = record.VisitTime(input, path, func(val *time.Time, _ *api.Record) error {
				actualTime = append(actualTime, val)
				return nil
			})
			if c.expectTimeErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				require.Len(t, actualTime, len(c.expectedTime), "time slices must have equal length")
				for i := range c.expectedTime {
					if c.expectedTime[i] == nil {
						assert.Nil(t, actualTime[i])
					} else {
						require.NotNil(t, actualTime[i])
						assert.True(t, c.expectedTime[i].Equal(*actualTime[i]), "expected %v and %v to be equal", c.expectedTime[i], *actualTime[i])
					}
				}
			}

			actualArray := make([][]any, 0)
			err = record.VisitArray(input, path, func(val []any, _ *api.Record) error {
				actualArray = append(actualArray, val)
				return nil
			})
			if c.expectArrayErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.expectedArray, actualArray)
			}
		})
	}
}
