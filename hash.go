package tcache

import xxhash "github.com/cespare/xxhash/v2"

type Hasher interface {
	Sum64(string) uint64
}

type defaultHash struct{}

func newDefaultHash() defaultHash {
	return defaultHash{}
}

func (d defaultHash) Sum64(key string) uint64 {
	return xxhash.Sum64(string2slice(key))
}

var defaultHasher = newDefaultHash()
