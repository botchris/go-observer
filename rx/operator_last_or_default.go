package rx

type operatorLastOrDefault struct {
	defaultValue interface{}
	last         interface{}
	empty        bool
}

func (o *operatorLastOrDefault) next(item interface{}, dst chan<- interface{}) bool {
	o.last = item
	o.empty = false

	return false
}

func (o *operatorLastOrDefault) end(dst chan<- interface{}) {
	value := o.last
	if o.empty {
		value = o.defaultValue
	}

	send(dst, value)
}

// LastOrDefault emit only the last item. If fails to emit any items, it emits a default value.
func (o *Operable[T]) LastOrDefault(defaultValue interface{}) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorLastOrDefault{
		defaultValue: defaultValue,
		empty:        true,
	})

	return o
}
