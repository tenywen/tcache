package cache

import "errors"

const (
	undefined = -1

	idLen    = 8
	keyLen   = 8
	valueLen = 8
	headLen  = idLen + keyLen + valueLen
)

var (
	errIdNotExist  = errors.New("id not exist")
	errKeyNotExist = errors.New("key not exist")
)

type shared struct {
	id int64
	// id -> start index
	ids map[int64]int64

	// hash -> id
	keys map[uint64]int64

	buffer buffer

	collision map[string][]byte
}

type head struct {
	id int64
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
		id:        1,
		ids:       make(map[int64]int64),
		keys:      make(map[uint64]int64),
		collision: make(map[string][]byte),
	}
}

func (s shared) get(key string, hasher Hasher) ([]byte, error) {
	hash := hasher.Sum64(key)
	id, ok := s.keys[hash]
	if !ok {
		return nil, errKeyNotExist
	}

	body, err := s.decode(id)
	if err != nil {
		return nil, err
	}

	if string2slice(body.k) != key {
		return nil, nil
	}
}

func (s *shared) set(k string, v []byte, hasher Hasher) {
	hash := hasher.Sum64(k)
	id, ok := s.keys[hash]
	if ok {
		// TODO
		s.del(id)
		return
	}

	s.write(string2slice(k), v)
}

func (s *shared) del(id int64) {

}

func (s shared) decode(id int64) (body body, err error) {
	si, ok := s.ids[id]
	if !ok {
		err = errIdNotExist
		return
	}

	return s.body(si)
}

func (s shared) body(si int64) (body body, err error) {
	kl := s.buffer.int64(si+idLen, si+idLen+keyLen)
	vl := s.buffer.int64(si+idLen+keyLen, si+headLen)

	body.k, err = s.buffer.read(si+headLen, si+headLen+kl)
	if err != nil {
		return
	}

	body.v, err = s.buffer.read(si+headLen+kl, si+headLen+kl+vl)
}

func (s *shared) write(k, v []byte) {
	s.id++
	id := s.id
	s.ids[id] = s.buffer.tail()
	s.buffer.putInt64(id)
	s.buffer.putInt64(int64(len(k)))
	s.buffer.putInt64(int64(len(v)))
	s.buffer.write(k)
	s.buffer.write(v)
}
