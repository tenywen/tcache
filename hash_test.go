package tcache

import (
	"strconv"
	"testing"

	xxhash "github.com/cespare/xxhash/v2"
)

var (
	temp = "abcde"
)

func TestHash(t *testing.T) {
	total := 1000000
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

func TestHashDiff(t *testing.T) {
	k1 := []byte("\xc1\xee\x00\x00")
	k2 := []byte("\xff\x0d\x00\x00")
	t.Log(xxhash.Sum64(k1))
	t.Log(xxhash.Sum64(k2))
}
