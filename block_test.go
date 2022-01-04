package cache

import (
	"math/rand"
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

func TestBlock(t *testing.T) {
	buffer := newBuffer()
	block := block{
		kl:    1,
		vl:    4,
		total: 99,
	}

	t.Log(putBlock(block, &buffer))
	nblock, err := getBlock(0, &buffer)
	t.Logf("%+v %v\n", nblock, err)
}

func BenchmarkBlockEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		kv := [headLen]byte{}
		encode(i, kv[:])
	}
}

func BenchmarkBlockDecode(b *testing.B) {
	kv := [headLen]byte{}
	encode(rand.Int(), kv[:])
	b.ResetTimer()
	v := 0
	for i := 0; i < b.N; i++ {
		v = decode(kv[:])
		v++
	}
	v = 0
}
