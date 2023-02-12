package rx

import (
	"context"

	"github.com/botchris/observer"
)

// MakeOperable makes the given input Stream operable so operators can be applied to. This new operable instance will be
// valid as long as the given context keeps active.
//
// By default, operable streams are created in "Lazy" mode, which means they won't start reading the input Stream until
// the operable itself is "consumed" (calls to any Stream interface method). This is specially useful when
// chaining multiple operators at once, so your "operators pipeline" is correctly defined upfront.
func MakeOperable[T comparable](ctx context.Context, input observer.Stream[T], opts ...Option) *Operable[T] {
	var eof T

	p := observer.NewProperty[T](eof)
	o := &Operable[T]{
		ctx:   ctx,
		input: input.Clone(),
		eof:   eof,

		done:      make(chan struct{}),
		output:    p.Observe(),
		operators: make([]operator[T], 0),
		surrogate: p,
	}

	ops := &options{
		startStrategy: Lazy,
	}

	for _, o := range opts {
		o.apply(ops)
	}

	if ops.startStrategy == Eager {
		o.Start()
	}

	return o
}

// Concat emit the emissions from two or more source streams without interleaving them.
func Concat[T comparable](ctx context.Context, sources []observer.Stream[T], opts ...Option) *Operable[T] {
	var eof T

	p := observer.NewProperty(eof)
	s := p.Observe()
	ready := make(chan struct{})
	done := ctx.Done()

	go func() {
		defer p.Update(eof)

		inputs := make([]observer.Stream[T], 0)
		for _, in := range sources {
			inputs = append(inputs, in.Clone())
		}

		close(ready)

		for _, src := range inputs {
		loop:
			for {
				select {
				case <-done:
					return
				case <-src.Changes():
					v := src.Next()
					if v == eof {
						break loop
					}

					p.Update(v)
				}
			}
		}
	}()

	<-ready

	return MakeOperable(ctx, s, opts...)
}
