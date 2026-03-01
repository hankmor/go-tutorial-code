package main

import "fmt"

// 坑 2：接口的 nil 判断

type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	if e == nil {
		return "nil error"
	}
	return e.msg
}

// 错误的做法
func doSomethingWrong() error {
	var err *MyError // nil 指针
	return err       // 返回的是接口类型
}

// 正确的做法
func doSomethingRight() error {
	var err *MyError
	if err != nil {
		return err
	}
	return nil // 返回 nil 接口，而不是 nil 指针
}

func main() {
	fmt.Println("=== 接口的 nil 判断陷阱 ===")

	// 错误示例
	err1 := doSomethingWrong()
	fmt.Printf("err1 == nil: %v\n", err1 == nil) // false！
	fmt.Printf("err1 type: %T, value: %v\n", err1, err1)

	// 正确示例
	err2 := doSomethingRight()
	fmt.Printf("err2 == nil: %v\n", err2 == nil) // true

	// 原因：接口包含类型和值两部分
	// 只有两者都为 nil，接口才为 nil
	// err1 的类型是 *MyError，值是 nil，所以接口不为 nil

	// 安全的检查方式
	if err1 != nil {
		fmt.Println("Error occurred (but actually it's nil pointer)")
		// 如果调用 err1.Error() 可能会 panic
		// fmt.Println(err1.Error())
	}
}
