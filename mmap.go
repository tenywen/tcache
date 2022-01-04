package cache

const (
	chunkSize = 1 << 16
)

func getChunk() []byte {
	return make([]byte, chunkSize)
}

func putChunk(chunk []byte) {
	if chunk == nil {
		return
	}
}
