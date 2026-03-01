package main

import (
	"fmt"
	"time"
)

// 限流器示例

func main() {
	fmt.Println("=== 示例 1: 简单限流器 ===")
	// 每 200ms 处理一个请求（每秒 5 个）
	limiter := time.Tick(200 * time.Millisecond)

	requests := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		requests <- i
	}
	close(requests)

	for req := range requests {
		<-limiter // 等待限流器
		fmt.Printf("Processing request %d at %v\n", req, time.Now().Format("15:04:05.000"))
	}

	fmt.Println("\n=== 示例 2: 突发限流器 ===")
	// 允许突发 3 个请求，然后每 200ms 一个
	burstyLimiter := make(chan time.Time, 3)

	// 预填充 3 个令牌
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	// 每 200ms 添加一个令牌
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	burstyRequests := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)

	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Printf("Processing bursty request %d at %v\n", req, time.Now().Format("15:04:05.000"))
	}
}
