package cache

import (
	"errors"
)

const (
	undefined = -1
)

var (
	errIdNotExist  = errors.New("id not exist")
	errKeyNotExist = errors.New("key not exist")
)

type shared struct {
	id int64

	keys map[uint64]block // hash -> block

	recycle sortBlocks

	buffer buffer

	collision map[string][]byte
}

type body struct {
	k []byte
	v []byte
}

func newShared() shared {
	return shared{
		keys:      make(map[uint64]block),
		collision: make(map[string][]byte),
		buffer:    newBuffer(1, 1024),
	}
}

func (s shared) get(key string, hasher Hasher) ([]byte, error) {
	// hit collision
	b, ok := s.collision[key]
	if ok {
		return b, nil
	}

	hash := hasher.Sum64(key)
	block, ok := s.keys[hash]
	if !ok {
		return nil, errKeyNotExist
	}

	si := block.si + block.kl
	v, err := s.buffer.read(si, block.vl+si)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *shared) set(key string, v []byte, hasher Hasher) {
	hash := hasher.Sum64(key)
	if _, ok := s.keys[hash]; ok {
		s.collision[key] = v
		return
	}

	k := string2slice(key)
	size := len(k) + len(v)

	b, ok := s.recycle.getBlock(size)
	if !ok {
		b = block{
			si:    s.buffer.off,
			total: size,
		}
	}

	b.kl = len(k)
	b.vl = len(v)
	s.buffer.write(b.si, size, k, v)
	s.keys[hash] = b
}

func (s *shared) del(key string, hasher Hasher) {
	if _, ok := s.collision[key]; ok {
		delete(s.collision, key)
		return
	}

	hash := hasher.Sum64(key)
	block, ok := s.keys[hash]
	if !ok {
		return
	}

	s.recycle.add(block)
	delete(s.keys, hash)
}

/*
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
*/
