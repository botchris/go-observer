package rx

type (
	startStrategy int

	// Predicate defines a func that returns a bool from an input value.
	Predicate[T any] func(v T) bool

	// Mapper defines a function that computes a new value from an input value.
	Mapper[I any, O any] func(input I) O

	// Comparator defines a func that returns an int:
	// - 0 if two elements are equals
	// - A negative value if the first argument is less than the second
	// - A positive value if the first argument is greater than the second
	Comparator[T comparable] func(a T, b T) int
)
