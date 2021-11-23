package cache

import (
	"testing"
)

func TestAddBlock(t *testing.T) {
	sb := sortBlocks{}
	sb.add(block{
		si:    1,
		total: 10,
	})
	sb.add(block{
		si:    11,
		total: 20,
	})
	sb.add(block{
		si:    32,
		total: 2,
	})
	sb.add(block{
		si:    35,
		total: 1,
	})
}

func TestGetBlock(t *testing.T) {
	sb := sortBlocks{}
	sb.add(block{
		si:    1,
		total: 10,
	})
	sb.add(block{
		si:    11,
		total: 20,
	})
	sb.add(block{
		si:    32,
		total: 2,
	})
	sb.add(block{
		si:    35,
		total: 1,
	})
	b, ok := sb.getBlock(1)
	t.Log(b, ok)
	b, ok = sb.getBlock(1)
	t.Log(b, ok)
	b, ok = sb.getBlock(1)
	t.Log(b, ok)

}

func TestChunk(t *testing.T) {
	chunk := chunk{
		used: 1,
		kl:   1<<16 - 1,
		vl:   1<<31 - 1,
	}

	bytes := chunk.encode()

	newChunk := decodeChunk(bytes)
	t.Log(newChunk)
}

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

func BenchmarkDecodeEncode(b *testing.B) {
	var n int
	var bytes [4]byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encode(i, bytes[:])
		n = decode(bytes[:])
	}
	n++
}
