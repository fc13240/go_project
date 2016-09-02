package main

import (
	"os"
	"fmt"
	"bufio"
	"io"
)
var filePath = "E:\\source\\nodejs_project\\GraphTool\\shell\\build\\target\\GraphTool_v0.5.3_32\\main.gt"
func main() {
	file, err := os.Open(filePath);
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	lineCounter := 0
	inputReader := bufio.NewReader(file)
	for{
		inputString, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {
			return
		}
		lineCounter++

		fmt.Printf("%d %s", lineCounter, inputString)
	}
}