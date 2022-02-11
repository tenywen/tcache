package cache

import (
	"errors"
	"log"
	"sync/atomic"
)

type Cache struct {
	opt     opt
	stat    stat
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

	c := Cache{
		opt:     opt,
		shareds: make([]*shared, opt.nShared),
	}

	for i := 0; i < opt.nShared; i++ {
		c.shareds[i] = newShared(c.opt)
	}

	return c
}

func (c *Cache) Get(key string, dst []byte) ([]byte, error) {
	if len(key) > keyLimit {
		return nil, errKeyLimit
	}

	hash := defaultHasher.Sum64(key)
	return c.shareds[hash%uint64(c.opt.nShared)].get(c.opt.neverConflict, hash, key, dst)
}

func (c *Cache) Set(key string, value []byte) error {
	if len(key) > keyLimit {
		return errKeyLimit
	}

	if len(value) > valueLimit {
		return errValueLimit
	}

	hash := defaultHasher.Sum64(key)
	return c.shareds[hash%uint64(c.opt.nShared)].set(c.opt.neverConflict, hash, key, value)
}

func (c *Cache) Delete(key string) error {
	if len(key) > keyLimit {
		return errKeyLimit
	}

	hash := defaultHasher.Sum64(key)
	c.shareds[hash%uint64(c.opt.nShared)].delete(hash, key)
	return nil
}

func (c Cache) Debug() {
	var calls int64
	var miss int64
	var removes int64
	var totals int64

	for k := range c.shareds {
		calls += atomic.LoadInt64(&c.shareds[k].stat.calls)
		miss += atomic.LoadInt64(&c.shareds[k].stat.missCnt)
		removes += atomic.LoadInt64(&c.shareds[k].stat.removeBytes)
		totals += atomic.LoadInt64(&c.shareds[k].stat.totalBytes)
	}

	log.Println("cache stat debug")
	log.Printf("call:%d\n", calls)
	log.Printf("miss:%d  ratio:%.2f\n", miss, float64(miss)/float64(calls))
	log.Printf("remove:%s\n", printBytes(removes))
	log.Printf("total:%s\n", printBytes(totals))
	log.Println("")
}
