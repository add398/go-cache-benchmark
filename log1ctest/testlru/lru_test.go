/**
 * @Author: Log1c
 * @Description:
 * @File:  lru_test
 * @Version: 1.0.0
 * @Date: 2022/4/18 17:34
 */

package simplelru

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestLRU(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k interface{}, v interface{}) {
		if k != v {
			t.Fatalf("Evict values not equal (%v!=%v)", k, v)
		}
		evictCounter++
	}
	l, err := NewLRU(128, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}

	if evictCounter != 128 {
		t.Fatalf("bad evict count: %v", evictCounter)
	}

	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || v != k || v != i+128 {
			t.Fatalf("bad key: %v", k)
		}
	}
	for i := 0; i < 128; i++ {
		_, ok := l.Get(i)
		if ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := l.Get(i)
		if !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		ok := l.Remove(i)
		if !ok {
			t.Fatalf("should be contained")
		}
		ok = l.Remove(i)
		if ok {
			t.Fatalf("should not be contained")
		}
		_, ok = l.Get(i)
		if ok {
			t.Fatalf("should be deleted")
		}
	}

	l.Get(192) // expect 192 to be last key in l.Keys()

	for i, k := range l.Keys() {
		if (i < 63 && k != i+193) || (i == 63 && k != 192) {
			t.Fatalf("out of order key: %v", k)
		}
	}

	l.Purge()
	if l.Len() != 0 {
		t.Fatalf("bad len: %v", l.Len())
	}
	if _, ok := l.Get(200); ok {
		t.Fatalf("should contain nothing")
	}
}



func TestLRU_log1c(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k interface{}, v interface{}) {
		if k != v {
			t.Fatalf("Evict values not equal (%v!=%v)", k, v)
		}
		evictCounter++
	}


	size := 100000

	l, err := NewLRU(size, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	chan1 := make(chan string, 5)
	sum := 100e5
	fmt.Println(sum)

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

	hits := 0
	misses := 0
	count := 0
	for i := 0; i < size; i++ {
		value := <- chan1
		l.Add(value,value)
	}
	for  {
		count++
		value := <- chan1
		if _, ok := l.Get(value); ok {
			hits++
		} else {
			misses++
			l.Add(value,value)
		}

		if count == int(sum) {
			fmt.Println(count, hits, misses, float64(float64(hits) / float64(count)))
			return
		}
	}




}
