package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

func Sum(records []*api.Record, path string) (*float64, error) {
	if len(records) == 0 {
		return nil, nil
	}
	sum := 0.0
	counted := 0.0
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			sum += *number
			counted++
		}
	}
	if counted == 0 {
		return nil, nil
	}
	return pointer(sum), nil
}
