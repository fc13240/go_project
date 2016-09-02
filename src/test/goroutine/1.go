package main

import "fmt"

func add(a, b int, ch chan int) {
	sum := a + b
	fmt.Printf("%d + %d = %d\n", a, b, sum);

	ch <- sum
}
func main() {
	ch := make(chan int)
	// add(1, 2, ch);
	fmt.Println("---");
	go add(10, 11, ch);
	fmt.Println("---2");

	sum := <- ch
	fmt.Printf("sum = %d", sum);
}