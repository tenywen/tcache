package cache

import (
	"reflect"
	"unsafe"
)

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type stringHeader struct {
}

const (
	max = 1 << 32
)

func power2(cap int) int {
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

func slice2string(k []byte) string {
	return *(*string)(unsafe.Pointer(&k))
}
