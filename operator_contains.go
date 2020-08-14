package observer

type operatorContains struct {
	predicate Predicate
	contains  bool
	stopped   bool
}

func (o *operatorContains) next(item interface{}, dst chan<- interface{}) bool {
	if !o.stopped && o.predicate(item) {
		o.contains = true
		o.stopped = true
		send(dst, true)

		return true
	}

	return false
}

func (o *operatorContains) end(dst chan<- interface{}) {
	if !o.stopped && !o.contains {
		send(dst, false)
	}
}

// Contains determine whether a particular item was emitted or not.
func (o *Operable) Contains(predicate Predicate) *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorContains{
		predicate: predicate,
	})

	return o
}
