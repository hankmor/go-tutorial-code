package main

import (
	"fmt"
	"sync"
	"time"
)

// 扇入扇出模式示例

// 生成器：生成数据
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// 处理器：处理数据
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			time.Sleep(100 * time.Millisecond) // 模拟耗时操作
			out <- n * n
		}
	}()
	return out
}

// Fan-out: 启动多个 worker 处理同一个输入
func fanOut(input <-chan int, workers int) []<-chan int {
	outputs := make([]<-chan int, workers)
	for i := 0; i < workers; i++ {
		outputs[i] = square(input)
	}
	return outputs
}

// Fan-in: 合并多个 channel 到一个
func fanIn(channels ...<-chan int) <-chan int {
	output := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for val := range c {
				output <- val
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func main() {
	fmt.Println("=== 示例 1: 简单管道 ===")
	input := generator(1, 2, 3, 4, 5)
	output := square(input)

	for result := range output {
		fmt.Println("Result:", result)
	}

	fmt.Println("\n=== 示例 2: Fan-out/Fan-in ===")
	input2 := generator(1, 2, 3, 4, 5, 6, 7, 8)

	// Fan-out: 3 个 worker 并发处理
	workers := fanOut(input2, 3)

	// Fan-in: 合并结果
	results := fanIn(workers...)

	for result := range results {
		fmt.Println("Result:", result)
	}
}
