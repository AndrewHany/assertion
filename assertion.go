package assertion

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/smarty/assertions"
)

var defaultAssertionFunc = assertions.ShouldEqual
var removeIndexRegex = regexp.MustCompile(`\[\d+\]`)

// Assert compares the actual and expected values using the custom assertions defined
// and returns the result and message
// actual is the actual value to be compared
// expected is the expected value to be compared
// customAssertions is the map of custom assertions defined for the path or type
// Example usage:
//
//	customAssertions := map[string]AssertionFunc{
//		"field1": customAssertionFunc1,
//		"field2": customAssertionFunc2,
//		"int":    customAssertionFunc3,
//	}
//
// match, message := Assert(actual, expected, customAssertions)
// returns the result and message
func Assert(actual any, expected any, customAssertions map[string]AssertionFunc) (bool, string) {
	// If not custom assertion defined, use default assertion for the whole object
	return assertWithPaths(reflect.ValueOf(actual), reflect.ValueOf(expected), customAssertions, "$")
}

// assertWithPaths recursively compares the actual and expected values
// using the custom assertions defined for the path or type of the field
// and returns the result and message
// path is the path of the field in the object
// customAssertions is the map of custom assertions defined for the path or type
// actual is the actual value to be compared
// expected is the expected value to be compared
// returns the result and message
func assertWithPaths(
	actual reflect.Value,
	expected reflect.Value,
	customAssertions map[string]AssertionFunc,
	path string,
) (bool, string) {
	match, message := true, ""
	// Dereference pointers
	for actual.Kind() == reflect.Ptr && expected.Kind() == reflect.Ptr {
		actual = actual.Elem()
		expected = expected.Elem()
	}

	// handle nil pointers
	if !actual.IsValid() && !expected.IsValid() {
		return true, ""
	}
	typ := getType(actual, expected)

	// check if custom assertion is defined for the path
	if customAssertionFunc, ok := hasCustomAssertion(path, typ, customAssertions); ok {
		return assertValue(path, customAssertionFunc, actual, expected)
	}

	if !actual.IsValid() || !expected.IsValid() {
		return false, formatMessage(message, "Path: %s\nExpected: %v\nActual: %v", path, expected, actual)
	}

	switch actual.Kind() {
	case reflect.Struct:
		// handle time.Time
		if actual.Type() == reflect.TypeOf(time.Time{}) {
			return assertValue(path, nil, actual, expected)
		}
		// handle structs not matching in fields
		if actual.NumField() != expected.NumField() {
			return assertValue(path, defaultAssertionFunc, actual, expected)
		}

		for i := 0; i < actual.NumField(); i++ {
			field := actual.Type().Field(i)
			fieldPath := path + "." + field.Name
			// check if expected has the same field
			if !expected.FieldByName(field.Name).IsValid() {
				return false, formatMessage(message, "Path: %s\nField %s not found in expected", fieldPath, field.Name)
			}
			if listMatch, listMessage := assertWithPaths(actual.Field(i), expected.FieldByName(field.Name), customAssertions, fieldPath); !listMatch {
				match = false
				message = formatMessage(message, "%s", listMessage)
			}
		}
	case reflect.Slice, reflect.Array:
		if actual.Len() != expected.Len() {
			return assertValue(path, defaultAssertionFunc, actual, expected)
		}
		for i := 0; i < actual.Len(); i++ {
			if listMatch, listMessage := assertWithPaths(actual.Index(i), expected.Index(i), customAssertions, fmt.Sprintf("%s[%d]", path, i)); !listMatch {
				match = false
				message = formatMessage(message, "%s", listMessage)
			}
		}
	case reflect.Map:
		if actual.Len() != expected.Len() {
			return assertValue(path, defaultAssertionFunc, actual, expected)
		}
		for _, key := range actual.MapKeys() {
			// check if expected has the same key
			if !expected.MapIndex(key).IsValid() {
				return false, formatMessage(message, "Path: %s\nKey %s not found in expected", path+"."+key.String(), key.String())
			}
			if listMatch, listMessage := assertWithPaths(actual.MapIndex(key), expected.MapIndex(key), customAssertions, fmt.Sprintf("%s.%v", path, key.Interface())); !listMatch {
				match = false
				message = formatMessage(message, "%s", listMessage)
			}
		}
	default:
		// check for custom assertions with path
		return assertValue(path, defaultAssertionFunc, actual, expected)
	}
	return match, message
}

// hasCustomAssertion checks if custom assertion is defined for the path or type of the field
func hasCustomAssertion(path string, fieldType reflect.Type, customAssertions map[string]AssertionFunc) (AssertionFunc, bool) {
	// check if custom assertion is defined for the path
	// replace index with [] to match the path
	if customAssertionByPath, ok := customAssertions[removeIndexRegex.ReplaceAllString(path, "[]")]; ok {
		return customAssertionByPath, true
	}
	// check if custom assertion is defined for the type
	if fieldType != nil {
		if customAssertionByType, ok := customAssertions[fieldType.String()]; ok {
			return customAssertionByType, true
		}
	}
	return nil, false
}

// assertValue checks if the actual and expected values are matching
// using the custom assertion function or default assertion function shouldEqual
// and returns the result and message
func assertValue(path string, customAssertion AssertionFunc, actual reflect.Value, expected reflect.Value) (bool, string) {
	// if both are nil, return true
	if customAssertion == nil {
		customAssertion = defaultAssertionFunc
	}

	isMatching, newMessage := assertions.So(getValue(actual), assertions.SoFunc(customAssertion), getValue(expected))

	if !isMatching {
		return false, fmt.Sprintf("Path: %s\n%s", path, newMessage)
	}
	return true, ""
}

// getValue returns the interface value of the reflect value or nil if not valid
func getValue(value reflect.Value) any {
	if !value.IsValid() {
		return nil
	}
	return value.Interface()
}

// getType returns the type of the actual value if valid, else returns the type of the expected value if valid else nil
func getType(act reflect.Value, exp reflect.Value) reflect.Type {
	if act.IsValid() {
		return act.Type()
	}
	if exp.IsValid() {
		return exp.Type()
	}
	return nil
}

// formatMessage formats the message with the format and arguments using fmt.Sprintf
func formatMessage(message string, format string, a ...any) string {
	if message == "" {
		return fmt.Sprintf(format, a...)
	}
	return fmt.Sprintf("%s\n%s", message, fmt.Sprintf(format, a...))
}
