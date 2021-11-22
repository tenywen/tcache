package cache

import "sort"

const (
	chunkBit  = 64 // bit
	keyBit    = 16 // bit
	valueBit  = chunkBit - keyBit
	maxBuffer = 1 << 40 // 256TB
)

type chunk struct {
	kl int
	vl int
}

type block struct {
	si    int
	total int
	kl    int
	vl    int
}

func decodeChunk(bytes [5]byte) (chunk chunk) {
	chunk.kl = int(bytes[0])<<8 | int(bytes[1])
	chunk.vl = int(bytes[2])<<16 | int(bytes[3])<<8 | int(bytes[4])
	return
}

func (chunk chunk) encode() [5]byte {
	return [5]byte{byte(chunk.kl >> 8), byte(chunk.kl), byte(chunk.vl >> 16), byte(chunk.vl >> 8), byte(chunk.vl)}
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
