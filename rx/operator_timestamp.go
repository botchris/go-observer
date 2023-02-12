package rx

import (
	"time"
)

type operatorTimestamp[T any] struct{}

// TimestampItem attach a timestamp to an item.
type TimestampItem[T any] struct {
	Timestamp time.Time
	Item      T
}

func (o *operatorTimestamp[T]) next(item T, dst chan<- T) bool {
	send[T](dst, TimestampItem[T]{
		Timestamp: time.Now().UTC(),
		Item:      item,
	})

	return true
}

func (o *operatorTimestamp[T]) end(dst chan<- T) {}

// Timestamp attaches a timestamp to each item indicating when it was emitted.
func (o *Operable[T]) Timestamp() *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorTimestamp[T]{})

	return o
}
