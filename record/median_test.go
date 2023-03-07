package record_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tilotech/tilores-insights/record"
	api "github.com/tilotech/tilores-plugin-api"
)

func TestMedian(t *testing.T) {
	cases := map[string]struct {
		records     []*api.Record
		expected    *float64
		expectError bool
	}{
		"empty list": {
			records:  []*api.Record{},
			expected: nil,
		},
		"nil list": {
			records:  nil,
			expected: nil,
		},
		"list with all nil values": {
			records: []*api.Record{
				nil,
				nil,
			},
			expected: nil,
		},
		"list with odd number of different values": {
			records: []*api.Record{
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "5",
					},
				},
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "10.0",
					},
				},
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "1",
					},
				},
				{
					ID:   "someid",
					Data: map[string]interface{}{},
				},
				nil,
			},
			expected: pointer(5.0),
		},
		"list with even number of different values": {
			records: []*api.Record{
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "4",
					},
				},
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "10.0",
					},
				},
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "6",
					},
				},
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "1",
					},
				},
				{
					ID:   "someid",
					Data: map[string]interface{}{},
				},
				nil,
			},
			expected: pointer(5.0),
		}, "list with one value": {
			records: []*api.Record{
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "4",
					},
				},
			},
			expected: pointer(4.0),
		},
		"list with non numbers values causes an error": {
			records: []*api.Record{
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "10.0",
					},
				},
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "not number",
					},
				},
				{
					ID: "someid",
					Data: map[string]interface{}{
						"num": "0",
					},
				},
			},
			expectError: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := record.Median(c.records, "num")
			if c.expectError {
				assert.Error(t, err)
			} else if c.expected == nil {
				assert.Nil(t, actual)
			} else {
				assert.NotNil(t, actual)
				assert.Equal(t, *c.expected, *actual)
			}
		})
	}
}

func pointer[T any](v T) *T {
	return &v
}
