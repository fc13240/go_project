package main

import "fmt"
import "./test"

/*
package以目录为单位，一个目录下的GO文件中package一般与目录一致，
同个package下的文件中方法在另个一个文件中可不用import直接使用，
但在package之外用import引用后方可使用

同个package下的文件编译完后相当于一个文件
*/
func main() {
	str := "hello wrold"
	fmt.Printf("str = %s\n", str)

	test.Reverse(&str)

	fmt.Printf("str = %s\n", str)

	str = "我们是中国人"
	fmt.Printf("str = %s\n", str)

	test.Reverse(&str)

	fmt.Printf("str = %s\n", str)

	test.Log("str use test.log package method Log")

	test.Logf("str use test.log package method %s", "Logf")
}
