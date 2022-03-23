package tcache

import (
	"fmt"
	"math/rand"
	"strconv"
	"tcache"
	"testing"

	"github.com/VictoriaMetrics/fastcache"
)

const (
	maxBytes = 1 << 31
)

func fastcacheStat(c *fastcache.Cache) {
	stat := fastcache.Stats{}
	c.UpdateStats(&stat)
	fmt.Printf("%+v\n", stat)
}

func initKeys(n int, size int) map[string][]byte {
	var total int
	keys := make(map[string][]byte, n)
	if size == 0 {
		size = 1
	}

	for i := 0; i < n; i++ {
		key := strconv.FormatInt(int64(i), 10)
		keys[key] = make([]byte, size)
		total += len(key)
		total += size
	}

	println("initKeys total:", total, " n:", n)

	return keys
}

func BenchmarkFastCacheWrite64(b *testing.B) {
	keys := initKeys(b.N, 64)
	c := fastcache.New(maxBytes)
	defer c.Reset()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(tcache.S2B(k), keys[k])
		}
	})
	fastcacheStat(c)
}

func BenchmarkFastCacheWrite256(b *testing.B) {
	keys := initKeys(b.N, 256)
	c := fastcache.New(maxBytes)
	defer c.Reset()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(tcache.S2B(k), keys[k])
		}
	})
	fastcacheStat(c)
}

func BenchmarkFastCacheWrite512(b *testing.B) {
	keys := initKeys(b.N, 512)
	c := fastcache.New(maxBytes)
	defer c.Reset()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(tcache.S2B(k), keys[k])
		}
	})
	fastcacheStat(c)
}

func BenchmarkFastCacheWrite1024(b *testing.B) {
	keys := initKeys(b.N, 1024)
	c := fastcache.New(maxBytes)
	defer c.Reset()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(tcache.S2B(k), keys[k])
		}
	})
	fastcacheStat(c)
}

func BenchmarkFastCacheWrite4196(b *testing.B) {
	keys := initKeys(b.N, 4196)
	c := fastcache.New(maxBytes)
	defer c.Reset()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(tcache.S2B(k), keys[k])
		}
	})
	fastcacheStat(c)
}

func BenchmarkFastCacheWrite10240(b *testing.B) {
	keys := initKeys(b.N, 10240)
	c := fastcache.New(maxBytes)
	defer c.Reset()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(tcache.S2B(k), keys[k])
		}
	})
	fastcacheStat(c)
}

func BenchmarkFastCacheWrite30000(b *testing.B) {
	keys := initKeys(b.N, 30000)
	c := fastcache.New(maxBytes)
	defer c.Reset()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(tcache.S2B(k), keys[k])
		}
	})
	fastcacheStat(c)
}

func BenchmarkCacheWrite128(b *testing.B) {
	keys := initKeys(b.N, 128)
	c := tcache.New(tcache.WithShared(512))

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(k, keys[k])
		}
	})
	c.Debug()
}

func BenchmarkCacheWrite1024(b *testing.B) {
	keys := initKeys(b.N, 1024)
	c := tcache.New(tcache.WithShared(512))

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(k, keys[k])
		}
	})
	c.Debug()
}

func BenchmarkCacheWrite4196(b *testing.B) {
	keys := initKeys(b.N, 4196)
	c := tcache.New(tcache.WithShared(512))

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(k, keys[k])
		}
	})

	c.Debug()
}

func BenchmarkCacheWrite10240(b *testing.B) {
	keys := initKeys(b.N, 10240)
	c := tcache.New(tcache.WithShared(512))

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			c.Set(k, keys[k])
		}
	})
	c.Debug()
}

func BenchmarkCacheWrite30000(b *testing.B) {
	keys := initKeys(b.N, 30000)
	c := tcache.New(tcache.WithShared(2048))

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := strconv.FormatInt(int64(rand.Intn(b.N)), 10)
			err := c.Set(k, keys[k])
			if err != nil {
				panic(err.Error())
			}
		}
	})
	c.Debug()
}
