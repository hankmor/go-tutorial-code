package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	age := 18
	user := User{Name: "老墨", Age: age}
	// user.Age = validAge(age) ? age: -1 // 编译错误，没有三元运算符
	// Go 没有三元运算符，而是使用 if else
	if validAge(age) {
		user.Age = age
	} else {
		user.Age = -1
	}

	// 非泛型方法
	// ret := choose(validAge(user.Age), age, -1)
	// user.Age = ret.(int) // 返回ret是 any 类型，需要类型断言
	// 泛型方法
	ret := choose1(validAge(user.Age), age, -1)
	user.Age = ret // 泛型方法无需断言
	fmt.Printf("%v\n", user)

	// 运行:
	// go run 02-style/choose.go
}

func choose(b bool, v1 any, v2 any) any {
	if b {
		return v1
	}
	return v2
}

func validAge(age int) bool {
	return age >= 0 && age <= 130
}

// 泛型方法
func choose1[T any](b bool, v1 T, v2 T) T {
	if b {
		return v1
	}
	return v2
}
