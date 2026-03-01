package main

import "fmt"

// 定义接口
type Speaker interface {
	Speak() string
}

// Dog 实现 Speaker
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Woof!"
}

// Cat 实现 Speaker
type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "Meow!"
}

// 接收任何 Speaker
func introduce(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	d := Dog{Name: "Buddy"}
	c := Cat{Name: "Kitty"}

	// 多态：同一个函数，不同的行为
	introduce(d) // Woof!
	introduce(c) // Meow!

	// 接口变量
	var s Speaker
	s = d
	fmt.Println("Speaker says:", s.Speak())

	s = c
	fmt.Println("Speaker says:", s.Speak())
}
