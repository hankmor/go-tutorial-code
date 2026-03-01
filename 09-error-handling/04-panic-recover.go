package main

import "fmt"

// Panic 和 Recover 示例

// 示例 1：基本的 panic
func example1() {
	fmt.Println("Start")
	panic("something went wrong")
	fmt.Println("End") // 不会执行
}

// 示例 2：使用 recover 捕获 panic
func safeCall(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	fn() // 执行可能 panic 的函数
}

// 示例 3：panic 的传播
func level3() {
	fmt.Println("Level 3: about to panic")
	panic("error at level 3")
}

func level2() {
	fmt.Println("Level 2: calling level 3")
	level3()
	fmt.Println("Level 2: after level 3") // 不会执行
}

func level1() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Caught at level 1:", r)
		}
	}()

	fmt.Println("Level 1: calling level 2")
	level2()
	fmt.Println("Level 1: after level 2") // 不会执行
}

// 示例 4：什么时候使用 panic
func mustOpen(filename string) {
	// 这是一个不好的例子，实际应该返回 error
	if filename == "" {
		panic("filename cannot be empty")
	}
	fmt.Printf("Opening file: %s\n", filename)
}

func main() {
	// 示例 1：基本 panic（会导致程序崩溃，所以注释掉）
	// example1()

	// 示例 2：使用 recover 捕获 panic
	fmt.Println("=== Example 2: Recover ===")
	safeCall(func() {
		panic("something went wrong")
	})
	fmt.Println("Program continues after recover\n")

	// 示例 3：panic 的传播
	fmt.Println("=== Example 3: Panic Propagation ===")
	level1()
	fmt.Println("Program continues after catching panic\n")

	// 示例 4：什么时候使用 panic
	fmt.Println("=== Example 4: When to Use Panic ===")
	safeCall(func() {
		mustOpen("") // 会 panic
	})
	fmt.Println("Program continues")
}
