package main

import (
	"fmt"
	"os"
)

// os 包示例

func main() {
	fmt.Println("=== 示例 1: 文件操作 ===")

	// 写入文件
	err := os.WriteFile("test.txt", []byte("Hello, World!"), 0644)
	if err != nil {
		fmt.Println("写入错误:", err)
		return
	}
	fmt.Println("文件已写入")

	// 读取文件
	data, err := os.ReadFile("test.txt")
	if err != nil {
		fmt.Println("读取错误:", err)
		return
	}
	fmt.Println("文件内容:", string(data))

	// 删除文件
	err = os.Remove("test.txt")
	if err != nil {
		fmt.Println("删除错误:", err)
		return
	}
	fmt.Println("文件已删除")

	fmt.Println("\n=== 示例 2: 目录操作 ===")

	// 创建目录
	err = os.Mkdir("testdir", 0755)
	if err != nil {
		fmt.Println("创建目录错误:", err)
	} else {
		fmt.Println("目录已创建")
	}

	// 创建多级目录
	err = os.MkdirAll("path/to/dir", 0755)
	if err != nil {
		fmt.Println("创建多级目录错误:", err)
	} else {
		fmt.Println("多级目录已创建")
	}

	// 删除目录
	err = os.Remove("testdir")
	if err != nil {
		fmt.Println("删除目录错误:", err)
	} else {
		fmt.Println("目录已删除")
	}

	// 删除目录及其内容
	err = os.RemoveAll("path")
	if err != nil {
		fmt.Println("删除目录树错误:", err)
	} else {
		fmt.Println("目录树已删除")
	}

	fmt.Println("\n=== 示例 3: 环境变量 ===")

	// 获取环境变量
	home := os.Getenv("HOME")
	fmt.Println("HOME:", home)

	// 设置环境变量
	os.Setenv("MY_VAR", "value")
	fmt.Println("MY_VAR:", os.Getenv("MY_VAR"))

	fmt.Println("\n=== 示例 4: 文件信息 ===")

	// 创建测试文件
	os.WriteFile("info.txt", []byte("test"), 0644)

	// 获取文件信息
	info, err := os.Stat("info.txt")
	if err != nil {
		fmt.Println("获取信息错误:", err)
	} else {
		fmt.Println("文件名:", info.Name())
		fmt.Println("大小:", info.Size())
		fmt.Println("是否目录:", info.IsDir())
		fmt.Println("修改时间:", info.ModTime())
	}

	// 清理
	os.Remove("info.txt")
}
