package assertion

import (
	"testing"

	"github.com/smarty/assertions"
)

func Test_AssertUnorderedSlice(t *testing.T) {
	testTable := []struct {
		name             string
		actual           []int
		expected         []int
		compare          func(actual, expected int) bool
		expectedActual   []int
		expectedExpected []int
	}{
		{
			name:             "Test unordered slices",
			actual:           []int{1, 2, 3},
			expected:         []int{3, 2, 1},
			compare:          func(actual, expected int) bool { return actual < expected },
			expectedActual:   []int{1, 2, 3},
			expectedExpected: []int{1, 2, 3},
		},
		{
			name:             "Test ordered slices",
			actual:           []int{1, 2, 3},
			expected:         []int{1, 2, 3},
			compare:          func(actual, expected int) bool { return actual < expected },
			expectedActual:   []int{1, 2, 3},
			expectedExpected: []int{1, 2, 3},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			actual, expected := AssertUnorderedSlice(test.compare)(test.actual, test.expected)
			ok, message := assertions.So(test.expectedActual, assertions.ShouldEqual, actual)
			if !ok {
				t.Errorf(message)
			}
			ok, message = assertions.So(test.expectedExpected, assertions.ShouldEqual, expected)
			if !ok {
				t.Errorf(message)
			}
		})
	}

	testTable2 := []struct {
		name             string
		actual           []string
		expected         []string
		compare          func(actual, expected string) bool
		expectedActual   []string
		expectedExpected []string
	}{
		{
			name:             "Test unordered slices",
			actual:           []string{"a", "b", "c"},
			expected:         []string{"c", "b", "a"},
			compare:          func(actual, expected string) bool { return actual < expected },
			expectedActual:   []string{"a", "b", "c"},
			expectedExpected: []string{"a", "b", "c"},
		},
		{
			name:             "Test ordered slices",
			actual:           []string{"a", "b", "c"},
			expected:         []string{"a", "b", "c"},
			compare:          func(actual, expected string) bool { return actual < expected },
			expectedActual:   []string{"a", "b", "c"},
			expectedExpected: []string{"a", "b", "c"},
		},
	}

	for _, test := range testTable2 {
		t.Run(test.name, func(t *testing.T) {
			actual, expected := AssertUnorderedSlice(test.compare)(test.actual, test.expected)
			ok, message := assertions.So(test.expectedActual, assertions.ShouldEqual, actual)
			if !ok {
				t.Errorf(message)
			}
			ok, message = assertions.So(test.expectedExpected, assertions.ShouldEqual, expected)
			if !ok {
				t.Errorf(message)
			}
		})
	}
}
