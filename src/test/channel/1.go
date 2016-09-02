package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("begin")
	ch := make(chan int)

	go func() {
		fmt.Println("inner func")
		time.Sleep(time.Second * 5)
		fmt.Println("after some seconds")
		ch <- 1
	}()

	<-ch
	fmt.Println("after func")
}
