package main

import "fmt"

// 坑 1：值接收者 vs 指针接收者

type Counter struct {
	count int
}

// 值接收者：无法修改原对象
func (c Counter) Increment() {
	c.count++ // 只修改了副本
	fmt.Println("Inside Increment:", c.count)
}

// 指针接收者：可以修改原对象
func (c *Counter) IncrementPtr() {
	c.count++ // 修改了原对象
	fmt.Println("Inside IncrementPtr:", c.count)
}

func main() {
	fmt.Println("=== 值接收者 vs 指针接收者 ===")

	c := Counter{count: 0}
	fmt.Println("Initial count:", c.count)

	c.Increment()
	fmt.Println("After Increment:", c.count) // 0（没变！）

	c.IncrementPtr()
	fmt.Println("After IncrementPtr:", c.count) // 1（变了）

	// 教训：
	// - 需要修改对象状态 → 用指针接收者
	// - 对象很大（避免拷贝）→ 用指针接收者
	// - 只读操作且对象小 → 用值接收者
	// - 同一类型的方法，接收者类型要保持一致
}
