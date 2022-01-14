package cache

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"sync/atomic"
)

const (
	undefined = -1
)

var (
	errKeyNotExist = errors.New("key not exist")
	errStartPos    = errors.New("read start pos or end pos error")
	errCapLimit    = errors.New("shared buffer cap exceeds limit 4G")
	errNotContent  = errors.New("not content in this start pos")
)

type shared struct {
	keys map[uint64]block // (key's hash -> block)

	off     int32
	chunks  [][]byte
	removes []block

	stat      stat
	collision map[string][]byte

	mu sync.RWMutex
}

type body struct {
	k []byte
	v []byte
}

func newShared(opt opt) *shared {
	shared := &shared{
		keys: make(map[uint64]block),
	}

	if opt.neverConflict {
		shared.collision = make(map[string][]byte)
	}

	return shared
}

func (shared *shared) get(neverConflict bool, hash uint64, key string, dst []byte) (v []byte, err error) {
	var k []byte
	shared.stat.call()
	shared.mu.RLock()

	block, ok := shared.keys[hash]
	if ok {
		s := block.s
		k, err = shared.read(s, s+int32(block.kl), nil)
		if err != nil {
			goto END
		}

		s += int32(block.kl)
		if slice2string(k) == key {
			v, err = shared.read(s, s+int32(block.vl), dst)
			goto END
		}
	}

	if !neverConflict {
		err = errKeyNotExist
		goto END
	}

	// hit collision
	v, ok = shared.collision[key]
	if !ok {
		err = errKeyNotExist
	}

END:
	shared.mu.RUnlock()
	if err != nil {
		shared.stat.miss()
	}
	return v, err
}

func (shared *shared) set(neverConflict bool, hash uint64, key string, v []byte) (err error) {
	total := uint16(len(key) + len(v))
	shared.stat.call()
	shared.mu.Lock()
	block, ok := shared.keys[hash]
	if ok {
		if neverConflict {
			k, err := shared.read(block.s, block.s+int32(block.kl), nil)
			if err != nil {
				goto END
			}

			if slice2string(k) != key {
				shared.collision[key] = v
				goto END
			}
		}

		if total <= block.total {
			goto SET
		}

		ok = false
		shared.remove(hash, key)
		if !neverConflict {
			shared.stat.collision()
		}
	}

SET:
	if !ok {
		block, ok = shared.getBlock(total)
		if ok {
			shared.stat.remove(int64(-total))
		}
	}

	block.kl = int16(len(key))
	block.vl = int16(len(v))

	if !ok {
		block.s = shared.off
		block.total = uint16(total)
		err = shared.write(block.s, string2slice(key))
		if err != nil {
			goto END
		}
	}

	err = shared.write(block.s+int32(block.kl), v)
	if err != nil {
		goto END
	}

	if !ok {
		shared.stat.add(int64(block.total))
	}
	shared.keys[hash] = block

END:
	shared.mu.Unlock()
	return
}

func (shared shared) debug() {
	for _, chunk := range shared.chunks {
		if int32(len(chunk)) > shared.off {
			fmt.Printf("%s\n", chunk[:shared.off])
		} else {
			fmt.Printf("%s\n", chunk[:])
		}

		shared.off -= int32(len(chunk))
		if shared.off <= 0 {
			break
		}
	}
}

func (shared *shared) delete(hash uint64, key string) {
	shared.mu.Lock()
	defer shared.mu.Unlock()
	shared.remove(hash, key)
}

func (shared *shared) remove(hash uint64, key string) {
	_, ok := shared.collision[key]
	if ok {
		delete(shared.collision, key)
		return
	}

	if block, ok := shared.keys[hash]; ok {
		shared.removes = append(shared.removes, block)
		shared.stat.remove(int64(block.total))
		delete(shared.keys, hash)
	}
}

func (shared *shared) getBlock(size uint16) (block block, ok bool) {
	if len(shared.removes) == 0 {
		return
	}

	sort.Sort(sortBlocks(shared.removes))
	length := len(shared.removes)
	i := sort.Search(length, func(i int) bool {
		return shared.removes[i].total >= size
	})

	if i < length {
		block = shared.removes[i]
		shared.removes[i], shared.removes[length-1] = shared.removes[length-1], shared.removes[i]
		shared.removes = shared.removes[:length-1]
		ok = true
		shared.stat.remove(int64(block.total))
		return
	}

	return
}

func (shared *shared) write(s int32, v []byte) error {
	if int(s)+len(v) > math.MaxInt32 {
		return errCapLimit
	}

	if s > shared.off {
		return errStartPos
	}

	offset := s % chunkSize
	idx := s / chunkSize

	for i := 0; i < len(v); i++ {
		if int32(len(shared.chunks)) == idx {
			shared.chunks = append(shared.chunks, getChunk())
		}

		chunk := shared.chunks[idx]
		i += copy(chunk[offset:], v[i:])
		offset = 0
		idx++
	}

	if s+int32(len(v)) > shared.off {
		shared.off = s + int32(len(v))
	}
	return nil
}

func (shared *shared) read(s, e int32, dst []byte) ([]byte, error) {
	if s > e || s < 0 || e > shared.off {
		return nil, errStartPos
	}

	length := int32(e - s)

	if dst == nil {
		dst = make([]byte, length)
	}

	if int32(len(dst)) < length {
		length = int32(len(dst))
	}

	offset := s % chunkSize
	idx := s / chunkSize

	for i := 0; i < int(length); {
		chunk := shared.chunks[idx]
		i += copy(dst[i:], chunk[offset:])
		offset = 0
		idx++
	}

	return dst[:length], nil
}

func (shared *shared) recycle() error {
	shared.mu.Lock()
	defer shared.mu.Unlock()

	blocks := make([]block, len(shared.keys))
	for _, block := range shared.keys {
		blocks = append(blocks, block)
	}

	rs := recycleSorts(blocks)
	sort.Sort(rs)

	s := int32(0)
	for _, block := range blocks {
		key, _ := shared.read(block.s, block.s+int32(block.kl), nil)
		value, _ := shared.read(block.s+int32(block.kl), block.s+int32(block.kl)+int32(block.vl), nil)
		block.s = s
		shared.write(s, key)
		shared.write(s+int32(block.kl), value)
		block.total = uint16(len(key) + len(value))
		s += int32(block.total)
		shared.keys[defaultHasher.Sum64(slice2string(key))] = block
	}

	shared.off = s
	idx := int(s / chunkSize)
	for i := idx + 1; i < len(shared.chunks); i++ {
		chunk := shared.chunks[i]
		putChunk(chunk)
		shared.chunks[i] = nil
	}

	if idx < len(shared.chunks) {
		shared.chunks = shared.chunks[:idx+1]
	}

	atomic.StoreInt64(&shared.stat.removeBytes, 0)
	atomic.StoreInt64(&shared.stat.totalBytes, int64(s))

	return nil
}

type recycleSorts []block

func (rs recycleSorts) Len() int {
	return len(rs)
}

func (rs recycleSorts) Less(i, j int) bool {
	return rs[i].s < rs[j].s
}

func (rs recycleSorts) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}
