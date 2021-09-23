package cache

import (
	"reflect"
	"unsafe"
)

const (
	max = 1 << 32
)

func power2(cap int64) int64 {
	cap = cap - 1
	cap |= cap >> 1
	cap |= cap >> 2
	cap |= cap >> 4
	cap |= cap >> 8
	cap |= cap >> 16

	if cap < 0 {
		println("xxx", cap)
		return max
	}

	if cap > max {
		return max
	}

	return cap + 1
}

func string2slice(k string) []byte {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&k))
	sh.Cap = sh.Len

	return *(*[]byte)(unsafe.Pointer(sh))
}

func string2slicev2(k string) []byte {
	return []byte(k)
}