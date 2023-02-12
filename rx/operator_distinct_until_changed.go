package rx

import (
	"context"
)

type operatorDistinctUntilChanged struct {
	ctx     context.Context
	apply   Mapper
	current interface{}
}

func (o *operatorDistinctUntilChanged) next(item interface{}, dst chan<- interface{}) bool {
	key := o.apply(o.ctx, item)

	if o.current != key {
		o.current = key
		send(dst, item)

		return true
	}

	return false
}

func (o *operatorDistinctUntilChanged) end(dst chan<- interface{}) {}

// DistinctUntilChanged suppresses consecutive duplicate items.
func (o *Operable[T]) DistinctUntilChanged(apply Mapper) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorDistinctUntilChanged{
		ctx:   o.ctx,
		apply: apply,
	})

	return o
}
