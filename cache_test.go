package cache

import (
	"fmt"
	"testing"
)

func TestMyCacheSet(t *testing.T) {
	cache := New(WithShared(1024))
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xyzx")
	//k[0]++
	err := cache.Set(slice2string(k), v)
	if err != nil {
		t.Log(err)
	}
	//k[0]++
	err = cache.Set(slice2string(k), v)
	if err != nil {
		t.Log(err)
	}
}

func TestMyCacheGet(t *testing.T) {
	c := New(WithShared(1024))
	k := []byte("\x00\x00\x00\x00")
	v := []byte("xyzx")
	err := c.Set(string(k), v)
	t.Logf("set %s %s err:%v", string(k), string(v), err)
	if err != nil {
		return
	}

	result, err := c.Get(string(k), nil)
	t.Logf("get %s %s err:%v", string(k), string(result), err)
	if err != nil {
		return
	}

	if string(result) != string(v) {
		panic(fmt.Errorf("get key:%q  want:%s got:%s", k, string(v), string(result)))
	}
}
