package main

import (
	"os"
	"fmt"
	"bufio"
	"strconv"
)

var filePath = "E:\\source\\nodejs_project\\GraphTool\\shell\\build\\target\\GraphTool_v0.5.3_32\\main.gt"
func main() {
	file, err := os.Open(filePath);
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}

	inputFile, inputError := os.Open(filePath)
	if inputError != nil {
		fmt.Println(inputError)
	}

	inputReader := bufio.NewReader(inputFile)

	lineCounter := 0

	for {
		buf := make([]byte, 1024)

		n, _ := inputReader.Read(buf)

		lineCounter++

		if n == 0 {
			break
		}

		fmt.Printf("%d==>\n%s\n", n, string(buf[:n]))
	}

	fmt.Println("\n\nlineCounter = "+strconv.Itoa(lineCounter));
}