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

type buffer struct {
	id int64
	// id -> start index
	ids map[int64]int64

	// hash -> id
	keys  map[uint64]int64
	bytes bytes
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

func (b buffer) get(key uint64) ([]byte, error) {
	id, ok := b.keys[key]
	if !ok {
		return nil, errKeyNotExist
	}

	data, err := b.getData(id)
	if err != nil {
		return nil, err
	}
	return data.body.v, nil
}

func (b *buffer) set(k string, v []byte, hasher Hasher) {
	hash := hasher.Sum64(k)
	id, ok := b.keys[hash]
	if ok {
		// TODO
		b.del(id)
		return
	}

	b.write(string2slice(k), v)
}

func (b *buffer) del(id int64) {

}

func (b buffer) getData(id int64) (data data, err error) {
	var si int64
	si, ok := b.ids[id]
	if !ok {
		err = errIdNotExist
		return
	}
	return b.decode(si)
}

func (b buffer) decode(si int64) (data data, err error) {
	data.head.id = b.bytes.int64(si, si+idLen)
	data.head.kl = b.bytes.int64(si+idLen, si+idLen+keyLen)
	data.head.vl = b.bytes.int64(si+idLen+keyLen, si+idLen+keyLen+valueLen)

	startBody := si + headLen
	data.body.k, err = b.bytes.read(startBody, startBody+data.head.kl)
	if err != nil {
		return
	}
	data.body.v, err = b.bytes.read(startBody+data.head.kl, startBody+data.head.kl+data.head.vl)
	return
}

func (b *buffer) write(k, v []byte) {
	b.id++
	id := b.id
	b.ids[id] = b.bytes.tail()
	b.bytes.putInt64(id)
	b.bytes.putInt64(int64(len(k)))
	b.bytes.putInt64(int64(len(v)))
	b.bytes.write(k)
	b.bytes.write(v)
}
