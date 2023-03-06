package record

import (
	"math"

	"github.com/tilotech/tilores-insights/helpers"
	api "github.com/tilotech/tilores-plugin-api"
)

func Min(records []*api.Record, path string) (*float64, error) {
	if len(records) == 0 {
		return nil, nil
	}
	min := 0.0
	counted := 0.0
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			if counted == 0 {
				min = *number
			} else {
				min = math.Min(min, *number)
			}
			counted++
		}
	}
	if counted == 0 {
		return nil, nil
	}
	return helpers.NullifyFloat(min), nil
}
