package main

import (
	"fmt"
	"strconv"
)

// strconv 包示例

func main() {
	fmt.Println("=== 示例 1: 字符串转数字 ===")

	// 字符串转整数
	i, err := strconv.Atoi("42")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Atoi:", i)

	// 字符串转 int64
	i64, err := strconv.ParseInt("42", 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("ParseInt:", i64)

	// 字符串转浮点数
	f, err := strconv.ParseFloat("3.14", 64)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("ParseFloat:", f)

	// 字符串转布尔值
	b, err := strconv.ParseBool("true")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("ParseBool:", b)

	fmt.Println("\n=== 示例 2: 数字转字符串 ===")

	// 整数转字符串
	s := strconv.Itoa(42)
	fmt.Println("Itoa:", s)

	// int64 转字符串
	s = strconv.FormatInt(42, 10)
	fmt.Println("FormatInt:", s)

	// 浮点数转字符串
	s = strconv.FormatFloat(3.14, 'f', 2, 64)
	fmt.Println("FormatFloat:", s)

	// 布尔值转字符串
	s = strconv.FormatBool(true)
	fmt.Println("FormatBool:", s)

	fmt.Println("\n=== 示例 3: 错误处理 ===")

	// 转换失败的情况
	_, err = strconv.Atoi("abc")
	if err != nil {
		fmt.Println("转换失败:", err)
	}
}
