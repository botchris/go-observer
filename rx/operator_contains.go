package rx

import (
	"context"
)

type operatorContains struct {
	ctx       context.Context
	predicate Predicate
	contains  bool
	stopped   bool
}

func (o *operatorContains) next(item interface{}, dst chan<- interface{}) bool {
	if !o.stopped && o.predicate(o.ctx, item) {
		o.contains = true
		o.stopped = true
		send(dst, true)

		return true
	}

	return false
}

func (o *operatorContains) end(dst chan<- interface{}) {
	if !o.stopped && !o.contains {
		send(dst, false)
	}
}

// Contains determine whether a particular item was emitted or not.
func (o *Operable[T]) Contains(predicate Predicate) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorContains{
		ctx:       o.ctx,
		predicate: predicate,
	})

	return o
}
