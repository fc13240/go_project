package main

import "fmt"

func sum(a []int, c chan int) {
	sum := 0
	for _, v := range a {
		sum += v
	}
	fmt.Println(a)
	c <- sum // send sum to c
}
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
		fmt.Println("i = ", i)
	}
	close(c)
}
func mai1n() {
	a := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(a[:len(a)/2], c)
	go sum(a[len(a)/2:], c)
	x, y := <-c, <-c // 从 c 中接收

	fmt.Println(x, y, x+y)
	fmt.Println("========")

	c1 := make(chan int, 30)
	go fibonacci(cap(c1), c1)
	for i := range c1 {
		fmt.Println(i)
	}
}
