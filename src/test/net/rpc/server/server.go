package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

/*Echo ...*/
type Echo int

/*Hi something*/
func (t *Echo) Hi(args string, reply *string) error {
	*reply = "echo :" + args

	return nil
}
func main() {
	rpc.Register(new(Echo))
	rpc.HandleHTTP()

	// 必须保证addr有指定的IP及相关端口，否则client连接失败
	l, e := net.Listen("tcp", "127.0.0.1:1234")

	if e != nil {
		log.Fatal("Listen error: ", e)
	}

	fmt.Print(l)
	http.Serve(l, nil)
}
