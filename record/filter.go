package record

import (
	"fmt"
	"regexp"
	"strconv"
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

	Equals any
	IsNull *bool

	StartsWith *string
	EndsWith   *string
	LikeRegex  *string

	LessThan      *float64
	LessEquals    *float64
	GreaterThan   *float64
	GreaterEquals *float64

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
		keep, err := checkFilterCondition(record, condition)
		if err != nil {
			return false, err
		}
		if condition.Invert != nil && *condition.Invert {
			keep = !keep
		}
		if !keep {
			return false, nil
		}
	}
	return true, nil
}

func checkFilterCondition(record *api.Record, condition *FilterCondition) (bool, error) {
	if !checkFilterCriteriaIsNull(record, condition) {
		return false, nil
	}
	if keep, err := checkFilterStringCriteria(record, condition); !keep || err != nil {
		return keep, err
	}
	if keep, err := checkFilterNumericCriteria(record, condition); !keep || err != nil {
		return keep, err
	}
	if keep, err := checkFilterTimeCriteria(record, condition); !keep || err != nil {
		return keep, err
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

func hasFilterStringCriteria(condition *FilterCondition) bool {
	return condition.Equals != nil ||
		condition.StartsWith != nil ||
		condition.EndsWith != nil ||
		condition.LikeRegex != nil
}

func checkFilterStringCriteria(record *api.Record, condition *FilterCondition) (bool, error) {
	if !hasFilterStringCriteria(condition) {
		return true, nil
	}

	caseSensitive := false
	if condition.CaseSensitive != nil {
		caseSensitive = *condition.CaseSensitive
	}
	value, err := ExtractString(record, condition.Path, caseSensitive)
	if err != nil {
		return false, err
	}

	if keep, err := checkFilterCriteriaEqual(value, condition.Equals, caseSensitive); !keep || err != nil {
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

	valueF, err := strconv.ParseFloat(*value, 64)
	if err == nil {
		equalF, err := strconv.ParseFloat(*equalString, 64)
		if err != nil {
			return false, nil
		}
		return valueF == equalF, nil
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

func hasFilterNumericCriteria(condition *FilterCondition) bool {
	return condition.LessThan != nil ||
		condition.LessEquals != nil ||
		condition.GreaterThan != nil ||
		condition.GreaterEquals != nil
}

func checkFilterNumericCriteria(record *api.Record, condition *FilterCondition) (bool, error) {
	if !hasFilterNumericCriteria(condition) {
		return true, nil
	}

	value, err := ExtractNumber(record, condition.Path)
	if err != nil {
		return false, err
	}

	if !checkFilterCriteriaCompare(value, condition.LessThan, checkFilterOpLessThan) {
		return false, nil
	}
	if !checkFilterCriteriaCompare(value, condition.LessEquals, checkFilterOpLessEqual) {
		return false, nil
	}
	if !checkFilterCriteriaCompare(value, condition.GreaterThan, checkFilterOpGreaterThan) {
		return false, nil
	}
	if !checkFilterCriteriaCompare(value, condition.GreaterEquals, checkFilterOpGreaterEqual) {
		return false, nil
	}

	return true, nil
}

func checkFilterOpLessThan(a, b float64) bool {
	return a < b
}

func checkFilterOpLessEqual(a, b float64) bool {
	return a <= b
}

func checkFilterOpGreaterThan(a, b float64) bool {
	return a > b
}

func checkFilterOpGreaterEqual(a, b float64) bool {
	return a >= b
}

func hasFilterTimeCriteria(condition *FilterCondition) bool {
	return condition.After != nil ||
		condition.Since != nil ||
		condition.Before != nil ||
		condition.Until != nil
}

func checkFilterTimeCriteria(record *api.Record, condition *FilterCondition) (bool, error) {
	if !hasFilterTimeCriteria(condition) {
		return true, nil
	}

	value, err := ExtractTime(record, condition.Path)
	if err != nil {
		return false, err
	}

	if !checkFilterCriteriaCompare(value, condition.After, checkFilterOpAfter) {
		return false, nil
	}
	if !checkFilterCriteriaCompare(value, condition.Since, checkFilterOpSince) {
		return false, nil
	}
	if !checkFilterCriteriaCompare(value, condition.Before, checkFilterOpBefore) {
		return false, nil
	}
	if !checkFilterCriteriaCompare(value, condition.Until, checkFilterOpUntil) {
		return false, nil
	}

	return true, nil
}

func checkFilterOpAfter(a, b time.Time) bool {
	return a.After(b)
}

func checkFilterOpSince(a, b time.Time) bool {
	return a.Equal(b) || a.After(b)
}

func checkFilterOpBefore(a, b time.Time) bool {
	return a.Before(b)
}

func checkFilterOpUntil(a, b time.Time) bool {
	return a.Equal(b) || a.Before(b)
}

func checkFilterCriteriaCompare[T any](value *T, test *T, op func(a, b T) bool) bool {
	if test == nil {
		return true
	}
	if value == nil {
		return false
	}
	return op(*value, *test)
}
