package observer

import (
	"io"
	"sync"
)

// Property is an object that is continuously updated by one or more
// publishers. It is completely goroutine safe: you can use Property
// concurrently from multiple goroutines.
type Property interface {
	// Value returns the current value for this property.
	Value() interface{}

	// Update sets a new value for this property.
	// Updating with `io.EOF` will mark this property as "ended",
	// once ended further calls to Update will no-op
	Update(value interface{})

	// Observe returns a newly created Stream for this property.
	Observe() Stream

	// Emits a io.EOF message, shortcut for `Property.Update(io.EOF)`
	End()

	// Done returns a channel that is closed when property reaches a EOF
	Done() <-chan struct{}
}

// NewProperty creates a new Property with the initial value value.
// It returns the created Property.
func NewProperty(value interface{}) Property {
	return &property{
		state: newState(value),
		done:  make(chan struct{}),
	}
}

type property struct {
	sync.RWMutex
	ended bool
	done  chan struct{}
	state *state
}

func (p *property) Value() interface{} {
	p.RLock()
	defer p.RUnlock()
	return p.state.value
}

func (p *property) Update(value interface{}) {
	p.Lock()
	defer p.Unlock()

	if p.ended {
		return
	}

	if value == io.EOF {
		p.ended = true
		close(p.done)
	}

	p.state = p.state.update(value)
}

func (p *property) Observe() Stream {
	p.RLock()
	defer p.RUnlock()
	return &stream{state: p.state}
}

func (p *property) End() {
	p.Update(io.EOF)
}

func (p *property) Done() <-chan struct{} {
	return p.done
}
