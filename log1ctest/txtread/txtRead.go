/**
 * @Author: Log1c
 * @Description:
 * @File:  txtRead
 * @Version: 1.0.0
 * @Date: 2022/1/19 11:56
 */


package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	ans := 0
	//fi, err := os.Open("log1c/from.txt")
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
		if ans % 1000 == 0 {
			fmt.Println(string(s1))
			fmt.Println(string(s2))
		}
		ans++
		ans++


	}
	fmt.Println(ans)
}
