package rx

import (
	"time"
)

type operatorTimestamp struct{}

// TimestampItem attach a timestamp to an item.
type TimestampItem struct {
	Timestamp time.Time
	Item      interface{}
}

func (o *operatorTimestamp) next(item interface{}, dst chan<- interface{}) bool {
	send(dst, TimestampItem{
		Timestamp: time.Now().UTC(),
		Item:      item,
	})

	return true
}

func (o *operatorTimestamp) end(dst chan<- interface{}) {}

// Timestamp attaches a timestamp to each item indicating when it was emitted.
func (o *Operable) Timestamp() *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorTimestamp{})

	return o
}
