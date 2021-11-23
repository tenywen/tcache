package cache

import "sort"

const (
	keyBit    = 15 // bit
	valueBit  = 32
	chunkBit  = keyBit + valueBit + 1 // 1 = usedBit
	maxBuffer = 1 << 40               // 256TB
)

type chunk struct {
	used int
	kl   int
	vl   int
}

type block struct {
	si    int
	total int
	kl    int
	vl    int
}

func decodeChunk(bytes [chunkBit >> 3]byte) (chunk chunk) {
	chunk.used = decode(bytes[:1]) & 0x1
	chunk.kl = decode(bytes[:2]) >> 1
	chunk.vl = decode(bytes[2:])
	return
}

func (chunk chunk) encode() (bytes [chunkBit >> 3]byte) {
	encode(chunk.used, bytes[:1])
	encode(chunk.kl<<1, bytes[:2])
	encode(chunk.vl, bytes[2:])
	return
}

func encode(v int, bytes []byte) {
	for i := 0; i < len(bytes); i++ {
		if v == 0 {
			break
		}
		bytes[i] = byte(v)
		v >>= 8
	}
}

func decode(bytes []byte) (v int) {
	for i := len(bytes) - 1; i >= 0; i-- {
		v <<= 8
		v += int(bytes[i])
	}
	return
}

type sortBlocks []block

func (sb sortBlocks) Len() int {
	return len(sb)
}

func (sb sortBlocks) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}

func (sb sortBlocks) Less(i, j int) bool {
	return sb[i].total < sb[j].total
}

func (sb *sortBlocks) add(b block) {
	*sb = append(*sb, b)
}

func (sb *sortBlocks) getBlock(size int) (b block, ok bool) {
	sort.Sort(*sb)
	length := sb.Len()
	index := sort.Search(length, func(i int) bool {
		return (*sb)[i].total >= size
	})

	if index >= length || (*sb)[index].total < size {
		return
	}

	ok = true
	b = (*sb)[index]

	(*sb)[index] = (*sb)[length-1]
	(*sb) = (*sb)[:length-1]

	return
}
