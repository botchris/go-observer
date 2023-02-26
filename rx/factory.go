package rx

import (
	"context"
	"io"

	"github.com/botchris/go-observer"
)

// MakeOperable makes the given input Stream operable so operators can be applied to. This new operable instance will be
// valid as long as the given context keeps active.
//
// By default operable streams are created in "Lazy" mode, which means they won't start reading the input Stream until
// the operable itself is "consumed" (calls to any Stream interface method). This is specially useful when
// chaining multiple operators at once, so your "operators pipeline" is correctly defined upfront.
func MakeOperable(ctx context.Context, input observer.Stream, opts ...Option) *Operable {
	p := observer.NewProperty(nil)
	o := &Operable{
		ctx:   ctx,
		input: input.Clone(),

		done:      make(chan struct{}),
		output:    p.Observe(),
		operators: make([]operator, 0),
		surrogate: p,
	}

	options := &options{
		startStrategy: Lazy,
	}

	for _, o := range opts {
		o.apply(options)
	}

	if options.startStrategy == Eager {
		o.Start()
	}

	return o
}

// Concat emit the emissions from two or more source streams without interleaving them.
func Concat(ctx context.Context, sources []observer.Stream, opts ...Option) *Operable {
	p := observer.NewProperty(nil)
	s := p.Observe()
	ready := make(chan struct{})

	done := ctx.Done()
	go func() {
		defer p.Update(io.EOF)

		inputs := make([]observer.Stream, 0)
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
					if v == io.EOF {
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
