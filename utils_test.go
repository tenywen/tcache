package cache

import (
	"math"
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

func TestArrayToSlice(t *testing.T) {
	var a = [5]int{1, 2, 3, 4, 5}
	b := (&a)[:2:2]
	b[0] = 99
	t.Log("a=", a)
	t.Log("b=", b)
}

func Benchmark2SliceV1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		string2slice("123456")
	}
}
