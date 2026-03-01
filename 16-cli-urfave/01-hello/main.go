package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// 最简单的 CLI 示例

func main() {
	app := cli.NewApp()
	app.Name = "my-cli"
	app.Usage = "我的第一个 CLI 工具"
	app.Version = "0.1.0"

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
