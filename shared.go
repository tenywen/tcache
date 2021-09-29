package cache

import "errors"

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

	// hash -> start
	keys map[uint64]int64

	buffer buffer

	collision map[string][]byte
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

	if slice2string(body.k) != key {
		return s.getCollision(key)
	}

	return body.v, nil
}

func (s shared) getCollision(key string) ([]byte, error) {
	b, ok := s.collision[key]
	if !ok {
		return nil, errKeyNotExist
	}
	return b, nil
}

func (s *shared) setCollision(key string, v []byte) {

}

func (s *shared) set(k string, v []byte, hasher Hasher) {
	hash := hasher.Sum64(k)
	_, ok := s.keys[hash]
	if ok {
		s.collision[k] = v
		return
	}

	s.keys[hash] = s.write(string2slice(k), v)
}

func (s *shared) del(key string, hasher Hasher) {
	hash := hasher.Sum64(key)
	start, ok := s.keys[hash]
	if ok {
		body, err := s.body(start)
		if err == nil {
			if slice2string(body.k) == key {
				return s.getCollision(key)
			}
		}
		return nil, errKeyNotExist
	}

	body, err := s.body(start)
	if err != nil {
		return nil, err
	}

	if slice2string(body.k) != key {
		return s.getCollision(key)
	}
}

func (s shared) body(si int64) (body body, err error) {
	kl := s.buffer.int64(si)
	vl := s.buffer.int64(si + keyLen)

	body.k, err = s.buffer.read(si+headLen, si+headLen+kl)
	if err != nil {
		return
	}

	body.v, err = s.buffer.read(si+headLen+kl, si+headLen+kl+vl)
}

func (s *shared) write(k, v []byte) int64 {
	start := s.buffer.off
	s.buffer.putInt64(int64(len(k)))
	s.buffer.putInt64(int64(len(v)))
	s.buffer.write(k)
	s.buffer.write(v)
	return start
}
