package cache

import (
	"errors"
	"sync"
)

const (
	undefined = -1

	flagLen  = 1
	totalLen = 3
	keyLen   = 2
	valueLen = 2
	bufLen   = flagLen + totalLen + keyLen + valueLen
)

var (
	errKeyNotExist = errors.New("key not exist")
	errReadPos     = errors.New("read start pos or end pos error")
	errCapLimit    = errors.New("buffer cap exceeds limit ")
	errNotContent  = errors.New("not content in this start pos")
)

type shared struct {
	keys map[uint64]int // hash -> si

	//recycle sortBlocks

	off    int
	chunks [][]byte
	kvBuf  [bufLen]byte

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

	kvBuf, err := shared.read(s, s+bufLen, shared.kvBuf[:])
	if err != nil {
		return nil, err
	}

	kLen := int(kvBuf[4])<<8 + int(kvBuf[5])
	vLen := int(kvBuf[6])<<8 + int(kvBuf[7])

	v, err = shared.read(s+bufLen+kLen, s+bufLen+kLen+vLen, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (shared shared) totalLen() int {
	return int(shared.kvBuf[1]>>16) + int(shared.kvBuf[2]>>8) + int(shared.kvBuf[3])
}

func (shared shared) keyLen() int {
	return int(shared.kvBuf[4]>>8) + int(shared.kvBuf[5])
}

func (shared shared) valueLen() int {
	return int(shared.kvBuf[6]>>8) + int(shared.kvBuf[7])
}

func (shared *shared) set(hash uint64, key, v []byte) error {
	shared.m.Lock()
	defer shared.m.Unlock()

	s, ok := shared.keys[hash]
	if ok {
		_, err := shared.read(s, s+bufLen, shared.kvBuf[:])
		if err != nil {
			return err
		}

		keyLen := shared.keyLen()
		k, err := shared.read(s+bufLen, s+bufLen+keyLen, nil)
		if err != nil {
			return err
		}

		if slice2string(k) != slice2string(key) {
			shared.collision[slice2string(key)] = v
			return nil
		}

		if shared.valueLen() < len(v) {
			shared.delete()
			ok = false
		}
	}

	length := len(key) + len(v)
	shared.kvBuf[1] = byte(length >> 16)
	shared.kvBuf[2] = byte(length >> 8)
	shared.kvBuf[3] = byte(length)
	shared.kvBuf[6] = byte(len(v) >> 8)
	shared.kvBuf[7] = byte(len(v))

	if !ok {
		s = shared.off
		shared.kvBuf[0] = 1
		shared.kvBuf[4] = byte(len(key) >> 8)
		shared.kvBuf[5] = byte(len(key))
	}

	err := shared.write(s, shared.kvBuf[:])
	if err != nil {
		return err
	}

	if !ok {
		err = shared.write(s+bufLen, key)
		if err != nil {
			return err
		}
	}
	err = shared.write(s+bufLen+len(key), v)
	if err != nil {
		return err
	}

	shared.keys[hash] = s
	return nil
}

func (shared *shared) delete() {

}

/*
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
				if len(v) > int(chunk.vl) {
					s.recycle.add(chunk)
					delete(s.keys, hash)
					ok = false
				}
		}

	length := len(key) + len(v)
	shared.kvBuf[0] = 1
	shared.kvBuf[1] = byte(length >> 8)
	shared.kvBuf[2] = byte(length)
	shared.kvBuf[3] = byte(len(key) >> 8)
	shared.kvBuf[4] = byte(len(key))
	shared.kvBuf[5] = byte(len(v) >> 8)
	shared.kvBuf[6] = byte(len(v))

	s = shared.off
	shared.write(shared.off, shared.kvBuf[:])
	shared.write(shared.off, key)
	shared.write(shared.off, v)

	shared.keys[hash] = s
	return nil
}
*/

func (shared *shared) write(s int, v []byte) error {
	if s > shared.off {
		return errReadPos
	}

	length := len(v)
	offset := s % chunkSize
	idx := s / chunkSize
	for len(v) != 0 {
		if len(shared.chunks) == idx {
			shared.chunks = append(shared.chunks, getChunk())
		}

		chunk := shared.chunks[idx]
		v = v[copy(chunk[offset:], v):]
		offset = 0
		idx++
	}

	if s+length > shared.off {
		shared.off = s + length
	}
	return nil
}

func (shared *shared) read(s, e int, dst []byte) ([]byte, error) {
	if s > e || s < 0 || e > shared.off {
		return nil, errReadPos
	}

	length := e - s

	if dst == nil {
		dst = make([]byte, length)
	}

	if len(dst) > length {
		dst = dst[:length]
	}

	tmp := dst
	offset := s % chunkSize
	idx := s / chunkSize

	for len(tmp) > 0 {
		chunk := shared.chunks[idx]
		tmp = tmp[copy(tmp, chunk[offset:]):]
		offset = 0
		idx++
	}

	return dst, nil
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
