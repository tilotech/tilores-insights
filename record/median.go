package record

import (
	"sort"

	api "github.com/tilotech/tilores-plugin-api"
)

func Median(records []*api.Record, path string) (*float64, error) {
	if len(records) == 0 {
		return nil, nil
	}
	numbers := []float64{}
	for _, record := range records {
		number, err := ExtractNumber(record, path)
		if err != nil {
			return nil, err
		}
		if number != nil {
			numbers = append(numbers, *number)
		}
	}
	counted := len(numbers)
	if counted == 0 {
		return nil, nil
	}
	sort.Float64s(numbers)
	if counted%2 == 1 {
		return pointer(numbers[(counted / 2)]), nil
	}
	return pointer((numbers[counted/2] + numbers[(counted/2)-1]) / 2), nil
}
