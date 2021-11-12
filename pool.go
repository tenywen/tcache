package cache

import (
	"math"
	"sync"
)

type pools struct {
	min   int
	max   int
	pools []*sync.Pool
}

func newPool(size int) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return make([]byte, size)
		},
	}
}

func newPools(min, max int) pools {
	if min > max {
		panic("pool maxSize less than minSize")
	}

	ps := pools{
		min: min,
		max: max,
	}

	cur := ps.min
	for cur < max {
		ps.pools = append(ps.pools, newPool(cur))
		cur = cur << 1
	}

	ps.pools = append(ps.pools, newPool(max))
	return ps
}

func (ps *pools) get(size int) []byte {
	pool := ps.getPool(size)
	if pool == nil {
		return make([]byte, size)
	}

	b := pool.Get().([]byte)
	return b[:size]
}

func (ps *pools) put(b []byte) {
	pool := ps.getPool(cap(b))
	if pool == nil {
		return
	}
	pool.Put(b)
}

func (ps *pools) getPool(size int) *sync.Pool {
	if size > ps.max {
		return nil
	}

	if size < ps.min {
		return ps.pools[0]
	}

	idx := int(math.Ceil(math.Log2(float64(size) / float64(ps.min))))
	if idx < 0 {
		idx = 0
	}

	return ps.pools[idx]
}
