package record

import (
	"sort"

	api "github.com/tilotech/tilores-plugin-api"
)

// SortCriteria defines how to sort.
type SortCriteria struct {
	Path string
	ASC  bool
}

// Sort returns a new RecordInsights that contains the records ordered by the
// provided SortCriteria.
func Sort(records []*api.Record, criteria []*SortCriteria) ([]*api.Record, error) {
	if len(criteria) == 0 {
		return records, nil
	}

	data, err := sortCollectData(records, criteria)
	if err != nil {
		return nil, err
	}

	sortedIdx := make([]int, len(records))
	for i := range sortedIdx {
		sortedIdx[i] = i
	}

	sort.SliceStable(sortedIdx, func(i, j int) bool {
		i2 := sortedIdx[i]
		j2 := sortedIdx[j]
		for k := range criteria {
			if data[k].useNumber {
				a := data[k].numberValues[i2]
				b := data[k].numberValues[j2]
				if sortIsEqual(a, b) {
					continue
				}
				return sortLess(a, b, criteria[k].ASC)
			}
			a := data[k].stringValues[i2]
			b := data[k].stringValues[j2]
			if sortIsEqual(a, b) {
				continue
			}
			return sortLess(a, b, criteria[k].ASC)
		}
		return false
	})

	sortedRecords := make([]*api.Record, len(records))
	for i, j := range sortedIdx {
		sortedRecords[i] = records[j]
	}

	return sortedRecords, nil
}

func sortCollectData(records []*api.Record, criteria []*SortCriteria) ([]sortValues, error) {
	data := make([]sortValues, 0, len(criteria))
	for range criteria {
		data = append(data, sortValues{
			useNumber:    true,
			numberValues: make([]*float64, 0, len(records)),
			stringValues: make([]*string, 0, len(records)),
		})
	}

	for _, record := range records {
		for i, c := range criteria {
			if data[i].useNumber {
				val, err := ExtractNumber(record, c.Path)
				if err != nil {
					data[i].useNumber = false
				} else {
					data[i].numberValues = append(data[i].numberValues, val)
				}
			}
			val, err := ExtractString(record, c.Path, false)
			if err != nil {
				return nil, err
			}
			data[i].stringValues = append(data[i].stringValues, val)
		}
	}
	return data, nil
}

func sortIsEqual[T float64 | string](a, b *T) bool {
	if a == nil {
		return b == nil
	}
	return b != nil && *a == *b
}

func sortLess[T float64 | string](a, b *T, sortASC bool) bool {
	if a == nil {
		return !sortASC
	}
	if b == nil {
		return sortASC
	}
	if sortASC {
		return *a < *b
	}
	return *a > *b
}

type sortValues struct {
	useNumber    bool
	numberValues []*float64
	stringValues []*string
}
