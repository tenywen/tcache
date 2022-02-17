package tcache

import (
	"fmt"
	"testing"

	"github.com/VictoriaMetrics/fastcache"
)

const items = 1 << 16

func BenchmarkMyCacheSet(b *testing.B) {
	cache := New(WithShared(2048))
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyz")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				err := cache.Set(slice2string(k), v)
				if err != nil {
					panic(err.Error())
				}
			}
		}
	})
	//cache.Debug()
}

func BenchmarkFastCacheSet(b *testing.B) {
	c := fastcache.New(12 * items)
	defer c.Reset()
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyz")
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				c.Set(k, v)
			}
		}
	})
	/*
		stat := fastcache.Stats{}
		c.UpdateStats(&stat)
		b.Logf("%+v \n", stat)
	*/
}

func BenchmarkFastCacheGet(b *testing.B) {
	c := fastcache.New(12 * items << 10)
	defer c.Reset()
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xyza")
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		c.Set(k, v)
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		var buf []byte
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				buf = c.Get(buf[:0], k)
				if slice2string(buf) != slice2string(v) {
					panic(fmt.Errorf("BUG: invalid value obtained; got %q; want %q", buf, v))
				}
			}
		}
	})
}

func BenchmarkMyCacheGet(b *testing.B) {
	cache := New(WithShared(1024))
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xyza")
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		cache.Set(slice2string(k), v)
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		var buf []byte
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}

				buf, _ = cache.Get(slice2string(k), buf)
				if slice2string(buf) != slice2string(v) {
					panic(fmt.Errorf("BUG: key:%q want:%s got:%s", k, string(v), string(buf)))
				}
			}
		}
	})
}
