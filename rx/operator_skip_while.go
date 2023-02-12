package rx

import (
	"context"
)

type operatorSkipWhile[T comparable] struct {
	ctx       context.Context
	predicate Predicate[T]
	skip      bool
}

func (o *operatorSkipWhile[T]) next(item T, dst chan<- T) bool {
	if !o.skip {
		send(dst, item)

		return true
	}

	if !o.predicate(o.ctx, item) {
		o.skip = false
		send(dst, item)

		return true
	}

	return false
}

func (o *operatorSkipWhile[T]) end(dst chan<- T) {}

// SkipWhile discard items until a specified condition becomes false.
func (o *Operable[T]) SkipWhile(predicate Predicate[T]) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorSkipWhile[T]{
		ctx:       o.ctx,
		predicate: predicate,
		skip:      true,
	})

	return o
}
