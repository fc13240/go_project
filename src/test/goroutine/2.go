package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const (
	GOROUTINENUM = 10
	TASKNUM      = 10
)

func run(urls []string, callback func()) {
	chReq := make(chan string, GOROUTINENUM)
	chRes := make(chan string, GOROUTINENUM)

	for i := 0; i < GOROUTINENUM; i++ {
		// 相当于开启了一个服务
		go func(index int) {
			fmt.Printf("1goroutine index = %d\n", index)
			for {
				url := <-chReq

				fmt.Println(len(url) > 0 && url != "quit")
				if len(url) > 0 && url != "quit" {
					fmt.Println(reflect.TypeOf(url))
					delay := time.Second * time.Duration(rand.Intn(5)+1)
					fmt.Printf("delay = %d, url = %s\n", delay, url)
					time.Sleep(delay)
					fmt.Printf("after delay url = %s\n", url)
					chRes <- url
				} else {
					fmt.Printf("exit goroutine index = %d\n", index)
					return
				}
			}
		}(i)
	}

	lenUrls := len(urls)
	go func() {
		for _, url := range urls {
			chReq <- url
		}
	}()

	for i := 0; i < lenUrls; i++ {
		url := <-chRes

		fmt.Printf("res len = %d, url = %s\n", i, url)
	}

	fmt.Println("===========work end===========")

	if nil != callback {
		callback()
	}
}
func main() {
	var urls = []string{}
	for i := 0; i < TASKNUM; i++ {
		url := fmt.Sprintf("http://%d.test.com", i)
		urls = append(urls, url)
	}

	run(urls, nil)

	fmt.Println("after run outer")
	// chReq := make(chan string, GOROUTINENUM)
	// chRes := make(chan string, GOROUTINENUM)

	// for i := 0; i < GOROUTINENUM; i++ {
	// 	// 相当于开启了一个服务
	// 	go func(index int) {
	// 		fmt.Printf("1goroutine index = %d\n", index)
	// 		for {
	// 			url := <-chReq

	// 			fmt.Println(len(url) > 0 && url != "quit")
	// 			if len(url) > 0 && url != "quit" {
	// 				fmt.Println(reflect.TypeOf(url))
	// 				delay := time.Second * time.Duration(rand.Intn(5)+1)
	// 				fmt.Printf("delay = %d, url = %s\n", delay, url)
	// 				time.Sleep(delay)
	// 				fmt.Printf("after delay url = %s\n", url)
	// 				chRes <- url
	// 			} else {
	// 				fmt.Printf("exit goroutine index = %d\n", index)
	// 				return
	// 			}
	// 		}
	// 	}(i)
	// }
	// // var urls = []string{}
	// go func() {
	// 	for i := 0; i < TASKNUM; i++ {
	// 		url := fmt.Sprintf("http://%d.test.com", i)
	// 		chReq <- url
	// 	}
	// }()

	// for i := 0; i < TASKNUM; i++ {
	// 	url := <-chRes

	// 	fmt.Printf("res len = %d, url = %s\n", i, url)
	// }

	// fmt.Println("===========work end===========")

	// for i := 0; i < GOROUTINENUM; i++ {
	// 	chReq <- "quit"
	// }
}
