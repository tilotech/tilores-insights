package record

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	api "github.com/tilotech/tilores-plugin-api"
)

// Extract provides the value of a record for the given path.
func Extract(record *api.Record, path string) any {
	if record == nil {
		return nil
	}
	pathParts := strings.Split(path, ".")
	return extract(record.Data, pathParts)
}

func extract(data any, pathParts []string) any {
	if len(pathParts) == 0 {
		return data
	}
	subPath, pathParts := pathParts[0], pathParts[1:]

	if mapData, ok := data.(map[string]any); ok {
		mapValue, ok := mapData[subPath]
		if !ok {
			return nil
		}
		return extract(mapValue, pathParts)
	}
	if listData, ok := data.([]any); ok {
		i, err := strconv.Atoi(subPath)
		if err != nil {
			return nil
		}
		if i < 0 || len(listData) <= i {
			return nil
		}
		return extract(listData[i], pathParts)
	}
	return nil
}

// ExtractNumber provides a numeric value of a record for the given path.
func ExtractNumber(record *api.Record, path string) (*float64, error) {
	val := Extract(record, path)
	return validateNumber(val, path)
}

func validateNumber(val any, path string) (*float64, error) {
	if val == nil {
		return nil, nil
	}
	if number, ok := val.(float64); ok {
		return &number, nil
	}
	if s, ok := val.(string); ok {
		number, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return &number, nil
	}
	return nil, fmt.Errorf("invalid type while extracting number from path %v, expected numeric value but received %T", path, val)
}

// ExtractString provides a string value of a record for the given path.
//
// If the value of that path is an array or a map, it will stringify the value
// into JSON.
func ExtractString(record *api.Record, path string, caseSensitive bool) (*string, error) {
	val := Extract(record, path)
	return validateString(val, caseSensitive)
}

func validateString(val any, caseSensitive bool) (*string, error) {
	if val == nil {
		return nil, nil
	}
	return valueToString(val, caseSensitive)
}

func valueToString(val any, caseSensitive bool) (*string, error) {
	switch val.(type) {
	case map[string]any, []any:
		marshal, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		jsonString := string(marshal)
		if !caseSensitive {
			jsonString = strings.ToLower(jsonString)
		}
		return &jsonString, nil
	}
	stringVal := fmt.Sprintf("%v", val)
	if !caseSensitive {
		stringVal = strings.ToLower(stringVal)
	}
	return &stringVal, nil
}

// ExtractTime provides a time value of a record for the given path.
func ExtractTime(record *api.Record, path string) (*time.Time, error) {
	val, err := ExtractString(record, path, true)
	if err != nil {
		return nil, err
	}
	return validateTime(val)
}

func validateTime(val *string) (*time.Time, error) {
	if val == nil {
		return nil, nil
	}
	return parseTime(*val)
}

var supportedTimeFormats = [...]string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999",
}

func parseTime(t string) (*time.Time, error) {
	var parsed time.Time
	var err error
	for _, format := range supportedTimeFormats {
		parsed, err = time.Parse(format, t)
		if err == nil {
			return &parsed, nil
		}
	}
	return &parsed, err
}

// ExtractArray provides an array value of a record for the given path.
func ExtractArray(record *api.Record, path string) ([]any, error) {
	val := Extract(record, path)
	return validateArray(val, path)
}

func validateArray(val any, path string) ([]any, error) {
	if val == nil {
		return nil, nil
	}
	if arr, ok := val.([]any); ok {
		return arr, nil
	}
	return nil, fmt.Errorf("invalid type while extracting array from path %v, received %T", path, val)
}

func extractStringKeys(records []*api.Record, paths []string, caseSensitive bool) ([]string, error) {
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
