package main

import (
	"fmt"
	"sync"
	"time"
)

// 常见并发陷阱示例

// 陷阱 1: 忘记等待 Goroutine
func pitfall1Wrong() {
	fmt.Println("=== 陷阱 1: 忘记等待 (错误示例) ===")
	for i := 0; i < 5; i++ {
		go func(n int) {
			fmt.Println(n)
		}(i)
	}
	// 程序立即退出，可能什么都没打印
	time.Sleep(100 * time.Millisecond)
}

func pitfall1Right() {
	fmt.Println("\n=== 陷阱 1: 正确做法 ===")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}
	wg.Wait()
}

// 陷阱 2: 闭包捕获循环变量
func pitfall2Wrong() {
	fmt.Println("\n=== 陷阱 2: 闭包陷阱 (错误示例) ===")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i) // 可能都打印 5
		}()
	}
	wg.Wait()
}

func pitfall2Right() {
	fmt.Println("\n=== 陷阱 2: 正确做法 ===")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}
	wg.Wait()
}

// 陷阱 3: 向已关闭的 Channel 发送
func pitfall3() {
	fmt.Println("\n=== 陷阱 3: 向已关闭的 Channel 发送 ===")
	ch := make(chan int, 2)
	ch <- 1
	close(ch)

	// 这会 panic
	// ch <- 2

	// 但可以继续接收
	fmt.Println("Received:", <-ch)
	fmt.Println("Received:", <-ch) // 零值
}

// 陷阱 4: 无缓冲 Channel 死锁
func pitfall4Wrong() {
	fmt.Println("\n=== 陷阱 4: 无缓冲 Channel 死锁 (注释掉) ===")
	// ch := make(chan int)
	// ch <- 42 // 死锁！
	// fmt.Println(<-ch)
	fmt.Println("(示例已注释，避免死锁)")
}

func pitfall4Right() {
	fmt.Println("\n=== 陷阱 4: 正确做法 ===")
	ch := make(chan int, 1) // 使用缓冲
	ch <- 42
	fmt.Println(<-ch)
}

func main() {
	pitfall1Wrong()
	pitfall1Right()
	pitfall2Wrong()
	pitfall2Right()
	pitfall3()
	pitfall4Wrong()
	pitfall4Right()
}
