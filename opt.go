package cache

const (
	defaultKey    = 1 << 8  // 256B
	defaultValue  = 1 << 16 // 64K
	defaultShared = 1 << 8  // 256
	defaultMax    = 1 << 30 // 1GB
)

type opts func(*opt)

type opt struct {
	overwrite bool
	keyMax    int // bytes
	valueMax  int // bytes
	maxSize   int // bytes
	nShared   int
}

func WithKeyMax(limit int) opts {
	return func(o *opt) {
		o.keyMax = limit
	}
}

func WithValueMax(limit int) opts {
	return func(o *opt) {
		o.valueMax = limit
	}
}

func WithMaxBuffer(max int) opts {
	return func(o *opt) {
		o.maxSize = max
	}
}

func WithShared(n int) opts {
	return func(o *opt) {
		o.nShared = n
	}
}

func defaultOpt() opt {
	return opt{
		keyMax:    defaultKey,
		valueMax:  defaultValue,
		maxSize:   defaultMax,
		nShared:   defaultShared,
		overwrite: true,
	}
}
