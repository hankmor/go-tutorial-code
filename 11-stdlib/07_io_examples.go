package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// io 包示例

func main() {
	fmt.Println("=== 示例 1: 复制 ===")

	src := strings.NewReader("Hello, World!")
	dst := os.Stdout
	fmt.Print("复制输出: ")
	io.Copy(dst, src)
	fmt.Println()

	fmt.Println("\n=== 示例 2: 读取全部 ===")

	reader := strings.NewReader("Hello")
	data, _ := io.ReadAll(reader)
	fmt.Println("读取内容:", string(data))

	fmt.Println("\n=== 示例 3: 限制读取 ===")

	reader = strings.NewReader("Hello, World!")
	limited := io.LimitReader(reader, 5)
	data, _ = io.ReadAll(limited)
	fmt.Println("限制读取:", string(data))

	fmt.Println("\n=== 示例 4: 多个 Reader ===")

	r1 := strings.NewReader("Hello ")
	r2 := strings.NewReader("World!")
	multiReader := io.MultiReader(r1, r2)
	data, _ = io.ReadAll(multiReader)
	fmt.Println("多个 Reader:", string(data))

	fmt.Println("\n=== 示例 5: 写入多个 Writer ===")

	var buf1, buf2 strings.Builder
	multiWriter := io.MultiWriter(&buf1, &buf2)
	multiWriter.Write([]byte("Hello"))

	fmt.Println("Writer 1:", buf1.String())
	fmt.Println("Writer 2:", buf2.String())
}
