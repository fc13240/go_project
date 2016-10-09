package test

import "fmt"

// Log any message
func Log(msg ...interface{}) {
	fmt.Println(msg)
}

// Logf any message use fmt.Printf
func Logf(format string, msg ...interface{}) {
	fmt.Printf(format, msg)
}
