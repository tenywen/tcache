package tcache

import (
	"fmt"
	"log"
	"sync/atomic"
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
	calls        int64 // hit
	missCnt      int64 // miss
	collisionCnt int64
	removeBytes  int64
	totalBytes   int64
}

func (stat *stat) call() {
	atomic.AddInt64(&stat.calls, 1)
}

func (stat *stat) miss() {
	atomic.AddInt64(&stat.missCnt, 1)
}

func (stat *stat) collision() {
	atomic.AddInt64(&stat.collisionCnt, 1)
}

func (stat *stat) remove(bytes int64) {
	atomic.AddInt64(&stat.removeBytes, bytes)
}

func (stat *stat) add(bytes int64) {
	atomic.AddInt64(&stat.totalBytes, bytes)
}

func (stat *stat) debug() {
	calls := atomic.LoadInt64(&stat.calls)
	miss := atomic.LoadInt64(&stat.missCnt)
	removes := atomic.LoadInt64(&stat.removeBytes)
	totals := atomic.LoadInt64(&stat.totalBytes)
	log.Println("cache shared stat debug")
	log.Printf("call:%d\n", calls)
	log.Printf("miss:%d  %.2f\n", miss, float64(miss)/float64(calls))
	log.Printf("remove:%s\n", printBytes(removes))
	log.Printf("total:%s\n", printBytes(totals))
	log.Println("")
}

func printBytes(size int64) string {
	var index int
	v := float64(size)
	const base = 1 << 10 // 1k
	for v/base > 1 {
		index++
		v /= base
	}

	return fmt.Sprintf("%.2f%s", v, format[index])
}
