package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float32 = 10

	fmt.Println(reflect.TypeOf(x))

	v := reflect.ValueOf(x)
	fmt.Printf("kind = %s, val = %f\n", v.Kind(), v.Float())

	//系统默认得到float64的类型，如果右侧值为不带小数的数字要强制设置float64类型
	var x1 float64 = 3.4
	p := reflect.ValueOf(&x1) // 注意：得到X的地址
	fmt.Println("type of p:", p.Type())
	fmt.Println("settability of p:", p.CanSet())
	v1 := p.Elem()
	fmt.Println("settability of v:", v1.CanSet())
	v1.SetFloat(7.1)
	fmt.Println(v1.Interface())
	fmt.Println(x1)

	type T struct {
		A int
		B string
	}
	t := T{203, "mh203"}
	s := reflect.ValueOf(&t).Elem()

	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}
