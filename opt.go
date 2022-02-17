package tcache

const (
	keyLimit   = 1 << 16 // 64K
	valueLimit = 1 << 16 // 64K
)

type opts func(*opt)

type opt struct {
	neverConflict bool
	threshold     float64
	nShared       int
}

// avoid conflict
// default: allowed conflict
func WithNeverConflict() opts {
	return func(o *opt) {
		o.neverConflict = true
	}
}

// cache shared numbera
// defalut: 512
func WithShared(n int) opts {
	return func(o *opt) {
		o.nShared = n
	}
}

// trigger recycle ratio(removes/total)
// default: can't trigger recycle
func WithRecycleThreshold(v float64) opts {
	return func(o *opt) {
		o.threshold = v
	}
}

func defaultOpt() opt {
	return opt{
		neverConflict: false,
		nShared:       512,
		threshold:     undefined,
	}
}
