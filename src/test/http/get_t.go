package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
)

func main() {
	resp, _ := http.Get("http://www.baidu.com")
	fmt.Println(resp)

	defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        // handle error
    }
 
    fmt.Println(string(body))
}