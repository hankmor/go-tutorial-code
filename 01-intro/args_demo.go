package main

import (
	"fmt"
	"os"
)

func main() {
	// 获取命令行参数
	args := os.Args

	// 输出参数个数和内容
	fmt.Printf("参数个数: %d\n", len(args)-1)
	fmt.Printf("程序名称: %s\n", args[0])
	fmt.Println("参数列表:")
	for i, arg := range args[1:] {
		fmt.Printf("  参数 %d: %s\n", i+1, arg)
	}
	// 运行:
	// go run 01-intro/args_demo.go
	// go run 01-intro/args_demo.go -c 123 -d abc
	// 注意对比输出的区别
}
