package rx

// Option handles configurable options.
type Option interface {
	apply(*options)
}

type options struct {
	startStrategy startStrategy
}

type funcOption struct {
	fn func(*options)
}

func (f *funcOption) apply(o *options) {
	f.fn(o)
}

// WithStartStrategy using the Eager strategy means the Operable will start consuming the input right after being created.
func WithStartStrategy(m startStrategy) Option {
	return &funcOption{
		fn: func(o *options) {
			o.startStrategy = m
		},
	}
}
