package rx

import (
	"context"
)

type operatorAll[T comparable] struct {
	ctx       context.Context
	predicate Predicate[T]
	all       bool
}

func (o *operatorAll[T]) next(item T, dst chan<- T) bool {
	if !o.predicate(o.ctx, item) {
		o.all = false
	}

	return false
}

func (o *operatorAll[T]) end(dst chan<- T) {
	//send(dst, o.all)
}

// All determine whether all items emitted meet some criteria.
func (o *Operable[T]) All(predicate Predicate[T]) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorAll[T]{
		ctx:       o.ctx,
		predicate: predicate,
		all:       true,
	})

	return o
}
