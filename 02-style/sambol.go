package main

import "fmt"

// import _ "github.com/go-sql-driver/mysql"

func main() {
	_a := 1
	fmt.Println(_a)
	b := 2
	fmt.Println(b)
	// 3c := 3 // 不能编译
	你好 := 5
	fmt.Println(你好)

	// fmt := 4
	// fmt.Println(fmt)
	// fmt.println("a") // 不能编译

	name, _ := getName()
	// name, err := getName() // 声明了变量err却不使用，编译错误
	fmt.Println(name)
}

// sayHello is a function to say hello with the giving name.
func sayHello(name string) {
	fmt.Println("hello, ", name)
}

func getName() (string, error) {
	return "huzhou", nil
}
