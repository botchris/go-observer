package observer

import (
	"context"
	"io"
	"sync"
)

// Operable defines a stream to which operators can be applied.
type Operable struct {
	ctx   context.Context
	input Stream

	mu        sync.RWMutex
	running   bool
	done      chan struct{}
	output    Stream
	operators []operator
	surrogate Property
}

// operator represents a component capable to altering/transforming the flow of items emitted by a source Stream.
// They can be stacked and are executed sequentially, so the changes made by one operator are visible to the
// operators coming after.
type operator interface {
	// next is triggered when the "source" Stream emits a new item.
	// Retuning FALSE indicates that no further operators should be applied, AND NOTHING should be written to the
	// output Stream; basically "discards" the item producing no visible changes at the putput stream.
	//
	// - item: is the value emitted by the source Stream, or emitted by the previous operator
	// - dst: this channel is used to indicate that the operator wants to change the given "item" value
	//   for another value written on the channel.
	next(item interface{}, dst chan<- interface{}) (next bool)

	// end is invoked when "source" Stream reaches io.EOF.
	// The given channel is used by the operator to indicate it wants to write something to the
	// output stream.
	end(dst chan<- interface{})
}

// send sends the given item over the given channel in a non-blocking fashion
func send(dst chan<- interface{}, item interface{}) {
	select {
	case dst <- item:
	default:
	}
}

// rcv receives a value from the given source or returns the given default value if there is nothing to read
func rcv(src <-chan interface{}, defaultValue interface{}) interface{} {
	select {
	case v := <-src:
		return v
	default:
	}

	return defaultValue
}

// MakeOperable makes the given input Stream operable so operators can be applied to. This new operable instance will be
// valid as long as the given context keeps active.
//
// By default operable streams are created in "Lazy" mode, which means they won't start reading  the input Stream until
// the operable itself is "consumed" (calls to any Stream interface method). This is specially useful when
// chaining multiple operators at once, so your "operators pipeline" is correctly defined upfront.
func MakeOperable(ctx context.Context, input Stream, opts ...Option) *Operable {
	p := NewProperty(nil)
	o := &Operable{
		ctx:   ctx,
		input: input.Clone(),

		done:      make(chan struct{}),
		output:    p.Observe(),
		operators: make([]operator, 0),
		surrogate: p,
	}

	options := &options{
		startMode: Lazy,
	}

	for _, o := range opts {
		o.apply(options)
	}

	if options.startMode == Eager {
		o.Start()
	}

	return o
}

// Start starts
func (o *Operable) Start() *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.running {
		return o
	}

	ready := make(chan struct{})
	go o.run(ready)

	<-ready
	o.running = true

	return o
}

// Value returns the current value for this stream.
func (o *Operable) Value() interface{} {
	o.Start()

	return o.output.Value()
}

// Changes returns the channel that is closed when a new value is available.
func (o *Operable) Changes() chan struct{} {
	o.Start()

	return o.output.Changes()
}

// Next advances this stream to the next state.
// You should never call this unless Changes channel is closed.
func (o *Operable) Next() interface{} {
	o.Start()

	return o.output.Next()
}

// Done returns a channel that's closed when Operable stops running. A "done" operator will emit no further items.
func (o *Operable) Done() <-chan struct{} {
	o.Start()

	return o.done
}

// HasNext checks whether there is a new value available.
func (o *Operable) HasNext() bool {
	o.Start()

	return o.output.HasNext()
}

// WaitNext waits for Changes to be closed, advances the stream and returns
// the current value.
func (o *Operable) WaitNext() interface{} {
	o.Start()

	return o.output.WaitNext()
}

// Clone creates a copy of this stream
func (o *Operable) Clone() Stream {
	o.Start()

	return o.output.Clone()
}

func (o *Operable) run(ready chan struct{}) {
	defer func() {
		o.surrogate.Update(io.EOF)
		close(o.done)
	}()

	close(ready)

	for {
		select {
		case <-o.input.Changes():
			value := o.input.Next()

			o.mu.RLock()
			operators := o.operators
			o.mu.RUnlock()

			if value == io.EOF {
				for _, operator := range operators {
					dst := make(chan interface{}, 1)
					operator.end(dst)

					if v := rcv(dst, value); v != io.EOF {
						o.surrogate.Update(v)
					}
				}

				return
			}

			stopped := false
			for _, operator := range operators {
				dst := make(chan interface{}, 1)
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
		case <-o.ctx.Done():
			return
		}
	}
}
