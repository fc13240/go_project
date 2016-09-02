package test

import (
	"fmt"
	"sort"
)

// Reverse a string
func Reverse(str *string) {
	*str = reverseString(*str)
}

// 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	fmt.Println(runes)
	return string(runes)
}

func reverseString1(s string) string {
	runes := []rune(s)
	sort.Reverse(runes)
	return string(runes)
}
