package main

import (
	"fmt"
	"net/http"
)

func main() {
	url := "http://www.baidu.com"
	fmt.Printf("feting url = %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	// io.Copy(os.Stdout, res.Body)

	// fmt.Println(res.Header)

	for k, v := range res.Header {
		fmt.Println(k, v)
	}
}
