package main

import (
    "os/exec" // "os/exec" in go1
    "fmt"
)

func main(){
    cmd := exec.Command("notepad", "hello.txt")
    buf, err := cmd.Output()
    fmt.Printf("%s\n%s",buf,err)
}