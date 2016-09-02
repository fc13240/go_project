package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"strings"
	"github.com/natefinch/npipe"
)


func main() {
	// argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	var duration_Seconds time.Duration = (2) * time.Millisecond
	conn, err := npipe.DialTimeout(`\\.\pipe\BPA-GRAPHTOOL-COMMAND`, duration_Seconds)
	defer conn.Close()

	if err != nil {
	    // handle error
	    fmt.Println("{\"err\": \"请联系管理员启动服务\"}");
	    os.Exit(0)
	}
	// fmt.Println("1");
	// put msg to socket
	if _, err := fmt.Fprintln(conn, strings.Join(argsWithoutProg, " ")+"\n"); err != nil {
	    // handle error
	    fmt.Println("{\"err\": \"没有得到服务响应\"}");
	    os.Exit(0)
	}
	// fmt.Println("11");
	r := bufio.NewReader(conn)
	msg, err := r.ReadString('\n')
	fmt.Println(msg)
	if err != nil {
	    // handle eror
	    fmt.Println("{\"err\": \"响应没有EOF\"}");
	    os.Exit(0)
	}
}
