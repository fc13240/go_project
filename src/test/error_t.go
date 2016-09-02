package main

import (
	"fmt"
	"time"
)

type MyError struct {
	When time.Time
	What string
}
func (err *MyError) Error() string{
	return fmt.Sprintf("at %v, %s", err.When, err.What)
}
func run() error{
	return &MyError { 
		time.Now(), 
		"it not work!" }
}
func main() {
	if err := run(); err != nil {
		fmt.Println(err);
	}
}