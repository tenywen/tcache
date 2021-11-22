package cache

import (
	"errors"
	"math"
)

type Cache struct {
	opt     opt
	shareds []*shared
}

var (
	errKeyLimit   = errors.New("key too large")
	errValueLimit = errors.New("value too large")
)

func New(opts ...opts) Cache {
	opt := defaultOpt()
	for k := range opts {
		opts[k](&opt)
	}

	cache := Cache{
		opt:     opt,
		shareds: make([]*shared, opt.nShared),
	}

	max := math.Ceil(float64(opt.maxSize) / float64(opt.nShared))

	for i := 0; i < cache.opt.nShared; i++ {
		cache.shareds[i] = newShared(int(max))
	}

	return cache
}

func (cache *Cache) Get(key string) ([]byte, error) {
	if len(key) > cache.opt.keyMax {
		return nil, errKeyLimit
	}

	hash := defaultHasher.Sum64(key)
	return cache.shareds[hash%uint64(cache.opt.nShared)].get(hash, key)
}

func (cache *Cache) Set(key string, value []byte) error {
	if len(key) > cache.opt.keyMax {
		return errKeyLimit
	}

	if len(value) > cache.opt.valueMax {
		return errValueLimit
	}

	hash := defaultHasher.Sum64(key)
	return cache.shareds[hash%uint64(cache.opt.nShared)].set(hash, key, value)
}

func (cache *Cache) Delete(key string) error {
	if len(key) > cache.opt.keyMax {
		return errKeyLimit
	}

	hash := defaultHasher.Sum64(key)
	cache.shareds[hash%uint64(cache.opt.nShared)].del(hash, key)
	return nil
}
