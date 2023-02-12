package rx

type operatorLast struct {
	last  interface{}
	empty bool
}

func (o *operatorLast) next(item interface{}, dst chan<- interface{}) bool {
	o.last = item
	o.empty = false

	return false
}

func (o *operatorLast) end(dst chan<- interface{}) {
	if !o.empty {
		send(dst, o.last)
	}
}

// Last emit only the last item.
func (o *Operable[T]) Last() *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorLast{
		empty: true,
	})

	return o
}
