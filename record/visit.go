package record

import (
	"strconv"
	"strings"
	"time"

	api "github.com/tilotech/tilores-plugin-api"
)

// Visit visits each record according to the provided path and calls the visitor
// for each value.
//
// If a path refers to a non-existent value, the visitor will be called with a
// nil value.
//
// A path is a combination of of properties or array indices, separated by a
// dot, e.g. foo.bar.0.a
// A path may contain a wildcard (*) instead of an array indices. In that case
// each value of the array will be traversed, resulting in more than one
// invocations of the visitor.
//
// If a visitor returns an error, no further elements will be visited and the
// function will return that error.
func Visit(records []*api.Record, path string, visitor func(val any, record *api.Record) error) error {
	pathParts := strings.Split(path, ".")
	for _, record := range records {
		var data map[string]any
		if record != nil {
			data = record.Data
		}
		if err := visit(data, pathParts, record, visitor); err != nil {
			return err
		}
	}
	return nil
}

func visit(data any, pathParts []string, record *api.Record, visitor func(val any, record *api.Record) error) error {
	if len(pathParts) == 0 {
		return visitor(data, record)
	}
	subPath, pathParts := pathParts[0], pathParts[1:]

	if mapData, ok := data.(map[string]any); ok {
		mapValue, ok := mapData[subPath]
		if !ok {
			return visitor(nil, record)
		}
		return visit(mapValue, pathParts, record, visitor)
	}
	if listData, ok := data.([]any); ok {
		return visitSlice(listData, subPath, pathParts, record, visitor)
	}
	return visitor(nil, record)
}

func visitSlice(listData []any, subPath string, pathParts []string, record *api.Record, visitor func(val any, record *api.Record) error) error {
	if subPath == "*" {
		for i := range listData {
			if err := visit(listData[i], pathParts, record, visitor); err != nil {
				return err
			}
		}
		return nil
	}
	i, err := strconv.Atoi(subPath)
	if err != nil {
		return visitor(nil, record)
	}
	if i < 0 || len(listData) <= i {
		return visitor(nil, record)
	}
	return visit(listData[i], pathParts, record, visitor)
}

// VisitNumber is a type-safe variant of Visit.
//
// If a found value cannot be converted into a number, then an error is returned.
func VisitNumber(records []*api.Record, path string, visitor func(val *float64, record *api.Record) error) error {
	return Visit(records, path, func(val any, record *api.Record) error {
		v, err := validateNumber(val, path)
		if err != nil {
			return err
		}
		return visitor(v, record)
	})
}

// VisitString is a type-safe variant of Visit.
//
// If a found value cannot be converted into a string, then an error is returned.
func VisitString(records []*api.Record, path string, caseSensitive bool, visitor func(val *string, record *api.Record) error) error {
	return Visit(records, path, func(val any, record *api.Record) error {
		v, err := validateString(val, caseSensitive)
		if err != nil {
			return err
		}
		return visitor(v, record)
	})
}

// VisitTime is a type-safe variant of Visit.
//
// If a found value cannot be converted into a time, then an error is returned.
func VisitTime(records []*api.Record, path string, visitor func(val *time.Time, record *api.Record) error) error {
	return VisitString(records, path, true, func(val *string, record *api.Record) error {
		v, err := validateTime(val)
		if err != nil {
			return err
		}
		return visitor(v, record)
	})
}

// VisitArray is a type-safe variant of Visit.
//
// If a found value cannot be converted into an array, then an error is returned.
func VisitArray(records []*api.Record, path string, visitor func(val []any, record *api.Record) error) error {
	return Visit(records, path, func(val any, record *api.Record) error {
		v, err := validateArray(val, path)
		if err != nil {
			return err
		}
		return visitor(v, record)
	})
}
