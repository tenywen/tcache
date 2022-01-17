[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/tenywen/cache/blob/master/LICENSE)

# Cache - inmemory cache in Go

读写速度更快，更好适应并发的，尽可能避免gc的内置cache.


### benchmarks

使用fastcache测试用例，测试数据如下:

```
➜ go test -v -bench=. -benchmem  bench/fast_cache_test.go

goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkMyCacheSet
BenchmarkMyCacheSet-12               622           1624744 ns/op          40.34 MB/s      271685 B/op      65572 allocs/op
BenchmarkFastCacheSet
BenchmarkFastCacheSet-12             829           1340148 ns/op          48.90 MB/s        6877 B/op         14 allocs/op
BenchmarkFastCacheGet
BenchmarkFastCacheGet-12            1381            841510 ns/op          77.88 MB/s           3 B/op          0 allocs/op
BenchmarkMyCacheGet
BenchmarkMyCacheGet-12              1108           1081764 ns/op          60.58 MB/s           4 B/op          0 allocs/op
BenchmarkBigCacheSet
BenchmarkBigCacheSet-12              338           3885160 ns/op          16.87 MB/s     6281929 B/op         27 allocs/op
BenchmarkBigCacheGet
BenchmarkBigCacheGet-12              700           1624757 ns/op          40.34 MB/s      981064 B/op     131082 allocs/op
BenchmarkBigCacheSetGet
BenchmarkBigCacheSetGet-12           198           5472170 ns/op          23.95 MB/s     5189552 B/op     131113 allocs/op
BenchmarkStdMapSet
BenchmarkStdMapSet-12                120           9866489 ns/op           6.64 MB/s      366742 B/op      65555 allocs/op
BenchmarkStdMapGet
BenchmarkStdMapGet-12                424           2609145 ns/op          25.12 MB/s       30196 B/op        160 allocs/op
BenchmarkStdMapSetGet
BenchmarkStdMapSetGet-12              68          40845637 ns/op           3.21 MB/s      447027 B/op      65574 allocs/op
BenchmarkSyncMapSet
BenchmarkSyncMapSet-12                55          23539543 ns/op           2.78 MB/s     3577198 B/op     264569 allocs/op
BenchmarkSyncMapGet
BenchmarkSyncMapGet-12              1593            677392 ns/op          96.75 MB/s        7986 B/op        248 allocs/op
BenchmarkSyncMapSetGet
BenchmarkSyncMapSetGet-12            170           5912691 ns/op          22.17 MB/s     3462660 B/op     262928 allocs/op
PASS
ok      command-line-arguments  23.128s
```

