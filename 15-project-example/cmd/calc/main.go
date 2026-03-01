// main.go
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/hankmor/calc/pkg/calculator"
)

func main() {
	// 注意: * 号在终端中表示通配符，所以直接传入 a * b 参数导致无法识别，需要使用 a '*' b, 或者转义 a \* b
	if len(os.Args) < 4 {
		fmt.Println("Usage: calc <num1> <operator> <num2>")
		fmt.Println("Operators: +, -, *, /")
		fmt.Println("Note: Use quotes for * operator: calc 5 '*' 3")
		os.Exit(1)
	}

	// 解析参数
	a, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		color.Red("Error: invalid first number")
		os.Exit(1)
	}

	operator := os.Args[2]

	b, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		color.Red("Error: invalid second number")
		os.Exit(1)
	}

	// 创建计算器
	calc := calculator.New()

	// 执行运算
	var result float64
	switch operator {
	case "+":
		result = calc.Add(a, b)
	case "-":
		result = calc.Subtract(a, b)
	case "*":
		result = calc.Multiply(a, b)
	case "/":
		r, err := calc.Divide(a, b)
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}
		result = r
	default:
		color.Red("Error: invalid operator")
		os.Exit(1)
	}

	// 彩色输出结果
	color.Green("Result: %.2f", result)

	// 显示历史记录
	if len(calc.History()) > 0 {
		color.Cyan("\nHistory:")
		for _, record := range calc.History() {
			fmt.Println("  " + record)
		}
	}
}

// func main() {
// 	// 注意: * 号在终端中表示通配符，所以直接传入 a * b 参数导致无法识别，需要使用 a '*' b, 或者转义 a \* b
//     if len(os.Args) < 4 {
//         fmt.Println("Usage: calc <num1> <operator> <num2>")
//         fmt.Println("Operators: +, -, *, /")
//         fmt.Println("Note: Use quotes for * operator: calc 5 '*' 3")
//         os.Exit(1)
//     }

//     // 解析参数
//     a, err := strconv.ParseFloat(os.Args[1], 64)
//     if err != nil {
//         fmt.Println("Error: invalid first number")
//         os.Exit(1)
//     }

//     operator := os.Args[2]

//     // 检测可能的 shell 通配符展开问题
//     if len(os.Args) > 4 {
//         fmt.Println("Error: too many arguments")
//         fmt.Println("Hint: If using *, wrap it in quotes: calc 5 '*' 3")
//         os.Exit(1)
//     }

//     b, err := strconv.ParseFloat(os.Args[3], 64)
//     if err != nil {
//         fmt.Println("Error: invalid second number")
//         os.Exit(1)
//     }

//     // 创建计算器
//     calc := calculator.New()

//     // 执行运算
//     var result float64
//     switch operator {
//     case "+":
//         result = calc.Add(a, b)
//     case "-":
//         result = calc.Subtract(a, b)
//     case "*":
//         result = calc.Multiply(a, b)
//     case "/":
//         r, err := calc.Divide(a, b)
//         if err != nil {
//             fmt.Printf("Error: %v\n", err)
//             os.Exit(1)
//         }
//         result = r
//     default:
//         fmt.Println("Error: invalid operator")
//         os.Exit(1)
//     }

//     // 输出结果
//     fmt.Printf("Result: %.2f\n", result)
// }
