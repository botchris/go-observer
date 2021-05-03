package rx

import (
	"context"
)

type operatorMap struct {
	ctx    context.Context
	mapper Mapper
}

func (o *operatorMap) next(item interface{}, dst chan<- interface{}) bool {
	send(dst, o.mapper(o.ctx, item))

	return true
}

func (o *operatorMap) end(dst chan<- interface{}) {}

// Map transform the items by applying a function to each item.
func (o *Operable) Map(mapper Mapper) *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorMap{
		ctx:    o.ctx,
		mapper: mapper,
	})

	return o
}
