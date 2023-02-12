package rx

import (
	"context"
)

type operatorAll struct {
	ctx       context.Context
	predicate Predicate
	all       bool
}

func (o *operatorAll) next(item interface{}, dst chan<- interface{}) bool {
	if !o.predicate(o.ctx, item) {
		o.all = false
	}

	return false
}

func (o *operatorAll) end(dst chan<- interface{}) {
	send(dst, o.all)
}

// All determine whether all items emitted meet some criteria.
func (o *Operable[T]) All(predicate Predicate) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorAll{
		ctx:       o.ctx,
		predicate: predicate,
		all:       true,
	})

	return o
}
