package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello go!")
	for i, arg := range os.Args {
		fmt.Println(i, "=", arg)
	}
	// 运行:
	// go run 02-style/args.go a b
}
