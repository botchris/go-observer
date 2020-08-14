package observer

import (
	"context"
	"io"
)

// Concat emit the emissions from two or more source streams without interleaving them.
func Concat(ctx context.Context, sources []Stream, opts ...Option) *Operable {
	p := NewProperty(nil)
	s := p.Observe()
	ready := make(chan struct{})

	go func() {
		defer p.Update(io.EOF)

		inputs := make([]Stream, 0)
		for _, in := range sources {
			inputs = append(inputs, in.Clone())
		}

		close(ready)

		for _, src := range inputs {
		loop:
			for {
				select {
				case <-ctx.Done():
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
