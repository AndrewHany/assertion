# assertions
Assertion style package, supporting custom assertions

This package is build on top of `github.com/smarty/assertions` package
it allows custom assertions on fields and types
you basically do
```go
	customAssertions := map[string]assertions.AssertionFunc{
		"$.Field1.Field2":   assertions.SkipAssertion,
		assertion.TimeType:  assertions.AssertTimeToDuration(time.Second),
		assertion.FloatType: assertions.AssertFloat64ToDecimalPlaces(2),
		"$.Field1[][]":      customAssertionFunc,
	}
```
and pass this assertion map to the assert function
For example
- `"$.Field1.Field2": assertions.SkipAssertion,`
references nested fields in root struct, and skips assertion
- `assertion.TimeType: assertions.AssertTimeToDuration(time.Second)`
asserts till seconds (skip milliesconds)
- `assertion.FloatType: assertions.AssertFloat64ToDecimalPlaces(2)`
asserts till the first 2 decimal places
- `"$.Field1[][]": customAssertionFunc`
Or you can build your custom assertion method
