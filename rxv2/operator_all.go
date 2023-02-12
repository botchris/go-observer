package rx

import (
	"context"
	"sync"

	"github.com/botchris/observer"
)

type AllStream[T any] interface {
	AllMeet() bool
	observer.Stream[T]
}

type all[T any] struct {
	ctx       context.Context
	pipe      observer.Property[T]
	predicate Predicate[T]

	source observer.Stream[T]
	observer.Stream[T]

	allMeet bool
	mu      sync.RWMutex
}

func All[T any](ctx context.Context, source observer.Stream[T], pre Predicate[T]) AllStream[T] {
	op := &all[T]{
		ctx:       ctx,
		source:    source.Clone(),
		pipe:      observer.NewProperty[T](source.Value()),
		predicate: pre,
		allMeet:   true,
	}

	defer op.init()

	return op
}

func (o *all[T]) AllMeet() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.allMeet
}

func (o *all[T]) init() {
	o.Stream = o.pipe.Observe()

	go func() {
		for {
			select {
			case <-o.ctx.Done():
				return
			case <-o.source.Changes():
				o.mu.Lock()

				if o.allMeet && !o.predicate(o.source.Value()) {
					o.allMeet = false
				}

				o.mu.Unlock()

				o.pipe.Update(o.source.Value())
			}
		}
	}()
}
