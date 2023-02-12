package rx

import (
	"context"
)

type operatorDistinct struct {
	ctx    context.Context
	apply  Mapper
	keyset map[interface{}]struct{}
}

func (o *operatorDistinct) next(item interface{}, dst chan<- interface{}) bool {
	key := o.apply(o.ctx, item)

	if _, exists := o.keyset[key]; !exists {
		send(dst, item)
		o.keyset[key] = struct{}{}

		return true
	}

	return false
}

func (o *operatorDistinct) end(dst chan<- interface{}) {}

// Distinct suppresses duplicate items.
func (o *Operable[T]) Distinct(apply Mapper) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorDistinct{
		ctx:    o.ctx,
		apply:  apply,
		keyset: make(map[interface{}]struct{}),
	})

	return o
}
