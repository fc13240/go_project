package main

import "fmt"
import "./test"

func main() {
	str := "hello wrold"
	fmt.Printf("str = %s\n", str)

	test.Reverse(&str)

	fmt.Printf("str = %s\n", str)

	str = "我们是中国人"
	fmt.Printf("str = %s\n", str)

	test.Reverse(&str)

	fmt.Printf("str = %s\n", str)
}
