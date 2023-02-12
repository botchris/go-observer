package rx

import (
	"context"
)

type operatorMax[T comparable] struct {
	ctx        context.Context
	comparator Comparator[T]
	empty      bool
	max        T
}

func (o *operatorMax[T]) next(item T, dst chan<- T) bool {
	o.empty = false
	if o.max == nil {
		o.max = item

		return false
	}

	if o.comparator(o.ctx, o.max, item) < 0 {
		o.max = item
	}

	return false
}

func (o *operatorMax[T]) end(dst chan<- T) {
	if !o.empty {
		send(dst, o.max)
	}
}

// Max determines and emits the maximum-valued item according to a comparator.
func (o *Operable[T]) Max(comparator Comparator[T]) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorMax[T]{
		ctx:        o.ctx,
		comparator: comparator,
		empty:      true,
	})

	return o
}
