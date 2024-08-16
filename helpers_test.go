package assertion

import (
	"testing"
)

func TestIsNilExpected(t *testing.T) {
	testTable := []struct {
		name       string
		actual     any
		expected   any
		expectedOk bool
	}{
		{
			name:       "Test nil expected",
			actual:     "test",
			expected:   nil,
			expectedOk: true,
		},
		{
			name:       "Test not nil expected",
			actual:     "test",
			expected:   "test",
			expectedOk: false,
		},
		{
			name:       "Test nil actual",
			actual:     nil,
			expected:   "test",
			expectedOk: false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ok := IsNilExpected(tt.actual, tt.expected)
			if ok != tt.expectedOk {
				t.Errorf("IsNilExpected failed")
			}
		})
	}
}

func TestIsNilActual(t *testing.T) {
	testTable := []struct {
		name       string
		actual     any
		expected   any
		expectedOk bool
	}{
		{
			name:       "Test nil actual",
			actual:     nil,
			expected:   "test",
			expectedOk: true,
		},
		{
			name:       "Test not nil actual",
			actual:     "test",
			expected:   "test",
			expectedOk: false,
		},
		{
			name:       "Test nil expected",
			actual:     "test",
			expected:   nil,
			expectedOk: false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ok := IsNilActual(tt.actual, tt.expected)
			if ok != tt.expectedOk {
				t.Errorf("IsNilActual failed")
			}
		})
	}
}
