package main

import (
	"fmt"
	// "strconv"
)

// http://www.tuicool.com/articles/JNRBvu
func main() {
	ch := make(chan int)
	task("A", ch)
	task("B", ch)
	fmt.Printf("begin\n")
	<-ch
	<-ch
}
func task(name string, ch chan int) {
	go func() {
		i := 1
		for {
			fmt.Printf("%s %d\n", name, i)
			i++
		}
		ch <- 1
	}()
}
