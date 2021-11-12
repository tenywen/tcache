package cache

const (
	defaultKey    = 1 << 8
	defaultValue  = 1 << 16
	defaultShared = 64
	defaultMax    = 1 << 30
)

type opts func(*opt)

type opt struct {
	keyLimit   int // bytes
	valueLimit int // bytes
	max        int // bytes
	nShared    int
}

func WithKeyLimit(limit int) opts {
	return func(o *opt) {
		o.keyLimit = limit
	}
}

func WithValueLimit(limit int) opts {
	return func(o *opt) {
		o.valueLimit = limit
	}
}

func WithMaxBuffer(max int) opts {
	return func(o *opt) {
		o.max = max
	}
}

func WithShared(n int) opts {
	return func(o *opt) {
		o.nShared = n
	}
}

func defaultOpt() opt {
	return opt{
		keyLimit:   defaultKey,
		valueLimit: defaultValue,
		max:        defaultMax,
		nShared:    defaultShared,
	}
}
