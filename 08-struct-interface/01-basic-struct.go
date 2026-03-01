package main

import "fmt"

// 定义结构体
type Person struct {
	Name string
	Age  int
}

func main() {
	// 初始化方式 1：按字段顺序
	p1 := Person{"Alice", 30}
	fmt.Printf("p1: %+v\n", p1)

	// 初始化方式 2：指定字段名（推荐）
	p2 := Person{
		Name: "Bob",
		Age:  25,
	}
	fmt.Printf("p2: %+v\n", p2)

	// 初始化方式 3：部分字段（其余为零值）
	p3 := Person{Name: "Charlie"} // Age 为 0
	fmt.Printf("p3: %+v\n", p3)

	// 初始化方式 4：使用 new（返回指针）
	p4 := new(Person)
	p4.Name = "David"
	fmt.Printf("p4: %+v\n", p4)

	// 访问和修改字段
	p := Person{Name: "Alice", Age: 30}
	fmt.Println("Name:", p.Name)

	p.Age = 31
	fmt.Println("Updated Age:", p.Age)

	// 指针访问（Go 会自动解引用）
	ptr := &p
	ptr.Name = "Alicia"
	fmt.Println("Updated Name:", ptr.Name)
}
