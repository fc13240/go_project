package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

func main() {
	go func() {
		ch := time.Tick(time.Second)

		for true {
			<-ch
            // 时间layout格式说明
            // http://www.open-open.com/lib/view/open1372919144628.html
			fmt.Printf("%s\n", time.Now().Format("2006-08-09 15:04:05"))
		}
	}()

	http.HandleFunc("/hello", func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("method = %s, referer = %s, path = %s, query = %s\n", req.Method, req.Referer(), req.URL.Path, req.URL.Query())
		fmt.Fprintf(res, "hello, %q", html.EscapeString(req.URL.Path))
	})

	http.HandleFunc("/*", func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("method = %s, referer = %s, path = %s, query = %s\n", req.Method, req.Referer(), req.URL.Path, req.URL.Query())
		fmt.Fprintf(res, "any, %q", html.EscapeString(req.URL.Path))
	})

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		now := time.Now()
		fmt.Printf("%s, method = %s, referer = %s, path = %s, query = %s\n", now.Format("2006-08-09 10:20:30"), req.Method, req.Referer(), req.URL.Path, req.URL.Query())
		fmt.Fprintf(res, "root, %q", html.EscapeString(req.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
