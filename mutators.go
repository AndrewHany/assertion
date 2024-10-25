package assertion

import (
	"sort"
)

type MutatorFunc func(actual any, expected any) (any, any)

func AssertUnorderedSlice[T any](compareFunc func(actual, expected T) bool) MutatorFunc {
	return func(actual any, expected any) (any, any) {

		if expected == nil {
			return actual, expected
		}

		if act, ok := actual.([]T); ok {
			if exp, ok := expected.([]T); ok {
				actual := make([]T, len(act))
				expected := make([]T, len(exp))
				copy(actual, act)
				copy(expected, exp)

				sort.Slice(actual, func(i, j int) bool {
					return compareFunc(actual[i], actual[j])
				})
				sort.Slice(expected, func(i, j int) bool {
					return compareFunc(expected[i], expected[j])
				})
				return actual, expected
			}
		}
		return actual, expected
	}
}
