package main

import "fmt"

func main() {
	const (
		a = iota * 2 // 0
		b            // 2
		c            // 4
		d = iota     // 3（为什么？因为iota始终是从0开始增长的，d之前 iota 的值是2，到d这里就是3）
		e            // 4
	)
	fmt.Println(a, b, c, d, e)
}
