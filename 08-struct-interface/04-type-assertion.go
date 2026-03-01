package main

import "fmt"

type Speaker interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Woof!"
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "Meow!"
}

// 类型断言示例
func checkType(s Speaker) {
	// 安全的类型断言
	if dog, ok := s.(Dog); ok {
		fmt.Printf("It's a dog named %s\n", dog.Name)
	} else {
		fmt.Println("Not a dog")
	}
}

// 类型选择示例
func describe(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %s\n", v)
	case Dog:
		fmt.Printf("Dog: %s\n", v.Name)
	case Cat:
		fmt.Printf("Cat: %s\n", v.Name)
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

func main() {
	var s Speaker = Dog{Name: "Buddy"}

	// 类型断言
	checkType(s)

	// 类型选择
	describe(42)
	describe("hello")
	describe(Dog{Name: "Buddy"})
	describe(Cat{Name: "Kitty"})
	describe([]int{1, 2, 3})
}
