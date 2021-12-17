package cache

import (
	"errors"
	"sync"
)

const (
	undefined = -1

	flagLen  = 1
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
		keys:      make(map[uint64]int),
		collision: make(map[string][]byte),
		buffer:    newBuffer(),
	}
}

func (s shared) get(hash uint64, key string) ([]byte, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	// hit collision
	b, ok := s.collision[key]
	if ok {
		return b, nil
	}

	si, ok := s.keys[hash]
	if !ok {
		return nil, errKeyNotExist
	}

	chunk := getChunk()
	defer recycleChunk(chunk)

	chunk.s = si

	err := s.buffer.decode(chunk)
	if err != nil {
		return nil, err
	}
	return chunk.v, nil
}

func (shared *shared) set(hash uint64, key string, v []byte) error {
	shared.m.Lock()
	defer shared.m.Unlock()

	s, ok := shared.keys[hash]

	if ok {
		kvBuff := [headLen]byte{}
		err := shared.buffer.read(s, s+headLen, kvBuff[:])
		if err != nil {
			return err
		}

		err := s.buffer.decode(chunk)
		if err != nil {
			return err
		}

		if slice2string(chunk.k) != key {
			s.collision[key] = v
			return nil
		}

		if len(v) > int(chunk.vl) {
			s.recycle.add(chunk)
			delete(s.keys, hash)
			ok = false
		}
	}

	// 已经存在
	if ok {
		if slice2string(v) == slice2string(chunk.v) {
			return nil
		}
		chunk.kl = 0
		chunk.k = nil
		chunk.vl = int32(len(v))
		chunk.total = int32(chunk.kl) + chunk.vl
		chunk.v = v
	} else {
		k := string2slice(key)
		total := int32(len(k) + len(v))
		if !ok {
			ok = s.recycle.getBlock(total, chunk)
			if !ok {
				chunk.used = ^unused
				chunk.total = total
				chunk.s = s.buffer.off
			}
		}

		chunk.kl = int16(len(k))
		chunk.vl = int32(len(v))
		chunk.k = k
		chunk.v = v
	}

	err := s.buffer.encode(chunk)
	if err != nil {
		return err
	}
	s.keys[hash] = chunk.s
	return nil
}

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
