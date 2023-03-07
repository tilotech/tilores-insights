package record

import (
	"math"

	api "github.com/tilotech/tilores-plugin-api"
)

func Max(records []*api.Record, path string) (*float64, error) {
	if len(records) == 0 {
		return nil, nil
	}
	max := 0.0
	counted := 0.0
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			if counted == 0 {
				max = *number
			} else {
				max = math.Max(max, *number)
			}
			counted++
		}
	}
	if counted == 0 {
		return nil, nil
	}
	return pointer(max), nil
}
