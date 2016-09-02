package main

import "fmt"

func fn(argv ...interface{}) {
	fmt.Println(argv)

	for _, v := range argv {
		fmt.Println(v)
	}
}
func main() {
	fn(1)

	fn(1, 2, 3)

	fn([]int{1, 2, 3})
}
