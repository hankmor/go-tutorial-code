package main

import "fmt"

// 定义接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 矩形
type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 圆形
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14159 * c.Radius
}

// 打印形状信息
func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())

	// 类型断言
	switch v := s.(type) {
	case Rectangle:
		fmt.Printf("Rectangle: %.2f x %.2f\n", v.Width, v.Height)
	case Circle:
		fmt.Printf("Circle: radius %.2f\n", v.Radius)
	}
	fmt.Println()
}

func main() {
	r := Rectangle{Width: 10, Height: 5}
	c := Circle{Radius: 7}

	printShapeInfo(r)
	printShapeInfo(c)

	// 使用接口切片
	shapes := []Shape{
		Rectangle{Width: 5, Height: 3},
		Circle{Radius: 4},
		Rectangle{Width: 8, Height: 6},
	}

	totalArea := 0.0
	for _, shape := range shapes {
		totalArea += shape.Area()
	}
	fmt.Printf("Total area: %.2f\n", totalArea)
}
