package cache

import (
	"testing"
)

func TestEncode(t *testing.T) {
	var b [4]byte
	encode(16, b[:])
	t.Log(b)
	encode(15, b[:])
	t.Log(b)
	encode(13, b[:])
	t.Log(b)
	encode(255, b[:])
	t.Log(b)
	encode(256, b[:])
	t.Log(b)

}

func TestDecodeEncode(t *testing.T) {
	var b [5]byte
	encode(16, b[:])
	t.Log(decode(b[:]))
	encode(64, b[:])
	t.Log(decode(b[:]))
	encode(128, b[:])
	t.Log(decode(b[:]))
	encode(130, b[:])
	t.Log(decode(b[:]))
}

func BenchmarkGetPutChunk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chunk := getChunk()
		chunk.kl = 1
		recycleChunk(chunk)
	}
}

func BenchmarkChunkEncode(b *testing.B) {
	bytes := make([]byte, chunkBit)
	b.ResetTimer()
	k := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	v := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	for i := 0; i < b.N; i++ {
		chunk := getChunk()
		chunk.s = 0
		chunk.k = k
		chunk.v = v
		chunk.kl = int16(len(k))
		chunk.vl = int32(len(v))
		chunk.encode(bytes)
		recycleChunk(chunk)
		bytes[0]++
	}
}
