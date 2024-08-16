package assertion

import (
	"strings"
	"testing"
	"time"

	"github.com/smarty/assertions"
)

func TestSkipAssertion(t *testing.T) {
	actual := "actual"
	expected := "expected"
	ok, _ := assertions.So(actual, SkipAssertion, expected)
	if !ok {
		t.Error("SkipAssertion failed")
	}
}
func TestAssertTimeToDuration(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")

	testTable := []struct {
		name       string
		actual     any
		expected   any
		duration   time.Duration
		expectedOk bool
	}{
		{
			name:       "Test matching time",
			actual:     testTime,
			expected:   testTime,
			duration:   time.Nanosecond,
			expectedOk: true,
		},
		{
			name:       "Test not matching time",
			actual:     testTime,
			expected:   testTime.Add(time.Second),
			duration:   time.Nanosecond,
			expectedOk: false,
		},
		{
			name:       "Test matching time with seconds difference",
			actual:     testTime,
			expected:   testTime.Add(time.Second),
			duration:   time.Minute,
			expectedOk: true,
		},
		{
			name:       "Test not matching time with minutes difference",
			actual:     testTime,
			expected:   testTime.Add(time.Minute),
			duration:   time.Hour,
			expectedOk: true,
		},
		{
			name:       "Test with non time type",
			actual:     "2021-01-01T00:00:00Z",
			expected:   testTime,
			duration:   time.Nanosecond,
			expectedOk: false,
		},
		{
			name:       "Test with missing expected value",
			actual:     testTime,
			duration:   time.Nanosecond,
			expectedOk: false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ok, message := assertions.So(tt.actual, assertions.SoFunc(AssertTimeToDuration(tt.duration)), tt.expected)
			if ok != tt.expectedOk {
				t.Errorf("AssertTimeToDuration failed: %s", message)
			}
		})
	}

}

func TestAssertFloat64ToDecimalPlaces(t *testing.T) {
	testTable := []struct {
		name       string
		actual     any
		expected   any
		decimal    int
		expectedOk bool
	}{
		{
			name:       "Test matching float",
			actual:     1.23456789,
			expected:   1.23456789,
			decimal:    8,
			expectedOk: true,
		},
		{
			name:       "Test not matching float",
			actual:     1.23456,
			expected:   1.2345678,
			decimal:    6,
			expectedOk: false,
		},
		{
			name:       "Test matching float with 2 decimal places",
			actual:     1.23456789,
			expected:   1.23,
			decimal:    2,
			expectedOk: true,
		},
		{
			name:       "Test not matching float with 2 decimal places",
			actual:     1.23456789,
			expected:   1.24,
			decimal:    2,
			expectedOk: false,
		},
		{
			name:       "Test with non float type",
			actual:     "1.23456789",
			expected:   1.23456789,
			decimal:    8,
			expectedOk: false,
		},
		{
			name:       "Test with missing expected value",
			actual:     1.23456789,
			decimal:    8,
			expectedOk: false,
		},
		{
			name:       "test with 0 decimal places",
			actual:     1.23456789,
			expected:   1.0,
			decimal:    0,
			expectedOk: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ok, message := assertions.So(tt.actual, assertions.SoFunc(AssertFloat64ToDecimalPlaces(tt.decimal)), tt.expected)
			if ok != tt.expectedOk {
				t.Errorf("AssertFloat64ToDecimalPlaces failed: %s", message)
			}
		})
	}
}

func TestRoundFloatToDecimalPlaces(t *testing.T) {
	testTable := []struct {
		name          string
		num           float64
		decimalPlaces int
		expected      float64
	}{
		{
			name:          "Test truncate float to 2 decimal places",
			num:           1.23456789,
			decimalPlaces: 2,
			expected:      1.23,
		},
		{
			name:          "Test truncate float to 3 decimal places",
			num:           1.23456789,
			decimalPlaces: 3,
			expected:      1.235,
		},
		{
			name:          "Test truncate float to 4 decimal places",
			num:           1.23,
			decimalPlaces: 4,
			expected:      1.23,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			actual := roundFloatToDecimalPlaces(tt.num, tt.decimalPlaces)
			if actual != tt.expected {
				t.Errorf("TruncateFloatDecimalPlaces failed: expected %f, got %f", tt.expected, actual)
			}
		})
	}

}

func TestAssertFloat64WithTolerance(t *testing.T) {
	testTable := []struct {
		name       string
		actual     any
		expected   any
		tolerance  float64
		expectedOk bool
	}{
		{
			name:       "Test matching float",
			actual:     1.23456789,
			expected:   1.23456789,
			tolerance:  0.00000001,
			expectedOk: true,
		},
		{
			name:       "Test not matching float",
			actual:     1.23456789,
			expected:   1.2345678,
			tolerance:  0.00000001,
			expectedOk: false,
		},
		{
			name:       "Test matching float with tolerance",
			actual:     1.23456789,
			expected:   1.2345678,
			tolerance:  0.0001,
			expectedOk: true,
		},
		{
			name:       "Test with non float type",
			actual:     "1.23456789",
			expected:   1.23456789,
			tolerance:  0.00000001,
			expectedOk: false,
		},
		{
			name:       "Test with missing expected value",
			actual:     1.23456789,
			tolerance:  0.00000001,
			expectedOk: false,
		},
		{
			name:       "Test with missing tolerance",
			actual:     1.23456789,
			expected:   1.23456789,
			expectedOk: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ok, message := assertions.So(tt.actual, assertions.SoFunc(AssertFloat64WithTolerance(tt.tolerance)), tt.expected)
			if ok != tt.expectedOk {
				t.Errorf("AssertFloat64WithTolerance failed: %s", message)
			}
		})
	}
}

func TestAssertStringWithCleanup(t *testing.T) {
	testTable := []struct {
		name       string
		actual     any
		expected   any
		cleanup    func(string) string
		expectedOk bool
	}{
		{
			name:       "Test matching string",
			actual:     "test",
			expected:   "test",
			cleanup:    nil,
			expectedOk: true,
		},
		{
			name:       "Test not matching string",
			actual:     "test",
			expected:   "test1",
			cleanup:    nil,
			expectedOk: false,
		},
		{
			name:       "Test matching string with cleanup",
			actual:     "test",
			expected:   "test",
			cleanup:    func(s string) string { return s },
			expectedOk: true,
		},
		{
			name:       "Test not matching string with cleanup",
			actual:     "  test  ",
			expected:   "test ",
			cleanup:    strings.TrimSpace,
			expectedOk: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ok, message := assertions.So(tt.actual, assertions.SoFunc(AssertStringWithCleanup(tt.cleanup)), tt.expected)
			if ok != tt.expectedOk {
				t.Errorf("AssertStringWithCleanup failed: %s", message)
			}
		})
	}
}

func TestSkipAssertionIf(t *testing.T) {
	testTable := []struct {
		name       string
		actual     any
		expected   any
		condition  func(actual any, expected any) bool
		expectedOk bool
	}{
		{
			name:       "Test skip assertion",
			actual:     "test",
			expected:   "test ",
			condition:  func(actual any, expected any) bool { return true },
			expectedOk: true,
		},
		{
			name:       "Test not skip assertion",
			actual:     "test",
			expected:   "test ",
			condition:  func(actual any, expected any) bool { return false },
			expectedOk: false,
		},
		{
			name:       "Test skip assertion with nil expected value",
			actual:     "test",
			expected:   nil,
			condition:  IsNilExpected,
			expectedOk: true,
		},
		{
			name:       "Test not skip assertion with nil actual value",
			actual:     nil,
			expected:   "test",
			condition:  IsNilActual,
			expectedOk: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ok, message := assertions.So(tt.actual, assertions.SoFunc(SkipAssertionIf(tt.condition, nil)), tt.expected)
			if ok != tt.expectedOk {
				t.Errorf("SkipAssertionIf failed: %s", message)
			}
		})
	}
}
