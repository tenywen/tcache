package cache

const (
	unused    = 0
	usedBit   = 1 // 1B
	keyBit    = 2 // 64KB
	valueBit  = 4 // 16MB
	totalBit  = 4
	chunkBit  = usedBit + totalBit + keyBit + valueBit
	maxBuffer = 1 << 40 // 256TB
)

type chunk struct {
	used  int8
	total int32
	kl    int16
	vl    int32
	s     int
	k     []byte
	v     []byte
}

type block struct {
	s     int
	total int
	kl    int
	vl    int
	k     []byte
	v     []byte
}

func getBlock(s int, b *buffer) (block block, err error) {
	kv := make([]byte, headLen)
	kv, err = b.read(s, s+headLen, kv)
	if err != nil {
		return
	}

	block.s = s
	block.total = decode(kv[:totalLen])
	block.kl = decode(kv[totalLen : totalLen+keyLen])
	block.vl = decode(kv[totalLen+keyLen:])
	return
}

func putBlock(block block, b *buffer) error {
	kv := [headLen]byte{}
	encode(block.total, kv[:totalLen])
	encode(block.kl, kv[totalLen:totalLen+keyLen])
	encode(block.vl, kv[totalLen+keyLen:])
	return b.write(block.s, kv[:])
}

func (chunk *chunk) decode(bytes []byte) error {
	var s int
	chunk.used = int8(decode(bytes[s : s+usedBit]))
	if chunk.used != ^unused {
		return errNotContent
	}

	s += usedBit
	chunk.total = int32(decode(bytes[s : s+totalBit]))

	s += totalBit
	chunk.kl = int16(decode(bytes[s : s+keyBit]))

	s += keyBit
	chunk.vl = int32(decode(bytes[s : s+valueBit]))

	return nil
}

func (chunk *chunk) encode(bytes []byte) {
	var s int
	if chunk.used != 0 {
		encode(int(chunk.used), bytes[s:usedBit])
	}
	s += usedBit
	if chunk.total != 0 {
		encode(int(chunk.total), bytes[s:s+totalBit])
	}
	s += totalBit
	if chunk.kl != 0 {
		encode(int(chunk.kl), bytes[s:s+keyBit])
	}
	s += keyBit
	if chunk.vl != 0 {
		encode(int(chunk.vl), bytes[s:s+valueBit])
	}
}

func encode(v int, bytes []byte) int {
	for i := 0; i < len(bytes); i++ {
		if v == 0 {
			return i
		}
		bytes[i] = byte(v)
		v >>= 8
	}
	return 0
}

func decode(bytes []byte) (v int) {
	for i := len(bytes) - 1; i >= 0; i-- {
		v <<= 8
		v += int(bytes[i])
	}
	return
}

/*
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

func (sb *sortBlocks) add(chunk *chunk) {
	*sb = append(*sb, block{s: chunk.s, total: chunk.total})
}

func (sb *sortBlocks) getBlock(size int32, chunk *chunk) bool {
	if len(*sb) == 0 {
		return false
	}

	sort.Sort(*sb)
	length := sb.Len()
	index := sort.Search(length, func(i int) bool {
		return (*sb)[i].total >= size
	})

	if index >= length || (*sb)[index].total < size {
		return false
	}

	chunk.total = int32((*sb)[index].total)

	(*sb)[index] = (*sb)[length-1]
	(*sb) = (*sb)[:length-1]

	return true
}
*/
