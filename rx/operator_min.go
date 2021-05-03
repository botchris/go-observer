package rx

import (
	"context"
)

type operatorMin struct {
	ctx        context.Context
	comparator Comparator
	empty      bool
	max        interface{}
}

func (o *operatorMin) next(item interface{}, dst chan<- interface{}) bool {
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

func (o *operatorMin) end(dst chan<- interface{}) {
	if !o.empty {
		send(dst, o.max)
	}
}

// Min determines and emits the minimum-valued item according to a comparator.
func (o *Operable) Min(comparator Comparator) *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorMin{
		ctx:        o.ctx,
		comparator: comparator,
		empty:      true,
	})

	return o
}
