package main

import "golang.org/x/sys/windows/registry"
import "fmt"

func main() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\BPA\GRAPHTOOL`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Printf("1")
	}
	defer k.Close()

	s, _, err := k.GetStringValue("NAME")
	if err != nil {
		fmt.Printf("2")
	}
	fmt.Printf("Windows system root is %q\n", s)
}
