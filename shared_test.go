package cache

import "testing"

var buf = buffer{
	ids: make(map[int64]int64),
}

func TestWrite(t *testing.T) {
	buf.write([]byte("this is a key !"), []byte("this is a value !@#"))
}

func TestRead(t *testing.T) {
	buf.write([]byte("this is a key !"), []byte("this is a value !@#"))
	data, err := buf.getData(1)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(string(data.body.k), string(data.body.v))
}
