package record

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	api "github.com/tilotech/tilores-plugin-api"
)

// FilterCondition defines the criteria for filtering a record list.
//
// Each filter condition must have a filter upon which the checks are applied and should have at least one criteria
// defined.
//
// Some criteria are mutually exclusive due to either logical reasons or type constraints. E.g. lessThan and after
// cannot be used together due to different type expectations.
type FilterCondition struct {
	Path string

	Equal  any
	IsNull *bool

	StartsWith *string
	EndsWith   *string
	LikeRegex  *string

	LessThan     *float64
	LessEqual    *float64
	GreaterThan  *float64
	GreaterEqual *float64

	After  *time.Time
	Since  *time.Time
	Before *time.Time
	Until  *time.Time

	Invert        *bool
	CaseSensitive *bool
}

// Filter returns a new RecordInsights that only contains the records for which
// the FilterCondition applies.
//
// If no records match the filter condition, then an empty RecordInsights is
// returned.
func Filter(records []*api.Record, conditions []*FilterCondition) ([]*api.Record, error) {
	if len(conditions) == 0 {
		return records, nil
	}

	filteredRecords := make([]*api.Record, 0, len(records))

	for _, record := range records {
		keep, err := checkFilterConditions(record, conditions)
		if err != nil {
			return nil, err
		}
		if keep {
			filteredRecords = append(filteredRecords, record)
		}
	}

	return filteredRecords, nil
}

func checkFilterConditions(record *api.Record, conditions []*FilterCondition) (bool, error) {
	for _, condition := range conditions {
		if !checkFilterCriteriaIsNull(record, condition) {
			return false, nil
		}
		if keep, err := checkFilterStringCriteria(record, condition); !keep || err != nil {
			return keep, err
		}
	}
	return true, nil
}

func checkFilterCriteriaIsNull(record *api.Record, condition *FilterCondition) bool {
	if condition.IsNull == nil {
		return true
	}
	value := Extract(record, condition.Path)
	if *condition.IsNull {
		return value == nil
	}
	return value != nil
}

func checkFilterStringCriteria(record *api.Record, condition *FilterCondition) (bool, error) {
	caseSensitive := false
	if condition.CaseSensitive != nil {
		caseSensitive = *condition.CaseSensitive
	}
	value, err := ExtractString(record, condition.Path, caseSensitive)
	if err != nil {
		return false, err
	}

	if keep, err := checkFilterCriteriaEqual(value, condition.Equal, caseSensitive); !keep || err != nil {
		return keep, err
	}
	if !checkFilterCriteriaStartsWith(value, condition.StartsWith, caseSensitive) {
		return false, nil
	}
	if !checkFilterCriteriaEndsWith(value, condition.EndsWith, caseSensitive) {
		return false, nil
	}
	if keep, err := checkFilterCriteriaLikeRegex(value, condition.LikeRegex, caseSensitive); !keep || err != nil {
		return keep, err
	}

	return true, nil
}

func checkFilterCriteriaEqual(value *string, equal any, caseSensitive bool) (bool, error) {
	if equal == nil {
		return true, nil
	}
	if value == nil {
		return false, nil
	}
	equalString, err := valueToString(equal, caseSensitive)
	if equalString == nil || err != nil {
		return false, err
	}
	return *value == *equalString, nil
}

func checkFilterCriteriaStartsWith(value *string, startsWith *string, caseSensitive bool) bool {
	if startsWith == nil {
		return true
	}
	if value == nil {
		return false
	}
	sw := *startsWith
	if !caseSensitive {
		sw = strings.ToLower(sw)
	}
	return strings.HasPrefix(*value, sw)
}

func checkFilterCriteriaEndsWith(value *string, endsWith *string, caseSensitive bool) bool {
	if endsWith == nil {
		return true
	}
	if value == nil {
		return false
	}
	sw := *endsWith
	if !caseSensitive {
		sw = strings.ToLower(sw)
	}
	return strings.HasSuffix(*value, sw)
}

func checkFilterCriteriaLikeRegex(value *string, likeRegex *string, caseSensitive bool) (bool, error) {
	if likeRegex == nil {
		return true, nil
	}
	if value == nil {
		return false, nil
	}
	likeRegexString := *likeRegex
	if !caseSensitive {
		likeRegexString = fmt.Sprintf("(?i)%v", likeRegexString)
	}
	r, err := regexp.Compile(likeRegexString)
	if err != nil {
		return false, err
	}
	return r.MatchString(*value), nil
}
