package rx

import (
	"context"
)

type operatorMap[T any] struct {
	ctx    context.Context
	mapper Mapper[T]
}

func (o *operatorMap[T]) next(item T, dst chan<- T) bool {
	send(dst, o.mapper(o.ctx, item))

	return true
}

func (o *operatorMap[T]) end(dst chan<- T) {}

// Map transform the items by applying a function to each item.
func (o *Operable[T]) Map(mapper Mapper[T]) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorMap[T]{
		ctx:    o.ctx,
		mapper: mapper,
	})

	return o
}
