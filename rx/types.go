package rx

import "context"

type (
	startStrategy int

	// Predicate defines a func that returns a bool from an input value.
	Predicate func(ctx context.Context, v interface{}) bool

	// Mapper defines a function that computes a value from an input value.
	Mapper func(ctx context.Context, i interface{}) interface{}

	// Comparator defines a func that returns an int:
	// - 0 if two elements are equals
	// - A negative value if the first argument is less than the second
	// - A positive value if the first argument is greater than the second
	Comparator func(ctx context.Context, a interface{}, b interface{}) int
)

// List of known starting strategies
const (
	Lazy  startStrategy = 0
	Eager startStrategy = 1
)
