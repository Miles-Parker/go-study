package main

import (
	"fmt"
	"math"
)

//题目：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。

type Shape interface {
	Area() float64
	Perimeter() float64
}
type Rectangle struct {
	Width  float64
	Height float64
}
type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
func (c Circle) Perimeter() float64 {
	return 2 * (c.Radius + c.Radius)
}

func main() {
	rect := Rectangle{Width: 23, Height: 34}
	circle := Circle{Radius: 23}

	shapes := []Shape{rect, circle}

	for _, shape := range shapes {
		fmt.Printf("图形类型: %T\n", shape)
		fmt.Printf("面积: %.2f\n", shape.Area())
		fmt.Printf("周长: %.2f\n\n", shape.Perimeter())
	}
}
