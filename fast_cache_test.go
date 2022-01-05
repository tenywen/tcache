package cache

import (
	"fmt"
	"testing"

	"github.com/VictoriaMetrics/fastcache"
)

func BenchmarkFastCacheSet(b *testing.B) {
	const items = 1 << 16
	c := fastcache.New(12 * items)
	defer c.Reset()
	b.ReportAllocs()
	b.SetBytes(items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
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
}

func BenchmarkFastCacheGet(b *testing.B) {
	const items = 1 << 16
	c := fastcache.New(12 * items)
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
				buf = c.Get(buf[:0], k)
				if string(buf) != string(v) {
					panic(fmt.Errorf("BUG: invalid value obtained; got %q; want %q", buf, v))
				}
			}
		}
	})
}

func BenchmarkFastCacheHas(b *testing.B) {
	const items = 1 << 16
	c := fastcache.New(12 * items)
	defer c.Reset()
	k := []byte("\x00\x00\x00\x00")
	for i := 0; i < items; i++ {
		k[0]++
		if k[0] == 0 {
			k[1]++
		}
		c.Set(k, nil)
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
				if !c.Has(k) {
					panic(fmt.Errorf("BUG: missing value for key %q", k))
				}
			}
		}
	})
}

/*
func BenchmarkFastCacheSetGet(b *testing.B) {
	const items = 1 << 16
	c := fastcache.New(12 * items)
	defer c.Reset()
	b.ReportAllocs()
	b.SetBytes(2 * items)
	b.RunParallel(func(pb *testing.PB) {
		k := []byte("\x00\x00\x00\x00")
		v := []byte("xyza")
		var buf []byte
		for pb.Next() {
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				c.Set(k, v)
			}
			for i := 0; i < items; i++ {
				k[0]++
				if k[0] == 0 {
					k[1]++
				}
				buf = c.Get(buf[:0], k)
				if string(buf) != string(v) {
					panic(fmt.Errorf("BUG: invalid value obtained; got %q; want %q", buf, v))
				}
			}
		}
	})
}
*/
