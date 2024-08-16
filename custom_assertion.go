package assertion

import (
	"math"
	"time"

	"github.com/smarty/assertions"
)

const (
	TimeType   = "time.Time"
	FloatType  = "float64"
	IntType    = "int"
	Int64Type  = "int64"
	stringType = "string"
)

type AssertionFunc assertions.SoFunc

// SkipAssertion is a custom assertion function that always returns true.
func SkipAssertion(actual any, expected ...any) string {
	return ""
}

// AssertTimeToDuration is a custom assertion function that truncates time to the specified duration before comparing.
// this function is using time.Truncate() to truncate the time to the specified duration.
func AssertTimeToDuration(duration time.Duration) AssertionFunc {
	return func(actual any, expected ...any) string {
		if len(expected) == 0 || expected[0] == nil {
			return "expected value is missing"
		}

		if act, ok := actual.(time.Time); ok {
			actual = act.Truncate(duration)
		}

		if exp, ok := expected[0].(time.Time); ok {
			expected[0] = exp.Truncate(duration)
		}
		return assertions.ShouldEqual(actual, expected[0])
	}
}

// AssertFloat64ToDecimalPlaces is a custom assertion function that rounds float64 to the specified decimal places before comparing.
// this function is rounding the float64 to the specified decimal places before comparing.
func AssertFloat64ToDecimalPlaces(decimalPlaces int) AssertionFunc {
	return func(actual any, expected ...any) string {
		if len(expected) == 0 || expected[0] == nil {
			return "expected value is missing"
		}
		if act, ok := actual.(float64); ok {
			actual = roundFloatToDecimalPlaces(act, decimalPlaces)
		}
		if exp, ok := expected[0].(float64); ok {
			expected[0] = roundFloatToDecimalPlaces(exp, decimalPlaces)
		}
		return assertions.ShouldEqual(actual, expected[0])
	}
}

// AssertFloat64WithTolerance is a custom assertion function that compares float64 values with a tolerance.
// tolerance is in format of 0.0001
// if tolerance is 0, it will compare the float64 values as is.
func AssertFloat64WithTolerance(tolerance float64) AssertionFunc {
	return func(actual any, expected ...any) string {
		if len(expected) == 0 || expected[0] == nil {
			return "expected value is missing"
		}
		if act, ok := actual.(float64); ok && tolerance > 0 {
			if exp, ok := expected[0].(float64); ok {
				return assertions.ShouldAlmostEqual(act, exp, tolerance)
			}
		}
		return assertions.ShouldEqual(actual, expected[0])
	}
}

func roundFloatToDecimalPlaces(num float64, decimalPlaces int) float64 {
	precision := math.Pow(10, float64(decimalPlaces))
	return math.Round(num*precision) / precision
}

// AssertStringWithCleanup is a custom assertion function that cleans up the string before comparing.
// cleanup function is a function that takes a string and returns a string after cleanup.
// if cleanup is nil, it will compare the string as is.
func AssertStringWithCleanup(cleanup func(string) string) AssertionFunc {
	return func(actual any, expected ...any) string {
		if len(expected) == 0 || expected[0] == nil {
			return "expected value is missing"
		}
		if act, ok := actual.(string); ok {
			if exp, ok := expected[0].(string); ok && cleanup != nil {
				return assertions.ShouldEqual(cleanup(act), cleanup(exp))
			}
		}
		return assertions.ShouldEqual(actual, expected[0])
	}
}
