package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(rw, e.Error(), http.StatusInternalServerError)
				log.Printf("WARN: panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()

		fn(rw, req)
	}
}

func test1() func(rw http.ResponseWriter, req *http.Request) {
	var isFetching = false
	now := time.Now()
	fmt.Println(now)

	var rwArr []http.ResponseWriter
	return func(rw http.ResponseWriter, req *http.Request) {
		rwArr = append(rwArr, rw)
		if !isFetching {
			isFetching = true
			time.AfterFunc(time.Second*5, func() {
				var result = "test" + strconv.Itoa(now.Nanosecond()) + " | " + strconv.Itoa(time.Now().Nanosecond())

				for _, v := range rwArr {
					v.Write([]byte(result))
				}

				fmt.Println("after fetching")

				isFetching = false
			})
		}

		fmt.Printf("isFetching = %v\n", isFetching)
	}
}

// golang里没有类似js里的setTimeout响应，且是同步操作，要实现类似延时响应的效果，需要借助chan 让程序挂起
func test() func(rw http.ResponseWriter, req *http.Request) {
	var isFetching = false
	now := time.Now()
	fmt.Println(now)

	ch := make(chan string)
	var lenWaiting int
	return func(rw http.ResponseWriter, req *http.Request) {
		url := req.RequestURI

		lenWaiting++

		fmt.Printf("%s isFetching = %v, len = %d\n\n", url, isFetching, len(ch))
		if !isFetching {
			isFetching = true
			// 之前所有的请求都同时返回结果
			time.AfterFunc(time.Second*5, func() {
				var result = "test" + strconv.Itoa(now.Nanosecond()) + " | " + strconv.Itoa(time.Now().Nanosecond())

				var i = 0
				for ; i < lenWaiting; i++ {
					ch <- result + " | " + strconv.Itoa(i)
				}
				fmt.Printf("after fetching %d\n\n", lenWaiting)

				isFetching = false
			})
		}

		fmt.Printf("isFetching = %v, len = %d\n\n", isFetching, lenWaiting)

		rs := <-ch
		lenWaiting--
		fmt.Printf("response %d\n\n", lenWaiting)

		rw.Write([]byte(rs + " | " + req.RequestURI))

	}
}

func main() {
	http.HandleFunc("/test", safeHandler(test()))

	http.ListenAndServe(":8080", nil)
}
