package assertion

import (
	"math"
	"time"

	"github.com/smarty/assertions"
)

const (
	TimeType  = "time.Time"
	FloatType = "float64"
)

type AssertionFunc assertions.SoFunc

func SkipAssertion(actual any, expected ...any) string {
	return ""
}

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

func roundFloatToDecimalPlaces(num float64, decimalPlaces int) float64 {
	precision := math.Pow(10, float64(decimalPlaces))
	return math.Round(num*precision) / precision
}
