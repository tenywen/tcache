package cache

import (
	"strconv"
	"testing"

	xxhash "github.com/cespare/xxhash/v2"
)

var (
	temp = "abcde"
)

func TestHash(t *testing.T) {
	total := 1000000000
	m := make(map[uint64]struct{}, total)
	collison := 0
	for i := 1; i < total; i++ {
		hash := xxhash.Sum64([]byte(strconv.FormatInt(int64(i), 10)))
		if _, ok := m[hash]; ok {
			collison++
			continue
		}

		m[hash] = struct{}{}
	}

	t.Log("collison:", collison)
}

func BenchmarkHash(b *testing.B) {
}
