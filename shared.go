package cache

import (
	"errors"
	"sync"
)

const (
	undefined = -1
)

var (
	errIdNotExist  = errors.New("id not exist")
	errKeyNotExist = errors.New("key not exist")
)

type shared struct {
	keys map[uint64]block // hash -> block

	recycle sortBlocks

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
		keys:      make(map[uint64]block),
		collision: make(map[string][]byte),
		buffer:    newBuffer(max),
	}
}

func (s *shared) get(hash uint64, key string) ([]byte, error) {
	s.m.RLock()
	// hit collision
	b, ok := s.collision[key]
	if ok {
		s.m.RUnlock()
		return b, nil
	}

	block, ok := s.keys[hash]
	if !ok {
		s.m.RUnlock()
		return nil, errKeyNotExist
	}

	si := block.si + block.kl
	v, err := s.buffer.read(si, block.vl+si)
	if err != nil {
		s.m.RUnlock()
		return nil, err
	}
	s.m.RUnlock()
	return v, nil
}

func (s *shared) set(hash uint64, key string, v []byte) error {
	s.m.Lock()
	b, ok := s.keys[hash]
	if ok {
		oldKey, _ := s.buffer.read(b.si, b.si+b.kl)
		if slice2string(oldKey) != key {
			s.collision[key] = v
			s.m.Unlock()
			return nil
		}

		if len(v) > b.vl {
			s.recycle.add(b)
			delete(s.keys, hash)
			ok = false
		}
	}

	k := string2slice(key)
	size := len(k) + len(v)

	if !ok {
		b, ok = s.recycle.getBlock(size)
		if !ok {
			b = block{
				si:    s.buffer.off,
				total: size,
			}
		}
	}

	b.kl = len(k)
	b.vl = len(v)
	err := s.buffer.write(b.si, size, k, v)
	if err != nil {
		s.m.Unlock()
		return err
	}
	s.keys[hash] = b
	s.m.Unlock()
	return nil
}

func (s *shared) del(hash uint64, key string) {
	s.m.Lock()
	if _, ok := s.collision[key]; ok {
		delete(s.collision, key)
		s.m.Unlock()
		return
	}

	block, ok := s.keys[hash]
	if !ok {
		s.m.Unlock()
		return
	}

	s.recycle.add(block)
	delete(s.keys, hash)
	s.m.Unlock()
}
