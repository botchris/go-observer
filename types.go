package observer

type startStrategy int

// List of known starting strategies
const (
	Lazy  startStrategy = 0
	Eager startStrategy = 1
)

// Predicate defines a func that returns a bool from an input value.
type Predicate func(v interface{}) bool

// Mapper defines a function that computes a value from an input value.
type Mapper func(i interface{}) interface{}
