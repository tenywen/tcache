package cache

import (
	"fmt"
	"testing"
)

func BenchmarkMyCacheGet(b *testing.B) {
	const items = 1 << 16
	cache := New(WithShared(512), WithMaxBuffer(1<<24))
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		cache.Set(k, v)
	}

	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				buf, _ := cache.Get(slice2string(k))
				if slice2string(buf) != slice2string(v) {
					panic(fmt.Errorf("BUG: key:%q got %q want:%q ", k, buf, v))
				}
			}
		}
	})
}

func BenchmarkMyCacheSet(b *testing.B) {
	const items = 1 << 16
	cache := New(WithShared(512), WithMaxBuffer(5))
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
				//add(k)
				err := cache.Set(k, v)
				if err != nil {
					panic(err.Error())
				}
			}
			//b.Log("len:", cache.shareds[0].buffer.off, "cap:", cap(cache.shareds[0].buffer.bytes))
		}
	})
}
