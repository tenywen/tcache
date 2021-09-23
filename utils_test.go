package cache

import (
	"math"
	"strconv"
	"testing"
)

func TestPower2(t *testing.T) {
	t.Log(power2(-1))
	t.Log(power2(2))
	t.Log(power2(3))
	t.Log(power2(4))
	t.Log(power2(5))
	t.Log(power2(118))
	t.Log(power2(128))
	t.Log(power2(1023))
	t.Log(power2(1024))
	t.Log(power2(1025))
	t.Log(power2(math.MaxInt64))
}

func Test2Slice(t *testing.T) {
	t.Log(string(string2slice("test is a test")))
}

func Benchmark2Slice(b *testing.B) {
	key := "test is a test"
	for i := 0; i < b.N; i++ {
		string2slice(key + strconv.FormatInt(int64(0), 10))
	}

}
