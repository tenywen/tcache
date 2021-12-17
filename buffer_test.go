package cache

import (
	"fmt"
	"testing"
)

func BenchmarkBufferEncode(b *testing.B) {
	buffer := newBuffer(1 << 30)
	key := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	value := []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	for i := 0; i < b.N; i++ {
		chunk := getChunk()
		chunk.used = ^unused
		chunk.s = 0
		chunk.kl = int16(len(key))
		chunk.vl = int32(len(value))
		chunk.total = int32(chunk.kl) + chunk.vl
		chunk.k = key
		chunk.v = value
		err := buffer.encode(chunk)
		if err != nil {
			panic(err.Error())
		}
		recycleChunk(chunk)
	}
}

func BenchmarkBufferEncodeDecode(b *testing.B) {
	buffer := newBuffer(1 << 30)
	key := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	value := []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")

	for i := 0; i < b.N; i++ {
		chunk := getChunk()
		chunk.used = ^unused
		chunk.s = 0
		chunk.kl = int16(len(key))
		chunk.vl = int32(len(value))
		chunk.total = int32(chunk.kl) + chunk.vl
		chunk.k = key
		chunk.v = value
		err := buffer.encode(chunk)
		if err != nil {
			panic(err.Error())
		}

		err = buffer.decode(chunk)
		if err != nil {
			panic(err.Error())
		}

		if slice2string(chunk.k) != slice2string(key) || slice2string(chunk.v) != slice2string(value) {
			panic(fmt.Errorf("got %q %q want %q %q", chunk.k, chunk.v, key, value))
		}

		recycleChunk(chunk)
	}
}
