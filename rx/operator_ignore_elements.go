package rx

type operatorIgnoreElements struct{}

func (o *operatorIgnoreElements) next(item interface{}, dst chan<- interface{}) bool {
	return false
}

func (o *operatorIgnoreElements) end(dst chan<- interface{}) {}

// IgnoreElements do not emit any items but mirror its termination notification.
func (o *Operable[T]) IgnoreElements() *Operable[T] {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.operators = append(o.operators, &operatorIgnoreElements{})

	return o
}
