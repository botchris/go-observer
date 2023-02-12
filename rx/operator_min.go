package rx

import (
	"context"
)

type operatorMin[T comparable] struct {
	ctx        context.Context
	comparator Comparator[T]
	empty      bool
	max        T
}

func (o *operatorMin[T]) next(item T, dst chan<- T) bool {
	o.empty = false
	if o.max == nil {
		o.max = item

		return false
	}

	if o.comparator(o.ctx, o.max, item) > 0 {
		o.max = item
	}

	return false
}

func (o *operatorMin[T]) end(dst chan<- T) {
	if !o.empty {
		send(dst, o.max)
	}
}

// Min determines and emits the minimum-valued item according to a comparator.
func (o *Operable[T]) Min(comparator Comparator[T]) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorMin[T]{
		ctx:        o.ctx,
		comparator: comparator,
		empty:      true,
	})

	return o
}
