package rx

import (
	"context"
	"sync"

	"github.com/botchris/observer"
)

// Operable defines a wrapped stream (input) on which operators can be applied to.
// Operable Streams live as long as the underlying context remains active. If context ends or gets
// cancelled operable will emit a io.EOF and no further items will be emitted.
type Operable[T any] struct {
	observer.Stream[T]
	ctx   context.Context
	input observer.Stream[T]
	eof   T

	mu         sync.RWMutex
	running    bool
	completed  bool
	done       chan struct{}
	output     observer.Stream[T]
	operators  []operator[T]
	surrogate  observer.Property[T]
	onStart    func()
	onNext     func(interface{})
	onComplete func()
}

// operator represents a component capable to altering/transforming the flow of items emitted by a source Stream.
// They can be stacked and are executed sequentially, so the changes made by one operator are visible to the
// operators coming after.
type operator[T any] interface {
	// next is triggered when the "source" Stream emits a new item.
	// Retuning FALSE indicates that no further operators should be applied, AND NOTHING should be written to the
	// output Stream; basically "discards" the item producing no visible changes at the output stream.
	//
	// - item: is the value emitted by the source Stream, or emitted by the previous operator
	// - dst: this channel is used to indicate that the operator wants to change the given "item" value
	//   for another value written on the channel.
	next(item T, dst chan<- T) (next bool)

	// end is invoked when "source" Stream reaches io.EOF.
	// The given channel is used by the operator to indicate it wants to write something to the
	// output stream.
	end(dst chan<- T)
}

// send sends the given item over the given channel in a non-blocking fashion
func send[T any](dst chan<- T, item T) {
	select {
	case dst <- item:
	default:
	}
}

// rcv receives a value from the given source or returns the given default value if there is nothing to read
func rcv[T any](src <-chan T, defaultValue T) T {
	select {
	case v := <-src:
		return v
	default:
	}

	return defaultValue
}

// Start starts reading the input stream, it will no-op if already started.
func (o *Operable[T]) Start() *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.running {
		return o
	}

	ready := make(chan struct{})
	go o.run(ready)

	<-ready
	o.running = true

	if o.onStart != nil {
		o.onStart()
	}

	return o
}

// OnStart registers a callback action that will be called once the Operable starts reading the input Stream.
func (o *Operable[T]) OnStart(startFunc func()) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.onStart = startFunc

	return o
}

// OnComplete registers a callback action that will be called after Next is invoked and reaches a io.EOF state.
func (o *Operable[T]) OnComplete(completedFunc func()) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.onComplete = completedFunc

	return o
}

// OnNext registers a callback action that will be called each time Next method is invoked.
func (o *Operable[T]) OnNext(nextFunc func(interface{})) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.onNext = nextFunc

	return o
}

// Value returns the current value for this stream.
func (o *Operable[T]) Value() T {
	o.Start()

	return o.output.Value()
}

// Changes returns the channel that is closed when a new value is available.
func (o *Operable[T]) Changes() chan struct{} {
	o.Start()

	return o.output.Changes()
}

// Next advances this stream to the next state.
// You should never call this unless Changes channel is closed.
func (o *Operable[T]) Next() T {
	o.Start()

	value := o.output.Next()
	if o.onNext != nil && value != o.eof {
		o.onNext(value)
	}

	if value == o.eof {
		o.complete()
	}

	return value
}

// Done returns a channel that's closed when Operable stops running. A "done" operator will emit no further items.
func (o *Operable[T]) Done() <-chan struct{} {
	o.Start()

	return o.done
}

// Clone creates a copy of this stream
func (o *Operable[T]) Clone() observer.Stream[T] {
	o.Start()

	return o.output.Clone()
}

// ToSlice collects every emitted item until EOF is reached an returns an slice holding each collected item.
func (o *Operable[T]) ToSlice() []interface{} {
	defer o.complete()
	out := make([]interface{}, 0)
	done := o.ctx.Done()
	for {
		select {
		case <-done:
			return out
		case <-o.Changes():
			v := o.Next()
			if v == o.eof {
				return out
			}

			out = append(out, v)
		}
	}
}

// ToMap convert the sequence of emitted items into a map keyed by a specified key function.
func (o *Operable[T]) ToMap(keySelector Mapper[T]) map[interface{}]interface{} {
	defer o.complete()
	out := make(map[interface{}]interface{})
	done := o.ctx.Done()
	for {
		select {
		case <-done:
			return out
		case <-o.Changes():
			v := o.Next()
			if v == o.eof {
				return out
			}

			key := keySelector(o.ctx, v)
			out[key] = v
		}
	}
}

func (o *Operable[T]) run(ready chan struct{}) {
	defer func() {
		o.complete()
	}()

	close(ready)
	done := o.ctx.Done()
	for {
		select {
		case <-o.input.Changes():
			value := o.input.Next()

			o.mu.RLock()
			operators := o.operators
			o.mu.RUnlock()

			if value == o.eof {
				for _, operator := range operators {
					dst := make(chan T, 1)
					operator.end(dst)

					if v := rcv(dst, value); v != o.eof {
						o.surrogate.Update(v)
					}
				}

				return
			}

			stopped := false
			for _, operator := range operators {
				dst := make(chan T, 1)
				next := operator.next(value, dst)
				value = rcv(dst, value)

				if !next {
					stopped = true
					break
				}
			}

			if !stopped {
				o.surrogate.Update(value)
			}
		case <-done:
			return
		}
	}
}

func (o *Operable[T]) complete() {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.completed {
		return
	}

	if o.onComplete != nil {
		o.onComplete()
	}

	o.completed = true
	close(o.done)
}
