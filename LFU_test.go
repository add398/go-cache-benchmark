/**
 * @Author: Log1c
 * @Description:
 * @File:  testlfu
 * @Version: 1.0.0
 * @Date: 2022/4/18 16:57
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

const size = 500

func TestLFU_log1c(t *testing.T)  {

	cache := NewTinyLFU(size)
	chan1 := make(chan string, 5)
	sum := 1000000
	fmt.Println(sum)

	go func() {
		fi, err := os.Open("log1ctest/from_to.txt")

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
		cache.Set(value)
	}
	for  {
		count++
		value := <- chan1
		if cache.Get(value) {
			hits++
		} else {
			misses++
			cache.Set(value)
		}

		if count == int(sum) {
			fmt.Println(count, hits, misses, float64(float64(hits) / float64(count)))
			return
		}
	}
}


func TestLRU_log1c(t *testing.T)  {
	cache := NewGroupCacheLRU(size)
	chan1 := make(chan string, 5)
	sum := 1000e3
	fmt.Println(sum)
	go func() {
		fi, err := os.Open("log1ctest/from_to.txt")

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
		cache.Set(value)
	}
	for  {
		count++
		value := <- chan1
		if cache.Get(value) {
			hits++
		} else {
			misses++
			cache.Set(value)
		}

		if count == int(sum) {
			fmt.Println(count, hits, misses, float64(float64(hits) / float64(count)))
			return
		}
	}

}

