package rx

type operatorBufferWithCount struct {
	size   int
	count  int
	buffer []interface{}
}

func (o *operatorBufferWithCount) next(item interface{}, dst chan<- interface{}) bool {
	o.buffer[o.count] = item
	o.count++

	if o.count == o.size {
		o.count = 0
		send(dst, o.buffer)
		o.buffer = make([]interface{}, o.size)

		return true
	}

	return false
}

func (o *operatorBufferWithCount) end(dst chan<- interface{}) {
	if o.count != 0 {
		send(dst, o.buffer[:o.count])
	}
}

// BufferWithCount periodically gather items emitted into bundles and emit these bundles rather than emitting the items one at a time.
func (o *Operable) BufferWithCount(size int) *Operable {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorBufferWithCount{
		size:   size,
		count:  0,
		buffer: make([]interface{}, size),
	})

	return o
}
