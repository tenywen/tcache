package cache

import (
	"errors"
	"sort"
)

const (
	undefined = -1

	keyLen   = 8
	valueLen = 8
	headLen  = keyLen + valueLen
)

var (
	errIdNotExist  = errors.New("id not exist")
	errKeyNotExist = errors.New("key not exist")
)

type shared struct {
	id int64

	keys map[uint64]int // hash -> start

	remove remove

	buffer buffer

	collision map[string][]byte

	ps pools
}

type head struct {
	kl int64 // key length
	vl int64 // value length
}

type body struct {
	k []byte
	v []byte
}

type data struct {
	head head
	body body
}

func newShared() shared {
	return shared{
		keys:      make(map[uint64]int),
		collision: make(map[string][]byte),
		ps:        newPools(8, 1<<16),
	}
}

func (s shared) get(key string, hasher Hasher) ([]byte, error) {
	// hit collision
	b, ok := s.collision[key]
	if ok {
		return b, nil
	}

	hash := hasher.Sum64(key)
	start, ok := s.keys[hash]
	if !ok {
		return nil, errKeyNotExist
	}

	body, err := s.body(start)
	if err != nil {
		return nil, err
	}

	return body.v, nil
}

func (s *shared) set(key string, v []byte, hasher Hasher) {
	hash := hasher.Sum64(key)
	if _, ok := s.keys[hash]; ok {
		s.collision[key] = v
		return
	}

	k := string2slice(key)
	size := headLen + len(k) + len(v)

	start := s.remove.getBlock(size)
	if start != -1 {
		s.buffer.rewrite(start, size, k, v)
		s.keys[hash] = start
		return
	}

	s.keys[hash] = s.write(k, v)
}

func (s *shared) del(key string, hasher Hasher) {
	if _, ok := s.collision[key]; ok {
		delete(s.collision, key)
		return
	}

	hash := hasher.Sum64(key)
	start, ok := s.keys[hash]
	if !ok {
		return
	}

	body, err := s.body(start)
	if err != nil || slice2string(body.k) != key {
		return
	}

	s.remove.add(int(start), headLen+len(body.k)+len(body.v))
	delete(s.keys, hash)
}

func (s shared) body(si int) (body body, err error) {
	kl := s.buffer.btoi(si)
	vl := s.buffer.btoi(si + keyLen)

	body.k, err = s.buffer.read(si+headLen, si+headLen+kl)
	if err != nil {
		return
	}

	body.v, err = s.buffer.read(si+headLen+kl, si+headLen+kl+vl)
	return
}

func (s *shared) shrink(hasher Hasher) {
	sort.Sort(sorts(s.removes))
	start := s.removes[0].start
	val := 0
	for i := 0; i < len(s.removes)-1; i++ {
		val += s.removes[i].end - s.removes[i].start
		err := s.adjust(int64(s.removes[i].end+1), int64(s.removes[i+1].start-1), int64(-val), hasher)
		if err != nil {
			panic(err)
		}
		println("start:", start)
		start += copy(s.buffer.bytes[start:], s.buffer.bytes[s.removes[i].end+1:s.removes[i+1].start-1])
	}
}

func (s *shared) adjust(start, end, val int64, hasher Hasher) error {
	for pos := start; pos < end; {
		body, err := s.body(pos)
		if err != nil {
			return err
		}

		hash := hasher.Sum64(slice2string(body.k))
		s.keys[hash] += val
		pos += int64(len(body.k) + len(body.v))
	}

	return nil
}

func (s *shared) write(k, v []byte) int {
	start := s.buffer.off
	s.buffer.writeInt(len(k))
	s.buffer.writeInt(len(v))
	s.buffer.write(k)
	s.buffer.write(v)
	return start
}
