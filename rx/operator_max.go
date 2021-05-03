package rx

import (
	"context"
)

type operatorMax struct {
	ctx        context.Context
	comparator Comparator
	empty      bool
	max        interface{}
}

func (o *operatorMax) next(item interface{}, dst chan<- interface{}) bool {
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

func (o *operatorMax) end(dst chan<- interface{}) {
	if !o.empty {
		send(dst, o.max)
	}
}

// Max determines and emits the maximum-valued item according to a comparator.
func (o *Operable) Max(comparator Comparator) *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorMax{
		ctx:        o.ctx,
		comparator: comparator,
		empty:      true,
	})

	return o
}
