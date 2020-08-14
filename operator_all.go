package observer

type operatorAll struct {
	predicate Predicate
	all       bool
}

func (o *operatorAll) next(item interface{}, dst chan<- interface{}) bool {
	if !o.predicate(item) {
		o.all = false
	}

	return false
}

func (o *operatorAll) end(dst chan<- interface{}) {
	send(dst, o.all)
}

// All determine whether all items emitted meet some criteria.
func (o *Operable) All(predicate Predicate) *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorAll{
		predicate: predicate,
		all:       true,
	})

	return o
}
