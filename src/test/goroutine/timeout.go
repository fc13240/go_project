package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("begin")
	timeout := make(chan bool)

	go func() {
		time.Sleep(time.Second * 2)
		timeout <- true
	}()
	ch := make(chan int)

	select {
	case <-ch:
		fmt.Println("read from chan")
	case <-timeout:
		fmt.Println("timeout")
	}

	fmt.Println("the 1")
	tc := time.Tick(time.Second) //返回一个time.C这个管道，1秒(time.Second)后会在此管道中放入一个时间点，
	//1秒后再放一个，一直反复，时间点记录的是放入管道那一刻的时间
	for i := 1; i <= 10; i++ {
		<-tc
		fmt.Println("hello")
	}
	//每隔1秒，打印一个hello
}
