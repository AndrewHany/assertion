package assertion

import (
	"math"
	"reflect"
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
		// check if actual is a time.Time
		if reflect.TypeOf(actual).String() != TimeType || reflect.TypeOf(expected[0]).String() != TimeType {
			return assertions.ShouldEqual(actual, expected[0])
		}
		actualTime := (actual.(time.Time)).Truncate(duration)
		expectedTime := (expected[0]).(time.Time).Truncate(duration)
		return assertions.ShouldEqual(actualTime, expectedTime)
	}
}

func AssertFloat64ToDecimalPlaces(decimalPlaces int) AssertionFunc {
	return func(actual any, expected ...any) string {
		if len(expected) == 0 || expected[0] == nil {
			return "expected value is missing"
		}
		// check if actual is a float64
		if reflect.TypeOf(actual).String() != FloatType || reflect.TypeOf(expected[0]).String() != FloatType {
			return assertions.ShouldEqual(actual, expected[0])
		}
		actualFloat := truncateFloatDecimalPlaces(actual.(float64), decimalPlaces)
		expectedFloat := truncateFloatDecimalPlaces(expected[0].(float64), decimalPlaces)
		return assertions.ShouldEqual(actualFloat, expectedFloat)
	}
}

func truncateFloatDecimalPlaces(num float64, decimalPlaces int) float64 {
	return math.Floor(num*math.Pow(10, float64(decimalPlaces))) / math.Pow(10, float64(decimalPlaces))
}
