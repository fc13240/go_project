package main

import "fmt"

// Person struct
type Person struct {
	Name string
}

// SayHello for Person
func (p *Person) SayHello() {
	fmt.Printf("hello world, my name is %s\n", p.Name)
}

// Man struct
type Man struct {
	Person
	Name string
}

func main() {
	p := Person{"tonny"}

	p.SayHello()

	// man := new(Man)
	man := Man{p, "man"}
	man.Name = "main1"
	man.SayHello()
}
