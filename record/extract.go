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
	value, err := ExtractString(record, path, true)
	if value == nil || err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339Nano, *value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// TODO: Add extract multi path string
