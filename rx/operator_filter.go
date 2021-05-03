package rx

import (
	"context"
)

type operatorFilter struct {
	ctx       context.Context
	predicate Predicate
}

func (o *operatorFilter) next(item interface{}, dst chan<- interface{}) bool {
	if o.predicate(o.ctx, item) {
		send(dst, item)

		return true
	}

	return false
}

func (o *operatorFilter) end(dst chan<- interface{}) {}

// Filter emit only those items that pass a predicate test.
func (o *Operable) Filter(predicate Predicate) *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorFilter{
		ctx:       o.ctx,
		predicate: predicate,
	})

	return o
}
