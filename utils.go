package cache

import (
	"unsafe"
)

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type stringHeader struct {
	Data unsafe.Pointer
	Len  int
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
		return max
	}

	if cap > max {
		return max
	}

	return cap + 1
}

func string2slice(k string) []byte {
	const max = 0x7fff0000
	if len(k) > max {
		panic("string too long")
	}
	return (*[max]byte)((*stringHeader)(unsafe.Pointer(&k)).Data)[:len(k):len(k)]
}

func slice2string(k []byte) string {
	return *(*string)(unsafe.Pointer(&k))
}
