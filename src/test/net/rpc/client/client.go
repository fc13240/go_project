package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")

	if err != nil {
		log.Fatal("dialing error: ", err)
	}

	var args = "hello from client" + os.Args[1]

	var reply string
	err = client.Call("Echo.Hi", args, &reply)
	if err != nil {
		log.Fatal("call error: ", err)
	}

	fmt.Printf("Arith: args = %s, reply = %s", args, reply)
}
