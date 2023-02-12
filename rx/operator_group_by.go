package rx

import (
	"github.com/botchris/observer"
)

// GroupBy divides an Operable into a set of Operable that each emit a different group of items from the original Operable, organized by key.
func (o *Operable[T]) GroupBy(length int, distribution func(item T) int) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	var def T

	root := observer.NewProperty[T](def)
	fork := MakeOperable[T](o.ctx, root.Observe())
	properties := make([]observer.Property[T], length)

	for i := 0; i < length; i++ {
		p := observer.NewProperty[T](def)
		properties[i] = p
		root.Update(MakeOperable[T](o.ctx, p.Observe()))
	}

	root.End()
	ready := make(chan struct{})
	done := o.ctx.Done()

	go func() {
		defer func() {
			for i := 0; i < length; i++ {
				properties[i].End()
			}
		}()

		close(ready)

		for {
			select {
			case <-done:
				return
			case <-o.Changes():
				value := o.Next()
				if value == o.eof {
					return
				}

				idx := distribution(value)
				if idx >= length {
					// TODO: handle error
					return
				}

				properties[idx].Update(value)
			}
		}
	}()

	<-ready

	return fork
}
