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

type pos struct {
	start int
	end   int
}

type shared struct {
	id int64

	keys map[uint64]int64 // hash -> start

	removes []pos

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
		keys:      make(map[uint64]int64),
		collision: make(map[string][]byte),
		ps:        newPools(8, 1<<16),
	}
}

func (s shared) get(key string, hasher Hasher) ([]byte, error) {
	hash := hasher.Sum64(key)
	start, ok := s.keys[hash]
	if !ok {
		return nil, errKeyNotExist
	}

	body, err := s.body(start)
	if err != nil {
		return nil, err
	}

	if slice2string(body.k) == key {
		return body.v, nil
	}

	// hit collision
	b, ok := s.collision[key]
	if !ok {
		return nil, errKeyNotExist
	}

	return b, nil
}

func (s *shared) set(k string, v []byte, hasher Hasher) {
	hash := hasher.Sum64(k)
	if _, ok := s.keys[hash]; ok {
		s.collision[k] = v
		return
	}

	s.keys[hash] = s.write(string2slice(k), v)
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

	s.removes = append(s.removes, pos{start: int(start), end: int(start) + headLen + len(body.k) + len(body.v) - 1})
	delete(s.keys, hash)
}

func (s shared) body(si int64) (body body, err error) {
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

func (s *shared) write(k, v []byte) int64 {
	start := s.buffer.off
	s.buffer.writeInt64(int64(len(k)))
	s.buffer.writeInt64(int64(len(v)))
	s.buffer.write(k)
	s.buffer.write(v)
	return start
}

func (s *shared) writeInt16(k int16) (start int64) {
	start = s.buffer.off
	//s.buffer.write()
	return
}

func (s *shared) writeInt32() (start int64) {
	return
}

func (s *shared) writeInt64() (start int64) {
	return
}

func (s shared) readInt() {}

type sorts []pos

func (s sorts) Less(i, j int) bool {
	return s[i].start > s[j].start
}

func (s sorts) Len() int {
	return len(s)
}

func (s sorts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
