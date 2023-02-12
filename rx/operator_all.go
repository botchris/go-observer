package rx

import (
	"github.com/botchris/observer"
)

type allOperator[T any] struct {
	observer.Property[T]
	predicate Predicate[T]
	allMeet   bool
}

type AllProperty[T any] interface {
	AllMeet() bool
	observer.Property[T]
}

func (o *allOperator[T]) AllMeet() bool {
	return o.allMeet
}

// Update sets a new value for this property.
func (o *allOperator[T]) Update(value ...T) {
	for _, v := range value {
		if o.allMeet && !o.predicate(v) {
			o.allMeet = false
		}
	}
}

func All[T any](pro observer.Property[T], pre Predicate[T]) AllProperty[T] {
	return &allOperator[T]{
		Property:  pro,
		predicate: pre,
		allMeet:   true,
	}
}
