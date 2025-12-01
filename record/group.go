package record

import (
	api "github.com/tilotech/tilores-plugin-api"
)

// Group returns a list of record lists where the records have been grouped by
// the provided paths.
func Group(records []*api.Record, paths []string, caseSensitive bool) ([][]*api.Record, error) {
	if len(paths) == 0 {
		return [][]*api.Record{records}, nil
	}
	groups := make([][]*api.Record, 0)
	idx := make(map[string]int, len(records))
	recordKeys, err := extractStringKeys(records, paths, caseSensitive)
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		key := recordKeys[i]
		i, ok := idx[key]
		if !ok {
			i = len(groups)
			groups = append(groups, []*api.Record{})
			idx[key] = i
		}
		groups[i] = append(groups[i], record)
	}

	return groups, nil
}
