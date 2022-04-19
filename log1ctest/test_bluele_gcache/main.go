/**
 * @Author: Log1c
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2022/4/18 17:48
 */

package main

import (
	"bufio"
	"fmt"
	"github.com/bluele/gcache"
	"io"
	"os"
)

func main() {
	size := 10000
	gc := gcache.New(size).
		LRU().
		Build()

	chan1 := make(chan string, 5)
	sum := 100e4
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
		gc.Set(value,"1")
	}
	for  {
		count++
		if count % 100000 == 0 {
			fmt.Println(count)
		}
		value := <- chan1
		if _, err := gc.Get(value); err == nil {
			hits++
		} else {
			misses++
			gc.Set(value,"1")
		}

		if count == int(sum) {
			fmt.Println(count, hits, misses, float64(float64(hits) / float64(count)))
			return
		}
	}


}


