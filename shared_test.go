package cache

import (
	"strconv"
	"testing"

	"github.com/cespare/xxhash/v2"
)

func TestSharedSet(t *testing.T) {
	shared := newShared(1)
	for i := 0; i < 1000; i++ {
		key := string2slice(strconv.FormatInt(int64(i), 10))
		hash := xxhash.Sum64(key)
		shared.set(hash, key, []byte{2, 2, 2})
	}
}

func TestSharedGet(t *testing.T) {
	shared := newShared(1)
	key := []byte("abcdefg")
	value := []byte("k22233ddd3213423243242343242342")
	hash := xxhash.Sum64(key)
	shared.set(hash, key, value)

	b, _ := shared.get(hash, slice2string(key))
	t.Logf("%s\n", slice2string(b))
	value = []byte("23424sdsfsfsdfdsfsfffffffffffffffff")
	hash = xxhash.Sum64(key)
	shared.set(hash, key, value)
	b, _ = shared.get(hash, slice2string(key))
	t.Logf("%s\n", slice2string(b))
}

func BenchmarkSharedSet(b *testing.B) {
	shared := newShared(1)

	const items = 1 << 20
	b.ReportAllocs()
	b.SetBytes(items)

	for i := 0; i < b.N; i++ {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyza")
		for n := 0; n < items; n++ {
			k[0]++
			if k[0] == 0 {
				k[1]++
				if k[1] == 0 {
					k[2]++
					if k[2] == 0 {
						k[3]++
					}
				}
			}
			hash := xxhash.Sum64(k)
			err := shared.set(hash, k, v)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	b.Log("len=", len(shared.collision))
}
