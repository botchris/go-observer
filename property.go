package observer

import (
	"sync"
)

// Property is an object that is continuously updated by one or more
// publishers. It is completely goroutine safe: you can use Property
// concurrently from multiple goroutines.
type Property[T any] interface {
	// Value returns the current value for this property.
	Value() T

	// Update sets a new value for this property.
	Update(value ...T)

	// End closes the property.
	End()

	// Observe returns a newly created Stream for this property.
	Observe() Stream[T]

	// Done returns a channel that is closed when property reaches a EOF
	Done() <-chan struct{}
}

// NewProperty creates a new Property with the initial value.
// It returns the created Property.
func NewProperty[T any](value T) Property[T] {
	var eof T

	return &property[T]{
		eof:   eof,
		state: newState[T](value),
		done:  make(chan struct{}),
	}
}

type property[T any] struct {
	sync.RWMutex
	eof   T
	ended bool
	done  chan struct{}
	state *state[T]
}

func (p *property[T]) Value() T {
	p.RLock()
	defer p.RUnlock()

	return p.state.value
}

func (p *property[T]) Update(values ...T) {
	p.Lock()
	defer p.Unlock()

	if p.ended {
		return
	}

	for _, value := range values {
		p.state = p.state.update(value)
	}
}

func (p *property[T]) Observe() Stream[T] {
	p.RLock()
	defer p.RUnlock()

	return &stream[T]{state: p.state}
}

func (p *property[T]) End() {
	close(p.done)
}

func (p *property[T]) Done() <-chan struct{} {
	return p.done
}
