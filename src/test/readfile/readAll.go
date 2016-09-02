package main

import (
	"os"
	"io/ioutil"
	"fmt"
)
func readAll(filePath string) ([]byte, error){
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}
func main(){
	content, err := readAll("E:\\source\\nodejs_project\\GraphTool\\shell\\build\\target\\GraphTool_v0.5.3_32\\main.gt")
	fmt.Println(string(content), err)
}