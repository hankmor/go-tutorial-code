package utils

// Add 两个数相加
func Add(a, b int) int {
	return a + b
}

// Subtract 两个数相减
func Subtract(a, b int) int {
	return a - b
}

// Multiply 两个数相乘
func Multiply(a, b int) int {
	return a * b
}

// Divide 两个数相除
func Divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}
