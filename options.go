package observer

// Option handles configurable options.
type Option interface {
	apply(*options)
}

type options struct {
	startMode startStrategy
}

type funcOption struct {
	fn func(*options)
}

func (f *funcOption) apply(o *options) {
	f.fn(o)
}

// WithStarStrategy using the Eager strategy means the Operable will start consuming the input right after being created.
func WithStarStrategy(m startStrategy) Option {
	return &funcOption{
		fn: func(o *options) {
			o.startMode = m
		},
	}
}
