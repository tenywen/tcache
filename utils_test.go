package tcache

import (
	"testing"
)

var bytes []byte

func TestArrayToSlice(t *testing.T) {
	var a = [5]int{1, 2, 3, 4, 5}
	b := (&a)[:2:2]
	b[0] = 99
	t.Log("a=", a)
	t.Log("b=", b)
}

func Benchmark2String(b *testing.B) {
	s := "tenywen"
	for i := 0; i < b.N; i++ {
		bytes = string2slice(s)
	}
}
