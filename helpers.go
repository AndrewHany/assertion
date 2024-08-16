package assertion

func IsNilExpected(actual any, expected any) bool {
	return expected == nil
}

func IsNilActual(actual any, expected any) bool {
	return actual == nil
}
