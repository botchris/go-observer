package rx

import "time"

type operatorDebounce struct {
	last     time.Time
	timespan time.Duration
}

func (o *operatorDebounce) next(item interface{}, dst chan<- interface{}) bool {
	now := time.Now()
	if now.After(o.last.Add(o.timespan)) {
		o.last = now
		send(dst, item)

		return true
	}

	return false
}

func (o *operatorDebounce) end(dst chan<- interface{}) {}

// Debounce only emit an item if a particular timespan has passed without it emitting another item.
func (o *Operable[T]) Debounce(timespan time.Duration) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorDebounce{
		last:     time.Now(),
		timespan: timespan,
	})

	return o
}
