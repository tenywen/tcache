package cache

import (
	"errors"
	"sync"
)

const (
	undefined = -1

	flagLen  = 0
	totalLen = 4
	keyLen   = 2
	valueLen = 4
	headLen  = flagLen + totalLen + keyLen + valueLen
)

var (
	errIdNotExist  = errors.New("id not exist")
	errKeyNotExist = errors.New("key not exist")
)

type shared struct {
	keys map[uint64]int // hash -> si

	//recycle sortBlocks

	buffer buffer

	collision map[string][]byte

	m sync.RWMutex
}

type body struct {
	k []byte
	v []byte
}

func newShared(max int) *shared {
	return &shared{
		keys:      make(map[uint64]int),
		collision: make(map[string][]byte),
		buffer:    newBuffer(),
	}
}

func (shared shared) get(hash uint64, key string) ([]byte, error) {
	shared.m.RLock()
	defer shared.m.RUnlock()

	// hit collision
	v, ok := shared.collision[key]
	if ok {
		return v, nil
	}

	s, ok := shared.keys[hash]
	if !ok {
		return nil, errKeyNotExist
	}

	block, err := getBlock(s, &shared.buffer)
	if err != nil {
		return nil, err
	}

	v, err = shared.buffer.read(s+headLen+block.kl, s+headLen+block.kl+block.vl, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (shared *shared) set(hash uint64, key, v []byte) error {
	shared.m.Lock()
	defer shared.m.Unlock()

	s, ok := shared.keys[hash]

	if ok {
		block, err := getBlock(s, &shared.buffer)
		if err != nil {
			return err
		}

		var k []byte

		k, err = shared.buffer.read(s+headLen, s+headLen+block.kl, k)
		if err != nil {
			return nil
		}

		if slice2string(k) != slice2string(key) {
			shared.collision[slice2string(key)] = v
			return nil
		}

		return nil
		/*
			if len(v) > int(chunk.vl) {
				s.recycle.add(chunk)
				delete(s.keys, hash)
				ok = false
			}
		*/
	}

	var block block

	block.s = shared.buffer.off
	block.kl = len(key)
	block.vl = len(v)
	block.total = block.kl + block.vl + headLen
	err := putBlock(block, &shared.buffer)
	if err != nil {
		return err
	}

	shared.buffer.write(block.s+headLen, key)
	shared.buffer.write(block.s+headLen+block.kl, v)

	shared.keys[hash] = block.s
	return nil
}

/*
func (s *shared) del(hash uint64, key string) {
	s.m.Lock()
	defer s.m.Unlock()
	_, ok := s.collision[key]
	if ok {
		delete(s.collision, key)
		return
	}

	chunk := getChunk()
	defer recycleChunk(chunk)

	chunk.s, ok = s.keys[hash]
	if !ok {
		return
	}

	err := s.buffer.decode(chunk)
	if err != nil {
		return
	}

	s.recycle.add(chunk)
	delete(s.keys, hash)
}
*/
