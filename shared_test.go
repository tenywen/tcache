package cache

import (
	"strconv"
	"testing"
	"time"
)

func TestSharedRecycle(t *testing.T) {
	shared := newShared(defaultOpt())
	for i := int64(0); i < 1<<20; i++ {
		key := strconv.FormatInt(i, 10)
		shared.set(false, defaultHasher.Sum64(key), key, string2slice(key))
	}

	for i := int64(1); i < 10000; i++ {
		key := strconv.FormatInt(i, 10)
		shared.delete(defaultHasher.Sum64(key), key)
	}

	start := time.Now()

	shared.recycle()

	key := "1111"
	v, err := shared.get(false, defaultHasher.Sum64(key), key)
	t.Logf("k=%s v=%s err=%v ts=%v\n", key, slice2string(v), err, time.Now().Sub(start))

	key = "10001"
	v, err = shared.get(false, defaultHasher.Sum64(key), key)
	t.Logf("k=%s v=%s err=%v ts=%v\n", key, v, err, time.Now().Sub(start))
}
