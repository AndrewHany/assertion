package assertion

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestAssertValue(t *testing.T) {
	testTable := []struct {
		name            string
		path            string
		customAssertion AssertionFunc
		value1          reflect.Value
		value2          reflect.Value
		expectedMatch   bool
		expectedMessage string
	}{
		{
			name:            "Test with nil values",
			value1:          reflect.ValueOf(nil),
			value2:          reflect.ValueOf(nil),
			expectedMatch:   true,
			expectedMessage: "",
		},
		{
			name:            "Test with nil value1",
			path:            "$",
			value1:          reflect.ValueOf(nil),
			value2:          reflect.ValueOf(1),
			expectedMatch:   false,
			expectedMessage: "Path: $\nExpected: 1\nActual:   nil\n(Should equal)!",
		},
		{
			name:            "Test with nil value2",
			path:            "$",
			value1:          reflect.ValueOf(1),
			value2:          reflect.ValueOf(nil),
			expectedMatch:   false,
			expectedMessage: "Path: $\nExpected: nil\nActual:   1\n(Should equal)!",
		},
		{
			name:            "Test with nil custom assertion",
			path:            "$",
			value1:          reflect.ValueOf(1),
			value2:          reflect.ValueOf(1),
			expectedMatch:   true,
			expectedMessage: "",
		},
		{
			name: "Test with custom assertion",
			path: "$",
			customAssertion: func(actual any, expected ...any) string {
				return "custom assertion"
			},
			value1:          reflect.ValueOf(1),
			value2:          reflect.ValueOf(1),
			expectedMatch:   false,
			expectedMessage: "Path: $\ncustom assertion",
		},
		{
			name:            "Test with different types",
			path:            "$",
			value1:          reflect.ValueOf(1),
			value2:          reflect.ValueOf("1"),
			expectedMatch:   false,
			expectedMessage: "Path: $\nExpected: \"1\"\nActual:   1\n(Should equal)!",
		},
		{
			name:   "Test invalid with custom assertion",
			path:   "$",
			value1: reflect.ValueOf(1),
			value2: reflect.ValueOf(nil),
			customAssertion: func(actual any, expected ...any) string {
				return ""
			},
			expectedMatch:   true,
			expectedMessage: "",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			match, message := assertValue(tt.path, tt.customAssertion, tt.value1, tt.value2)
			if match != tt.expectedMatch {
				t.Errorf("Expected match: %v, got: %v", tt.expectedMatch, match)
			}
			if message != tt.expectedMessage {
				t.Errorf("Expected message:\n%s\ngot:\n%s", tt.expectedMessage, message)
			}
		})
	}
}

func TestGetValue(t *testing.T) {
	testTable := []struct {
		name          string
		value         reflect.Value
		expectedValue any
	}{
		{
			name:          "Test with valid value",
			value:         reflect.ValueOf(1),
			expectedValue: 1,
		},
		{
			name:          "Test with invalid value",
			value:         reflect.ValueOf(nil),
			expectedValue: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			actualValue := getValue(tt.value)
			if actualValue != tt.expectedValue {
				t.Errorf("Expected value: %v, got: %v", tt.expectedValue, actualValue)
			}
		})
	}
}

func TestHasCustomAssertion(t *testing.T) {
	// Define some dummy assertion functions
	assertionFuncA := func(actual any, expected ...any) string { return "assertionFuncA" }
	assertionFuncB := func(actual any, expected ...any) string { return "assertionFuncB" }

	testTable := []struct {
		name             string
		path             string
		fieldType        reflect.Type
		customAssertions map[string]AssertionFunc
		expectedFuncResp string
		expectedOk       bool
	}{
		{
			name:      "Custom assertion by path",
			path:      "someField.subField",
			fieldType: reflect.TypeOf(""),
			customAssertions: map[string]AssertionFunc{
				"someField.subField": assertionFuncA,
			},
			expectedFuncResp: "assertionFuncA",
			expectedOk:       true,
		},
		{
			name:      "Custom assertion by path with index",
			path:      "someField[2].subField",
			fieldType: reflect.TypeOf(""),
			customAssertions: map[string]AssertionFunc{
				"someField[].subField": assertionFuncA,
			},
			expectedFuncResp: "assertionFuncA",
			expectedOk:       true,
		},
		{
			name:      "Custom assertion by type",
			path:      "someField.subField",
			fieldType: reflect.TypeOf(""),
			customAssertions: map[string]AssertionFunc{
				"string": assertionFuncB,
			},
			expectedFuncResp: "assertionFuncB",
			expectedOk:       true,
		},
		{
			name:      "No custom assertion",
			path:      "someField.subField",
			fieldType: reflect.TypeOf(123),
			customAssertions: map[string]AssertionFunc{
				"string": assertionFuncB,
			},
			expectedOk: false,
		},
		{
			name:      "Custom assertion by path and type, path takes precedence",
			path:      "someField[1].subField",
			fieldType: reflect.TypeOf(""),
			customAssertions: map[string]AssertionFunc{
				"someField[].subField": assertionFuncA,
				"string":               assertionFuncB,
			},
			expectedFuncResp: "assertionFuncA",
			expectedOk:       true,
		},
		{
			name: "Custom assertion with nil type",
			path: "$",
			customAssertions: map[string]AssertionFunc{
				"string": assertionFuncB,
			},
			expectedOk: false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			actualFunc, actualOk := hasCustomAssertion(tt.path, tt.fieldType, tt.customAssertions)
			if actualOk != tt.expectedOk {
				t.Errorf("Expected ok: %v, got: %v", tt.expectedOk, actualOk)
			}
			if actualOk && actualFunc(nil, nil) != tt.expectedFuncResp {
				t.Errorf("Expected function response: %s, got: %s", tt.expectedFuncResp, actualFunc(nil, nil))
			}
		})
	}
}

// Test Assert with paths

func TestAssertWithPaths_genericChecks(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")

	testTable := []struct {
		name             string
		actual           any
		expected         any
		customAssertions map[string]AssertionFunc
		expectedMatch    bool
		expectedMessage  string
	}{
		{
			name:          "Test with nil values",
			actual:        nil,
			expected:      nil,
			expectedMatch: true,
		},
		{
			name:            "Test with nil actual",
			actual:          nil,
			expected:        1,
			expectedMatch:   false,
			expectedMessage: "Path: $\nExpected: 1\nActual: <invalid reflect.Value>",
		},
		{
			name:            "Test with nil expected",
			actual:          1,
			expected:        nil,
			expectedMatch:   false,
			expectedMessage: "Path: $\nExpected: <invalid reflect.Value>\nActual: 1",
		},
		{
			name:          "Test with nil pointers",
			actual:        &struct{ Name string }{},
			expected:      &struct{ Name string }{},
			expectedMatch: true,
		},
		{
			name:            "Test with nil actual pointer",
			actual:          &struct{ Name string }{},
			expected:        nil,
			expectedMatch:   false,
			expectedMessage: "Path: $\nExpected: <invalid reflect.Value>\nActual: &{}",
		},
		{
			name:            "Test with nil expected pointer",
			actual:          nil,
			expected:        &struct{ Name string }{},
			expectedMatch:   false,
			expectedMessage: "Path: $\nExpected: &{}\nActual: <invalid reflect.Value>",
		},
		{
			name:     "test with custom assertion on type",
			actual:   time.Now(),
			expected: time.Now(),
			customAssertions: map[string]AssertionFunc{
				TimeType: AssertTimeToDuration(time.Second),
			},
			expectedMatch:   true,
			expectedMessage: "",
		},
		{
			name:     "Test with custom assertion on path",
			actual:   struct{ Time time.Time }{Time: time.Now()},
			expected: struct{ Time time.Time }{Time: time.Now()},
			customAssertions: map[string]AssertionFunc{
				"$.Time": AssertTimeToDuration(time.Second),
			},
			expectedMatch: true,
		},
		{
			name:     "Test with custom assertion on path not matching",
			actual:   struct{ TestValue int }{TestValue: 1},
			expected: struct{ TestValue int }{TestValue: 2},
			customAssertions: map[string]AssertionFunc{
				"$.TestValue": func(actual any, expected ...any) string {
					return "" // always return matching
				},
			},
			expectedMatch: true,
		},
		{
			name: "Test with time and no custom assertion",
			actual: struct {
				Time time.Time
			}{Time: testTime},
			expected: struct {
				Time time.Time
			}{Time: testTime.Add(time.Second)},
			expectedMatch:   false,
			expectedMessage: fmt.Sprintf("Path: $.Time\nExpected: time.Time{%v}\nActual:   time.Time{%v}\n(Should equal)!\n", testTime.Add(time.Second), testTime),
		},
		{
			name: "Test with invalid values and custom assertion",
			actual: struct {
				Time *time.Time
			}{Time: nil},
			expected: struct {
				Time *time.Time
			}{Time: &testTime},
			expectedMatch: true,
			customAssertions: map[string]AssertionFunc{
				"$.Time": func(actual any, expected ...any) string {
					return ""
				},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			match, message := assertWithPaths(reflect.ValueOf(tt.actual), reflect.ValueOf(tt.expected), tt.customAssertions, "$")
			if match != tt.expectedMatch {
				t.Errorf("Expected match: %v, got: %v", tt.expectedMatch, match)
			}
			// remove anything after the word "Diff" in the message till end of line
			re := regexp.MustCompile(`(?m)^.*Diff:.*?(\n|$)`)
			// Replace matched lines with an empty string.
			message = re.ReplaceAllString(message, "")

			if message != tt.expectedMessage {
				t.Errorf("Expected message:\n%s\ngot:\n%s", tt.expectedMessage, message)
			}
		})
	}
}

func TestAssertWithPaths_structs(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	type testStruct struct {
		Name string
		Age  int
	}
	type testNestedStruct struct {
		Sub testStruct
	}
	testTable := []struct {
		name             string
		actual           any
		expected         any
		customAssertions map[string]AssertionFunc
		expectedMatch    bool
		expectedMessage  string
	}{
		{
			name:            "Test with time and no custom assertion",
			actual:          testTime,
			expected:        testTime.Add(time.Second),
			expectedMatch:   false,
			expectedMessage: fmt.Sprintf("Path: $\nExpected: time.Time{%v}\nActual:   time.Time{%v}\n(Should equal)!\n", testTime.Add(time.Second), testTime),
		},
		{
			name:          "Test with nested struct",
			actual:        testNestedStruct{Sub: testStruct{Name: "test", Age: 1}},
			expected:      testNestedStruct{Sub: testStruct{Name: "test", Age: 1}},
			expectedMatch: true,
		},
		{
			name:            "Test with nested struct not matching",
			actual:          testNestedStruct{Sub: testStruct{Name: "test", Age: 1}},
			expected:        testNestedStruct{Sub: testStruct{Name: "test2", Age: 1}},
			expectedMatch:   false,
			expectedMessage: "Path: $.Sub.Name\nExpected: \"test2\"\nActual:   \"test\"\n(Should equal)!\n",
		},
		{
			name:            "Test with nested struct not matching",
			actual:          testNestedStruct{Sub: testStruct{Name: "test", Age: 1}},
			expected:        testNestedStruct{Sub: testStruct{Name: "test", Age: 2}},
			expectedMatch:   false,
			expectedMessage: "Path: $.Sub.Age\nExpected: 2\nActual:   1\n(Should equal)!",
		},
		{
			name:   "Test with nested struct missing field in expected",
			actual: testNestedStruct{Sub: testStruct{Name: "test", Age: 1}},
			expected: struct {
				Sub struct {
					Name string
					Test int
				}
			}{Sub: struct {
				Name string
				Test int
			}{Name: "test"}},
			expectedMatch:   false,
			expectedMessage: "Path: $.Sub.Age\nField Age not found in expected",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			match, message := assertWithPaths(reflect.ValueOf(tt.actual), reflect.ValueOf(tt.expected), tt.customAssertions, "$")
			if match != tt.expectedMatch {
				t.Errorf("Expected match: %v, got: %v", tt.expectedMatch, match)
			}
			// remove anything after the word "Diff" in the message till end of line
			re := regexp.MustCompile(`(?m)^.*Diff:.*?(\n|$)`)
			// Replace matched lines with an empty string.
			message = re.ReplaceAllString(message, "")

			if message != tt.expectedMessage {
				t.Errorf("Expected message:\n%s\ngot:\n%s", tt.expectedMessage, message)
			}
		})
	}
}

func TestAssertWithPaths_Slices(t *testing.T) {

	testTable := []struct {
		name             string
		actual           any
		expected         any
		customAssertions map[string]AssertionFunc
		expectedMatch    bool
		expectedMessage  string
	}{
		{
			name:          "Test with slices matching",
			actual:        []int{1, 2, 3},
			expected:      []int{1, 2, 3},
			expectedMatch: true,
		},
		{
			name: "Test with nested slices matching",
			actual: struct {
				NestedList []int
			}{NestedList: []int{1, 2, 3}},
			expected: struct {
				NestedList []int
			}{NestedList: []int{1, 2, 3}},
			expectedMatch: true,
		},
		{
			name:            "Test with slices not matching in length",
			actual:          []int{1, 2, 3},
			expected:        []int{1, 2},
			expectedMatch:   false,
			expectedMessage: "Path: $\nExpected: []int{1, 2}\nActual:   []int{1, 2, 3}\n(Should equal)!\n",
		},
		{
			name:            "Test with slices not matching in values",
			actual:          []int{1, 2, 3},
			expected:        []int{1, 2, 4},
			expectedMatch:   false,
			expectedMessage: "Path: $[2]\nExpected: 4\nActual:   3\n(Should equal)!",
		},
		{
			name:            "Test with slices not matching in values",
			actual:          []int{0, 2, 4},
			expected:        []int{1, 2, 4},
			expectedMatch:   false,
			expectedMessage: "Path: $[0]\nExpected: 1\nActual:   0\n(Should equal)!",
		},
		{
			name: "Test with nested slices not matching in values with custom assertion",
			actual: struct {
				NestedList []int
			}{NestedList: []int{1, 2, 3}},
			expected: struct {
				NestedList []int
			}{NestedList: []int{1, 2, 4}},
			customAssertions: map[string]AssertionFunc{
				"$.NestedList": func(actual any, expected ...any) string {
					return ""
				},
			},
			expectedMatch: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			match, message := assertWithPaths(reflect.ValueOf(tt.actual), reflect.ValueOf(tt.expected), tt.customAssertions, "$")
			if match != tt.expectedMatch {
				t.Errorf("Expected match: %v, got: %v", tt.expectedMatch, match)
			}
			// remove anything after the word "Diff" in the message till end of line
			re := regexp.MustCompile(`(?m)^.*Diff:.*?(\n|$)`)
			// Replace matched lines with an empty string.
			message = re.ReplaceAllString(message, "")

			if message != tt.expectedMessage {
				t.Errorf("Expected message:\n%s\ngot:\n%s", tt.expectedMessage, message)
			}
		})
	}
}

func TestAssertWithPaths_Maps(t *testing.T) {

	testTable := []struct {
		name             string
		actual           any
		expected         any
		customAssertions map[string]AssertionFunc
		expectedMatch    bool
		expectedMessage  string
	}{
		{
			name:          "Test with maps matching",
			actual:        map[string]int{"a": 1, "b": 2},
			expected:      map[string]int{"a": 1, "b": 2},
			expectedMatch: true,
		},

		{
			name: "Test with nested maps matching",
			actual: struct {
				NestedMap map[string]int
			}{NestedMap: map[string]int{"a": 1, "b": 2}},
			expected: struct {
				NestedMap map[string]int
			}{NestedMap: map[string]int{"a": 1, "b": 2}},
			expectedMatch: true,
		},
		{
			name:            "Test with maps not matching in length",
			actual:          map[string]int{"a": 1, "b": 2},
			expected:        map[string]int{"a": 1},
			expectedMatch:   false,
			expectedMessage: "Path: $\nExpected: map[string]int{\"a\":1}\nActual:   map[string]int{\"a\":1, \"b\":2}\n(Should equal)!\n",
		},
		{
			name:            "Test with maps not matching in values",
			actual:          map[string]int{"a": 1, "b": 2},
			expected:        map[string]int{"a": 1, "b": 3},
			expectedMatch:   false,
			expectedMessage: "Path: $.b\nExpected: 3\nActual:   2\n(Should equal)!",
		},
		{
			name:            "Test with maps not matching in values",
			actual:          map[string]int{"a": 0, "b": 3},
			expected:        map[string]int{"a": 1, "b": 3},
			expectedMatch:   false,
			expectedMessage: "Path: $.a\nExpected: 1\nActual:   0\n(Should equal)!",
		},
		{
			name: "Test with maps not matching in values with custom assertion",
			actual: struct {
				NestedMap map[string]int
			}{NestedMap: map[string]int{"a": 1, "b": 2}},
			expected: struct {
				NestedMap map[string]int
			}{NestedMap: map[string]int{"a": 1, "b": 3}},
			customAssertions: map[string]AssertionFunc{
				"$.NestedMap": func(actual any, expected ...any) string {
					return ""
				},
			},
			expectedMatch: true,
		},
		{
			name:            "Test with map not containing key",
			actual:          map[string]int{"a": 1, "b": 2},
			expected:        map[string]int{"a": 1, "c": 2},
			expectedMatch:   false,
			expectedMessage: "Path: $.b\nKey b not found in expected",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			match, message := assertWithPaths(reflect.ValueOf(tt.actual), reflect.ValueOf(tt.expected), tt.customAssertions, "$")
			if match != tt.expectedMatch {
				t.Errorf("Expected match: %v, got: %v", tt.expectedMatch, match)
			}
			// remove anything after the word "Diff" in the message till end of line
			re := regexp.MustCompile(`(?m)^.*Diff:.*?(\n|$)`)
			// Replace matched lines with an empty string.
			message = re.ReplaceAllString(message, "")

			if message != tt.expectedMessage {
				t.Errorf("Expected message:\n%s\ngot:\n%s", tt.expectedMessage, message)
			}
		})
	}
}

func TestAssertWithPaths_Complex(t *testing.T) {
	type subStruct struct {
		Field1 string
		Field2 []int
		Field3 map[string]int
	}
	type testStruct struct {
		Field1 string
		Field2 int
		Field3 []string
		Field4 map[string]string
		Field5 subStruct
	}

	testTable := []struct {
		name             string
		actual           any
		expected         any
		customAssertions map[string]AssertionFunc
		expectedMatch    bool
		expectedMessage  string
	}{
		{
			name: "Test with struct",
			actual: testStruct{
				Field1: "test",
				Field2: 1,
				Field3: []string{"a", "b"},
				Field4: map[string]string{"a": "1", "b": "2"},
				Field5: subStruct{
					Field1: "sub",
					Field2: []int{1, 2},
					Field3: map[string]int{"a": 1, "b": 2},
				},
			},
			expected: testStruct{
				Field1: "test",
				Field2: 1,
				Field3: []string{"a", "b"},
				Field4: map[string]string{"a": "1", "b": "2"},
				Field5: subStruct{
					Field1: "sub",
					Field2: []int{1, 2},
					Field3: map[string]int{"a": 1, "b": 2},
				},
			},
			expectedMatch: true,
		},
		{
			name: "Test with struct not matching",
			actual: testStruct{
				Field1: "test",
				Field2: 1,
				Field3: []string{"a", "b"},
				Field4: map[string]string{"a": "1", "b": "2"},
				Field5: subStruct{
					Field1: "sub",
					Field2: []int{1, 2},
					Field3: map[string]int{"a": 1, "b": 2},
				},
			},
			expected: testStruct{
				Field1: "test",
				Field2: 1,
				Field3: []string{"a", "b"},
				Field4: map[string]string{"a": "1", "b": "2"},
				Field5: subStruct{
					Field1: "sub",
					Field2: []int{1, 2},
					Field3: map[string]int{"a": 1, "b": 3},
				},
			},
			expectedMatch:   false,
			expectedMessage: "Path: $.Field5.Field3.b\nExpected: 3\nActual:   2\n(Should equal)!",
		},
		{
			name: "Test with struct not matching multiple values",
			actual: testStruct{
				Field1: "test",
				Field2: 1,
				Field3: []string{"a", "b"},
				Field4: map[string]string{"a": "1", "b": "2"},
				Field5: subStruct{
					Field1: "sub",
					Field2: []int{1, 2, 3},
					Field3: map[string]int{"a": 1, "b": 2},
				},
			},
			expected: testStruct{
				Field1: "test5",
				Field2: 5,
				Field3: []string{"a", "b"},
				Field4: map[string]string{"a": "1", "b": "2"},
				Field5: subStruct{
					Field1: "sub",
					Field2: []int{1, 2},
					Field3: map[string]int{"a": 1, "b": 3},
				},
			},
			expectedMatch:   false,
			expectedMessage: "Path: $.Field1\nExpected: \"test5\"\nActual:   \"test\"\n(Should equal)!\nPath: $.Field2\nExpected: 5\nActual:   1\n(Should equal)!\nPath: $.Field5.Field2\nExpected: []int{1, 2}\nActual:   []int{1, 2, 3}\n(Should equal)!\nPath: $.Field5.Field3.b\nExpected: 3\nActual:   2\n(Should equal)!",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			match, message := assertWithPaths(reflect.ValueOf(tt.actual), reflect.ValueOf(tt.expected), tt.customAssertions, "$")
			if match != tt.expectedMatch {
				t.Errorf("Expected match: %v, got: %v", tt.expectedMatch, match)
			}
			// remove anything after the word "Diff" in the message till end of line
			re := regexp.MustCompile(`(?m)^.*Diff:.*?(\n|$)`)
			// Replace matched lines with an empty string.
			message = re.ReplaceAllString(message, "")

			if message != tt.expectedMessage {
				t.Errorf("Expected message:\n%s\ngot:\n%s", tt.expectedMessage, message)
			}
		})
	}
}

func TestFormatMessage(t *testing.T) {
	testTable := []struct {
		name            string
		message         string
		format          string
		a               []any
		expectedMessage string
	}{
		{
			name:            "Test with message",
			message:         "message",
			format:          "format %v",
			a:               []any{"arg"},
			expectedMessage: "message\nformat arg",
		},
		{
			name:            "Test with empty message",
			format:          "format %v",
			a:               []any{"arg"},
			expectedMessage: "format arg",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			actualMessage := formatMessage(tt.message, tt.format, tt.a...)
			if actualMessage != tt.expectedMessage {
				t.Errorf("Expected message: %s, got: %s", tt.expectedMessage, actualMessage)
			}
		})
	}
}

func TestGetType(t *testing.T) {
	testTable := []struct {
		name         string
		actual       reflect.Value
		expected     reflect.Value
		expectedType reflect.Type
	}{
		{
			name:         "Test with valid values",
			actual:       reflect.ValueOf(1),
			expected:     reflect.ValueOf(1),
			expectedType: reflect.TypeOf(1),
		},
		{
			name:         "Test with invalid values",
			actual:       reflect.ValueOf(nil),
			expected:     reflect.ValueOf(nil),
			expectedType: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			actualType := getType(tt.actual, tt.expected)
			if actualType != tt.expectedType {
				t.Errorf("Expected type: %v, got: %v", tt.expectedType, actualType)
			}
		})
	}

}
