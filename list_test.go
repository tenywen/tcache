package cache

import (
	"container/list"
	"testing"
)

func BenchmarkListInsert(b *testing.B) {
	list := list.New()
	n := list.PushFront(block{})
	for i := 0; i < b.N; i++ {
		block := block{}
		n = list.InsertAfter(block, n)
	}

	//b.Log("len:", list.Len())
}
