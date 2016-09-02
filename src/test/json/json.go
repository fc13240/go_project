package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	type Person struct {
		Name   string
		Age    int
		height float32
	}
	p := Person{"tonny", 10, 1.72}
	fmt.Printf("name = %s, age = %d, height = %f\n", p.Name, p.Age, p.height)

	//相对json.Marshal方法，p为包外对象，所以要想访问到值的话必须保证属性名第一个字母大写
	b, err := json.Marshal(p)

	if err != nil {
		log.Fatal("json error: ", err)
	}

	fmt.Printf("jsonstr = %s\n", b)

	var man Person
	json.Unmarshal(b, &man)

	fmt.Printf("name = %s, age = %d, height = %f\n", man.Name, man.Age, man.height)

	type Book struct {
		Title       string
		Authors     []string
		Publisher   string
		IsPublished bool
		Price       float64
	}

	gobook := Book{"Go语言编程", []string{"XuShiwei", "HughLv", "Pandaman", "GuaguaSong", "HanTuo", "BertYuan", "XuDaoli"}, "ituring.com.cn", true, 9.99}

	bGobook, eGobook := json.Marshal(gobook)
	fmt.Printf("%s %s", bGobook, eGobook)
}
