/**
 * @Author: Log1c
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2022/4/18 22:33
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

type Benchmark struct {
	Generator
	N int
}


func main() {

	//cacheSize := []int{1e3, 10e3, 100e3, 1e6}
	cacheSize := []int{256, 512, 1000, 10000}
	multiplier := []int{10, 100, 1000}
	newCache := []NewCacheFunc{
		NewTinyLFU,
		NewClockPro,
		NewARC,
		NewRistretto,
		NewTwoQueue,
		NewGroupCacheLRU,
		NewHashicorpLRU,
		NewS4LRU,
		NewSLRU,
	}
	newGen := []NewGeneratorFunc{
		NewScrambledZipfian,
		// NewHotspot,
		// NewUniform,
	}

	var results []*BenchmarkResult

	for _, newGen := range newGen {
		for _, cacheSize := range cacheSize {
			for _, multiplier := range multiplier {
				numKey := cacheSize * multiplier

				if len(results) > 0 {
					printResults(results)
					results = results[:0]
				}

				for _, newCache := range newCache {
					result := run(newGen, cacheSize, numKey, newCache)
					results = append(results, result)
				}
			}
		}
	}
}

func run(newGen NewGeneratorFunc, cacheSize, numKey int, newCache NewCacheFunc) *BenchmarkResult {
	gen := newGen(numKey)
	b := &Benchmark{
		Generator: gen,
		N:         1e6,
	}

	alloc1 := memAlloc()
	cache := newCache(cacheSize)
	defer cache.Close()

	start := time.Now()
	hits, misses := bench(b, cache)
	dur := time.Since(start)

	alloc2 := memAlloc()

	return &BenchmarkResult{
		GenName:   gen.Name(),
		CacheName: cache.Name(),
		CacheSize: cacheSize,
		NumKey:    numKey,

		Hits:     hits,
		Misses:   misses,
		Duration: dur,
		Bytes:    int64(alloc2) - int64(alloc1),
	}
}

func bench(b *Benchmark, cache Cache) (hits, misses int) {

	chan1 := make(chan string, 5)

	go func() {
		fi, err := os.Open("/Users/log1c/Code/Golang/go-cache-benchmark/log1ctest/from_to.txt")

		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		defer fi.Close()

		br := bufio.NewReader(fi)
		i:=0
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			if i==0{
				// 过滤第一行 from，to
				i++
				continue
			}
			s1 := a[:42]
			s2:= a[43:]
			//fmt.Println(s1)
			//fmt.Println(s2)
			chan1 <- string(s1)
			chan1 <- string(s2)

		}
	}()

	for i := 0; i < b.N; i++ {
		//value := b.Next()
		value := <- chan1
		if cache.Get(value) {
			hits++
		} else {
			misses++
			cache.Set(value)
		}
	}

	return hits, misses
}

func memAlloc() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

