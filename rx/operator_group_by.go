package rx

import (
	"io"

	"github.com/botchris/go-observer"
)

// GroupBy divides an Operable into a set of Operable that each emit a different group of items from the original Operable, organized by key.
func (o *Operable) GroupBy(length int, distribution func(item interface{}) int) *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	root := observer.NewProperty(nil)
	fork := MakeOperable(o.ctx, root.Observe())
	properties := make([]observer.Property, length)

	for i := 0; i < length; i++ {
		p := observer.NewProperty(nil)
		properties[i] = p
		root.Update(MakeOperable(o.ctx, p.Observe()))
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
				if value == io.EOF {
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
