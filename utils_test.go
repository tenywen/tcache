package cache

import (
	"testing"
)

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
