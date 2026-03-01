package main

import "fmt"

// fmt 包示例

func main() {
	fmt.Println("=== 示例 1: 基本输出 ===")
	fmt.Print("Hello")
	fmt.Print("World\n")

	fmt.Println("Hello")
	fmt.Println("World")

	name := "Go"
	age := 15
	fmt.Printf("Name: %s, Age: %d\n", name, age)

	fmt.Println("\n=== 示例 2: 常用占位符 ===")
	// 字符串
	fmt.Printf("%s\n", "hello")
	fmt.Printf("%q\n", "hello")

	// 整数
	fmt.Printf("%d\n", 42)
	fmt.Printf("%b\n", 42)
	fmt.Printf("%x\n", 42)

	// 浮点数
	fmt.Printf("%f\n", 3.14)
	fmt.Printf("%.2f\n", 3.14159)

	// 布尔值
	fmt.Printf("%t\n", true)

	// 通用
	fmt.Printf("%v\n", 42)
	fmt.Printf("%+v\n", struct{ X int }{42})
	fmt.Printf("%#v\n", []int{1, 2})
	fmt.Printf("%T\n", 42)

	fmt.Println("\n=== 示例 3: 格式化字符串 ===")
	msg := fmt.Sprintf("Hello, %s!", "World")
	fmt.Println(msg)

	err := fmt.Errorf("failed to open file: %s", "config.json")
	fmt.Println(err)
}
