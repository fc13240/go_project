package main

import "github.com/natefinch/npipe"

func main() {
	conn, err := npipe.Dial(`\\.\pipe\BPA-GraphTool`)
	if err != nil {
	    // handle error
	}
	if _, err := fmt.Fprintln(conn, "Hi server!"); err != nil {
	    // handle error
	}
	r := bufio.NewReader(conn)
	msg, err := r.ReadString('\n')
	if err != nil {
	    // handle eror
	}
	fmt.Println(msg)
}