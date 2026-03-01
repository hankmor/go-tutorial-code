package main

import (
	"fmt"
	"time"
)

// 实战建议 4：空结构体的妙用

// 1. 实现 Set（集合）
type Set map[string]struct{}

func NewSet() Set {
	return make(Set)
}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Remove(key string) {
	delete(s, key)
}

func (s Set) Contains(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Size() int {
	return len(s)
}

// 2. 信号通道
func worker(id int, done chan struct{}) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
	done <- struct{}{} // 发送信号
}

// 3. 实现接口但不需要状态
type NoOpWriter struct{}

func (NoOpWriter) Write(p []byte) (n int, err error) {
	return len(p), nil // 什么都不做，直接返回
}

func main() {
	fmt.Println("=== 空结构体的妙用 ===")

	// 1. Set 示例
	fmt.Println("\n--- Set 示例 ---")
	set := NewSet()
	set.Add("apple")
	set.Add("banana")
	set.Add("apple") // 重复添加

	fmt.Println("Contains apple:", set.Contains("apple"))
	fmt.Println("Contains orange:", set.Contains("orange"))
	fmt.Println("Set size:", set.Size())

	set.Remove("apple")
	fmt.Println("After remove, contains apple:", set.Contains("apple"))

	// 2. 信号通道示例
	fmt.Println("\n--- 信号通道示例 ---")
	done := make(chan struct{})

	go worker(1, done)
	go worker(2, done)

	// 等待两个 worker 完成
	<-done
	<-done
	fmt.Println("All workers done")

	// 3. NoOpWriter 示例
	fmt.Println("\n--- NoOpWriter 示例 ---")
	writer := NoOpWriter{}
	n, err := writer.Write([]byte("test data"))
	fmt.Printf("Wrote %d bytes, error: %v\n", n, err)

	// 空结构体的优势：
	// - 不占用内存（sizeof(struct{}) == 0）
	// - 语义清晰（表示"无数据，只有信号"）
	// - 性能好（无内存分配）
}
