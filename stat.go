package cache

import (
	"fmt"
	"log"
)

var (
	format []string = []string{
		"B",
		"KB",
		"MB",
		"GB",
		"TB",
	}
)

type stat struct {
	hitCnt      int // hit
	missCnt     int //miss
	removeBytes int
	totalBytes  int
}

func (stat *stat) hit() {
	stat.hitCnt++
}

func (stat *stat) miss() {
	stat.missCnt++
}

func (stat *stat) remove(bytes int) {
	stat.removeBytes += bytes
}

func (stat *stat) add(bytes int) {
	stat.totalBytes += bytes
}

func (stat *stat) debug() {
	log.Println("cache stat debug")
	log.Printf("hit:%d\n", stat.hitCnt)
	log.Printf("miss:%d  %.2f\n", stat.missCnt, float64(stat.missCnt)/float64(stat.hitCnt))
	log.Printf("remove:%s\n", printBytes(stat.removeBytes))
	log.Printf("total:%s\n", printBytes(stat.totalBytes))
	log.Println("")
}

func printBytes(size int) string {
	var index int
	const base = 1 << 10 // 1k
	for size >= base {
		index++
		size /= base
	}

	return fmt.Sprintf("%d%s", size, format[index])
}
