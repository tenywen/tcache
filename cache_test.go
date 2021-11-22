package cache

import (
	"fmt"
	"testing"
)

func BenchmarkMyCacheGet(b *testing.B) {
	const items = 1 << 16
	cache := New(WithShared(512), WithMaxBuffer(1<<30))
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xyza")
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		cache.Set(slice2string(k), v)
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		var buf []byte
		k := []byte("\x00\x00\x00\x00")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				buf, _ = cache.Get(slice2string(k))
				if slice2string(buf) != slice2string(v) {
					panic(fmt.Errorf("BUG: got %q want:%q", buf, v))
				}
			}
		}
	})

}

func BenchmarkMyCacheSet(b *testing.B) {
	const items = 1 << 16
	cache := New(WithShared(512), WithMaxBuffer(1<<30))
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyza")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				cache.Set(slice2string(k), v)
			}
		}
	})
}
