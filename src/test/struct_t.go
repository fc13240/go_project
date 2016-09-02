package main

import (
	"fmt"
	"runtime"
	"math"
)

type Vertex struct {
	X, Y,z int
}

var (
	p = Vertex{1, 2, 3}  // has type Vertex
	q = &Vertex{1, 2, 3} // has type *Vertex
	r = Vertex{X: 1}  // Y:0 is implicit
	s = Vertex{}      // X:0 and Y:0
)

type Vertex_Float struct {
	X, Y float64
}
func (v *Vertex_Float) abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
	var f_num = 1.0
	fmt.Printf("%f.", f_num)
	p_new := &p
	p_new.X = 9
	fmt.Println(p, q, r, s, p_new)

	// map
	m := make(map[string]int)

	m["Answer"] = 42
	fmt.Println("The value:", m["Answer"])

	m["Answer"] = 48
	fmt.Println("The value:", m["Answer"])

	delete(m, "Answer")
	fmt.Println("The value:", m["Answer"])

	v, ok := m["Answer"]
	fmt.Println("The value:", v, "Present?", ok)

	m["test"] = 48
	if _, ok := m["test"]; ok {
		fmt.Println("hello")
	} else {
		fmt.Println("world")
	}

	// [] T 是一个元素类型为T的slice
	p := []int{2, 3, 5, 7, 11, 13}
	fmt.Println("p ==", p)

	p = p[1:cap(p)]
	for i := 0; i < len(p); i++ {
		fmt.Printf("p[%d] == %d\n",
			i, p[i])
	}

	for i,j,k := 0, 10, 20; i<j; i++ {
		fmt.Println(i, j, k)
	}

	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	case "windows":
		fmt.Println("windows1.")
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.", os)
	}

	v_Vertex_Float := &Vertex_Float{10, 4}
	fmt.Printf("vertex_float.Abs : %f", v_Vertex_Float.abs())
}
