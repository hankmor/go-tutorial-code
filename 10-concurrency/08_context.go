package main

import (
	"context"
	"fmt"
	"time"
)

// Context 控制 Goroutine 生命周期示例

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d stopping: %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("Worker %d working\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	fmt.Println("=== 示例 1: 使用 Context 超时控制 ===")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	<-ctx.Done()
	fmt.Println("All workers stopped")

	time.Sleep(time.Second)

	fmt.Println("\n=== 示例 2: 手动取消 ===")
	ctx2, cancel2 := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx2.Done():
				fmt.Println("Task cancelled")
				return
			default:
				fmt.Println("Task running...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	cancel2() // 手动取消
	time.Sleep(time.Second)
}
