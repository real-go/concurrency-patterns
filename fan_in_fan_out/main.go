package main

import (
	"fmt"
	"sync"
)

// 多个函数可以从同一个通道获得输入：扇出
// 一个函数可以从多个通道读取并继续：扇入

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		for num := range c {
			out <- num
		}
		wg.Done()
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, num := range nums {
			out <- num
		}
		close(out)
	}()

	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for num := range in {
			out <- num * num
		}
		close(out)
	}()
	return out
}

func main() {
	in := gen(1, 2, 3)
	c1 := sq(in)
	c2 := sq(in)
	for num := range merge(c1, c2) {
		fmt.Println(num)
	}
}
