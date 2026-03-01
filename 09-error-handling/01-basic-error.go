package main

import (
	"errors"
	"fmt"
)

// 创建错误的三种方式

// 方式 1：使用 errors.New
func example1() {
	err := errors.New("something went wrong")
	fmt.Println("Error 1:", err)
}

// 方式 2：使用 fmt.Errorf（支持格式化）
func example2() {
	filename := "config.json"
	err := fmt.Errorf("failed to open file: %s", filename)
	fmt.Println("Error 2:", err)
}

// 方式 3：自定义错误类型
type MyError struct {
	Code int
	Msg  string
}

func (e MyError) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Msg)
}

func example3() {
	err := MyError{Code: 404, Msg: "not found"}
	fmt.Println("Error 3:", err)
}

// 基本错误处理模式
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil
}

func main() {
	fmt.Println("=== 创建错误的三种方式 ===")
	example1()
	example2()
	example3()

	fmt.Println("\n=== 基本错误处理 ===")
	// 正常情况
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", result)

	// 错误情况
	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", result) // 不会执行
}
