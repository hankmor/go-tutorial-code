package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

// DemoTeeReader 演示 TeeReader: 边读边写
func DemoTeeReader() {
	fmt.Println("--- 演示 TeeReader: 边读边写 ---")
	// 模拟一个数据源 (Reader)
	src := strings.NewReader("老墨极客时间，关注不迷路！")

	// 模拟一个目的地 (Writer)，比如一个文件或缓冲区
	var buf bytes.Buffer

	// 创建 TeeReader: 从 tee 读数据时，会自动同步写到 buf
	tee := io.TeeReader(src, &buf)

	// 开始读取
	p := make([]byte, 100)
	n, _ := tee.Read(p)

	fmt.Printf("读取到的内容: %s\n", p[:n])
	fmt.Printf("自动同步到缓冲区的内容: %s\n", buf.String())
}

// DemoMultiReader 演示 MultiReader: 串联读取
func DemoMultiReader() {
	fmt.Println("\n--- 演示 MultiReader: 串联读取 ---")
	r1 := strings.NewReader("第一部分；")
	r2 := strings.NewReader("第二部分；")
	r3 := strings.NewReader("结语。")

	// 串联三个 Reader
	combined := io.MultiReader(r1, r2, r3)

	// 一次性读出来
	data, _ := io.ReadAll(combined)
	fmt.Printf("合并后的内容: %s\n", data)
}

// DemoLimitReader 演示 LimitReader: 限制读取长度
func DemoLimitReader() {
	fmt.Println("\n--- 演示 LimitReader: 限制读取长度 ---")
	src := strings.NewReader("这是一段很长的机密数据，我只想让你看前五个字")

	// 限制只能读取 15 个字节 (假设一个汉字 3 字节)
	limited := io.LimitReader(src, 15)

	data, _ := io.ReadAll(limited)
	fmt.Printf("受限读取内容: %s\n", data)
}

// DemoMemoryIO 演示内存中的 IO (bytes.Buffer)
func DemoMemoryIO() {
	fmt.Println("\n--- 演示内存中的 IO (bytes.Buffer) ---")
	var buf bytes.Buffer

	// 既可以像 Writer 一样写
	buf.Write([]byte("Hello "))
	fmt.Fprint(&buf, "Geek Mo!")

	// 也可以像 Reader 一样读
	result := buf.String()
	fmt.Printf("缓冲区最终结果: %s\n", result)
}

// DemoNetworkToLocal 演示网络流直接转存本地 (io.Copy)
func DemoNetworkToLocal() {
	fmt.Println("\n--- 演示网络流直接转存本地 (io.Copy) ---")

	// 模拟一个 HTTP 服务器返回数据
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "这是一段来自网络的数据流")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// 请求数据
	resp, err := http.Get(server.URL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// 创建本地文件
	file, _ := os.Create("network_data.txt")
	defer file.Close()

	// 核心：直接用 io.Copy 连接网络流和文件流
	_, err = io.Copy(file, resp.Body)
	if err == nil {
		fmt.Println("成功使用 io.Copy 将网络数据保存到 network_data.txt")
	}
}

func main() {
	DemoTeeReader()
	DemoMultiReader()
	DemoLimitReader()
	DemoMemoryIO()
	DemoNetworkToLocal()
}
