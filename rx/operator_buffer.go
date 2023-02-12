package rx

type operatorBufferWithCount[T any] struct {
	size   int
	count  int
	buffer []interface{}
}

func (o *operatorBufferWithCount[T]) next(item T, dst chan<- T) bool {
	o.buffer[o.count] = item
	o.count++

	if o.count == o.size {
		o.count = 0
		send(dst, o.buffer)
		o.buffer = make([]T, o.size)

		return true
	}

	return false
}

func (o *operatorBufferWithCount[T]) end(dst chan<- T) {
	if o.count != 0 {
		send(dst, o.buffer[:o.count])
	}
}

// BufferWithCount periodically gather items emitted into bundles and emit these bundles rather than emitting the items one at a time.
func (o *Operable[T]) BufferWithCount(size int) *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorBufferWithCount[T]{
		size:   size,
		count:  0,
		buffer: make([]T, size),
	})

	return o
}
