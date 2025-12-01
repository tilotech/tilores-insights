package record

import (
	"strings"

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
	recordKeys, err := groupRecordStringKeys(records, paths, caseSensitive)
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

func groupRecordStringKeys(records []*api.Record, paths []string, caseSensitive bool) ([]string, error) {
	recordKeys := make(map[string]*strings.Builder)
	for _, path := range paths {
		err := VisitString(records, path, caseSensitive, func(s *string, record *api.Record) error {
			if record == nil {
				return nil
			}
			sb, ok := recordKeys[record.ID]
			if !ok {
				sb = &strings.Builder{}
				recordKeys[record.ID] = sb
			}
			if s == nil {
				sb.WriteString(":n:")
			} else {
				sb.WriteString(":s:")
				sb.WriteString(*s)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	keys := make([]string, len(records))
	for i, record := range records {
		var key string
		if record != nil {
			sb, ok := recordKeys[record.ID]
			if ok {
				key = sb.String()
			}
		}
		keys[i] = key
	}
	return keys, nil
}
